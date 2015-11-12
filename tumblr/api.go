//Virtually all of this code is from ChimeraCoder's anaconda system for twitter, because it's awesome.
//Check out his link here: https://github.com/ChimeraCoder/anaconda

package tumblr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	_GET    = iota
	_POST   = iota
	BaseUrl = "https://api.twitter.com/1.1"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://www.tumblr.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://www.tumblr.com/oauth/authorize",
	TokenRequestURI:               "https://www.tumblr.com/oauth/access_token",
}

type Api struct {
	Credentials          *oauth.Credentials
	queryQueue           chan query
	bucket               *tokenbucket.Bucket
	returnRateLimitError bool
	HttpClient           *http.Client
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

const DEFAULT_DELAY = 0 * time.Second
const DEFAULT_CAPACITY = 5

//NewApi takes an user-specific access token and secret and returns a Api struct for that user.
//The Api struct can be used for accessing any of the endpoints available.
func NewApi(accessToken string, accessTokenSecret string) *Api {
	//TODO figure out how much to buffer this channel
	//A non-buffered channel will cause blocking when multiple queries are made at the same time
	queue := make(chan query)
	c := &Api{
		Credentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},
		queryQueue:           queue,
		bucket:               nil,
		returnRateLimitError: false,
		HttpClient:           http.DefaultClient,
	}
	go c.throttledQuery()
	return c
}

//SetConsumerKey will set the application-specific consumer_key used in the initial OAuth process
//This key is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerKey(consumer_key string) {
	oauthClient.Credentials.Token = consumer_key
}

//SetConsumerSecret will set the application-specific secret used in the initial OAuth process
//This secret is listed on https://dev.twitter.com/apps/YOUR_APP_ID/show
func SetConsumerSecret(consumer_secret string) {
	oauthClient.Credentials.Secret = consumer_secret
}

// ReturnRateLimitError specifies behavior when the Twitter API returns a rate-limit error.
// If set to true, the query will fail and return the error instead of automatically queuing and
// retrying the query when the rate limit expires
func (c *Api) ReturnRateLimitError(b bool) {
	c.returnRateLimitError = b
}

// Enable query throttling using the tokenbucket algorithm
func (c *Api) EnableThrottling(rate time.Duration, bufferSize int64) {
	c.bucket = tokenbucket.NewBucket(rate, bufferSize)
}

// Disable query throttling
func (c *Api) DisableThrottling() {
	c.bucket = nil
}

// SetDelay will set the delay between throttled queries
// To turn of throttling, set it to 0 seconds
func (c *Api) SetDelay(t time.Duration) {
	c.bucket.SetRate(t)
}

func (c *Api) GetDelay() time.Duration {
	return c.bucket.GetRate()
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
func (c Api) apiGet(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Get(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// apiPost issues a POST request to the Twitter API and decodes the response JSON to data.
func (c Api) apiPost(urlStr string, form url.Values, data interface{}) error {
	resp, err := oauthClient.Post(c.HttpClient, c.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

// decodeResponse decodes the JSON response from the Twitter API.
func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		fmt.Println("err!")
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

//query executes a query to the specified url, sending the values specified by form, and decodes the response JSON to data
//method can be either _GET or _POST
func (c Api) execQuery(urlStr string, form url.Values, data interface{}, method int) error {
	switch method {
	case _GET:
		return c.apiGet(urlStr, form, data)
	case _POST:
		return c.apiPost(urlStr, form, data)
	default:
		return fmt.Errorf("HTTP method not yet supported")
	}
}

// throttledQuery executes queries and automatically throttles them according to SECONDS_PER_QUERY
// It is the only function that reads from the queryQueue for a particular *Api struct

func (c *Api) throttledQuery() {
	for q := range c.queryQueue {
		url := q.url
		form := q.form
		data := q.data //This is where the actual response will be written
		method := q.method

		responseCh := q.responseCh

		if c.bucket != nil {
			<-c.bucket.SpendToken(1)
		}

		err := c.execQuery(url, form, data, method)

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
func (c *Api) Close() {
	close(c.queryQueue)
}
