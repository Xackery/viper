//Virtually all of this code is from ChimeraCoder's anaconda system for twitter, because it's awesome.
//Check out his link here: https://github.com/ChimeraCoder/anaconda

package tumblr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	_Get  = iota
	_Post = iota
	//BaseURL for all API calls
	BaseURL = "https://api.tumblr.com/v2"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://www.tumblr.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://www.tumblr.com/oauth/authorize",
	TokenRequestURI:               "https://www.tumblr.com/oauth/access_token",
}

//API wraps a session representing a access token/secret, and is the root object to call methods from
type API struct {
	Credentials          *oauth.Credentials
	queryQueue           chan query
	bucket               *tokenbucket.Bucket
	returnRateLimitError bool
	HTTPClient           *http.Client
}

type query struct {
	url        string
	form       url.Values
	data       interface{}
	method     int
	responseCh chan response
}

type response struct {
	data interface{}
	err  error
}

//NewAPI takes an user-specific access token and secret and returns a Api struct for that user.
func NewAPI(accessToken string, accessTokenSecret string) (api *API) {
	//TODO figure out how much to buffer this channel
	//A non-buffered channel will cause blocking when multiple queries are made at the same time
	queue := make(chan query)
	api = &API{
		Credentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},
		queryQueue:           queue,
		bucket:               nil,
		returnRateLimitError: false,
		HTTPClient:           http.DefaultClient,
	}
	go api.throttledQuery()
	return api
}

//SetConsumerKey will set the application-specific consumer_key used in the initial OAuth process
//This key is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerKey(consumerKey string) {
	oauthClient.Credentials.Token = consumerKey
}

//SetConsumerSecret will set the application-specific secret used in the initial OAuth process
//This secret is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerSecret(consumerSecret string) {
	oauthClient.Credentials.Secret = consumerSecret
}

// ReturnRateLimitError specifies behavior when the Twitter API returns a rate-limit error.
// If set to true, the query will fail and return the error instead of automatically queuing and
// retrying the query when the rate limit expires
func (a *API) ReturnRateLimitError(b bool) {
	a.returnRateLimitError = b
}

//EnableThrottling is used to enable query throttling with the tokenbucket algorithm
func (a *API) EnableThrottling(rate time.Duration, bufferSize int64) {
	a.bucket = tokenbucket.NewBucket(rate, bufferSize)
}

//DisableThrottling is used to enable query throttling with the tokenbucket algorithm
func (a *API) DisableThrottling() {
	a.bucket = nil
}

// SetDelay will set the delay between throttled queries
// To turn of throttling, set it to 0 seconds
func (a *API) SetDelay(t time.Duration) {
	a.bucket.SetRate(t)
}

//GetDelay retrives the delay currently set between throttled queries
func (a *API) GetDelay() time.Duration {
	return a.bucket.GetRate()
}

//AuthorizationURL generates the authorization URL for the first part of the OAuth handshake.
//Redirect the user to this URL.
//This assumes that the consumer key has already been set (using SetConsumerKey).
func AuthorizationURL(callback string) (string, *oauth.Credentials, error) {
	tempCred, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, callback, nil)
	if err != nil {
		return "", nil, err
	}
	return oauthClient.AuthorizationURL(tempCred, nil), tempCred, nil
}

//GetCredentials from Oauth provider
func GetCredentials(tempCred *oauth.Credentials, verifier string) (*oauth.Credentials, url.Values, error) {
	return oauthClient.RequestToken(http.DefaultClient, tempCred, verifier)
}

func cleanValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return v
}

// apiGet issues a GET request to the Twitter API and decodes the response JSON to data.
func (a API) apiGet(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(a.HTTPClient, a.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (a API) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Post(a.HTTPClient, a.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, &data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		fmt.Println("non 200 error")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

//query executes a query to the specified url, sending the values specified by form, and decodes the response JSON to data
//method can be either _GET or _POST
func (a API) execQuery(urlStr string, form url.Values, data interface{}, method int) error {
	switch method {
	case _Get:
		return a.apiGet(urlStr, form, data)
	case _Post:
		return a.apiPost(urlStr, form, &data)
	default:
		return fmt.Errorf("HTTP method not yet supported")
	}
}

// throttledQuery executes queries and automatically throttles them according to SECONDS_PER_QUERY
// It is the only function that reads from the queryQueue for a particular *Api struct

func (a *API) throttledQuery() {
	for q := range a.queryQueue {
		url := q.url
		form := q.form
		data := q.data //This is where the actual response will be written
		method := q.method

		responseCh := q.responseCh

		if a.bucket != nil {
			<-a.bucket.SpendToken(1)
		}

		err := a.execQuery(url, form, &data, method)

		// Check if Twitter returned a rate-limiting error
		if err != nil {
			/*if apiErr, ok := err.(*ApiError); ok {
				if isRateLimitError, nextWindow := apiErr.RateLimitCheck(); isRateLimitError && !c.returnRateLimitError {

					// If this is a rate-limiting error, re-add the job to the queue
					// TODO it really should preserve order
					go func() {
						c.queryQueue <- q
					}()

					delay := nextWindow.Sub(time.Now())
					<-time.After(delay)

					// Drain the bucket (start over fresh)
					if c.bucket != nil {
						c.bucket.Drain()
					}

					continue
				}
			}*/
			fmt.Println("Error on query", err.Error())
		}

		responseCh <- response{data, err}
	}
}

// Close query queue
func (a *API) Close() {
	close(a.queryQueue)
}
