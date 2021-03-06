//https://www.tumblr.com/docs/en/API/v2#blog-info
package tumblr

import (
	"fmt"
	"net/url"
)

type Blog struct {
	Title               string
	Name                string
	Posts               int
	Url                 string
	Updated             int
	Description         string
	IsNSFW              bool `json:"is_nsfw"`
	Ask                 bool
	AskPageTitle        string
	AskAnon             bool `json:"ask_anon"`
	SubmissionPageTitle string
	ShareLikes          bool
	Likes               int
}

type ResponseContainer struct {
	Blog Blog `json:"blog"`
}

type BodyContainer struct {
	Response ResponseContainer `json:"response"`
}

//This method returns general information about the blog, such as the title, number of posts, and other high-level data.
func (a *API) GetBlogInfo(baseHostname string) (blog Blog, err error) {
	body := BodyContainer{}

	responseCh := make(chan response)
	a.queryQueue <- query{BaseURL + "/blog/" + baseHostname + "/info", nil, &body, _Get, responseCh}
	fmt.Println(body)

	blog = body.Response.Blog
	return blog, (<-responseCh).err
}

//You can get a blog's avatar in 9 different sizes. The default size is 64x64.
func (a *API) GetBlogAvatar(baseHostname string, size int) (avatarUrl string, err error) {

	return
}

//This method can be used to retrieve the publicly exposed likes from a blog.
func (a *API) GetBlogLikes(baseHostname string, params url.Values) (likedPosts []Post, likedCount int, err error) {

	return
}

//Retrieve a Blog's Followers
func (a *API) GetBlogFollowers(baseHostname string, params url.Values) (users []User, err error) {

	return
}

//Retrieve Published Posts
func (a *API) GetBlogPosts(baseHostname string, params url.Values) (blogs []Blog, posts []Post, totalPosts int, err error) {

	return
}

// Retrieve Queued Posts
func (a *API) GetBlogPostsQueue(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *API) GetBlogPostsDraft(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *API) GetBlogPostsSubmission(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *API) PostBlog(bastHostname string, params url.Values) (err error) {

	return
}

func (a *API) PostBlogEdit(baseHostname string, postID int) (err error) {

	return
}

func (a *API) PostBlogReblog(baseHostname string, params url.Values) (err error) {

	return
}

func (a *API) PostBlogDelete(baseHostname string, postID int) (err error) {

	return
}
