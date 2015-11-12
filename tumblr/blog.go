//https://www.tumblr.com/docs/en/api/v2#blog-info
package tumblr

import (
	"net/url"
)

type Blog struct {
}

//This method returns general information about the blog, such as the title, number of posts, and other high-level data.
func (a *Api) GetBlogInfo(baseHostname string) (blog Blog, err error) {

	responseCh := make(chan response)
	a.queryQueue <- query{BaseUrl + "/blog/" + baseHostname + "/info", nil, &blog, _GET, responseCh}
	return blog, (<-responseCh).err
}

//You can get a blog's avatar in 9 different sizes. The default size is 64x64.
func (a *Api) GetBlogAvatar(baseHostname string, size int) (avatarUrl string, err error) {

	return
}

//This method can be used to retrieve the publicly exposed likes from a blog.
func (a *Api) GetBlogLikes(baseHostname string, params url.Values) (likedPosts []Post, likedCount int, err error) {

	return
}

//Retrieve a Blog's Followers
func (a *Api) GetBlogFollowers(baseHostname string, params url.Values) (users []User, err error) {

	return
}

//Retrieve Published Posts
func (a *Api) GetBlogPosts(baseHostname string, params url.Values) (blogs []Blog, posts []Post, totalPosts int, err error) {

	return
}

// Retrieve Queued Posts
func (a *Api) GetBlogPostsQueue(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *Api) GetBlogPostsDraft(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *Api) GetBlogPostsSubmission(baseHostname string, params url.Values) (blogs []Blog, posts []Post, err error) {

	return
}

func (a *Api) PostBlog(bastHostname string, params url.Values) (err error) {

	return
}

func (a *Api) PostBlogEdit(baseHostname string, postID int) (err error) {

	return
}

func (a *Api) PostBlogReblog(baseHostname string, params url.Values) (err error) {

	return
}

func (a *Api) PostBlogDelete(baseHostname string, postID int) (err error) {

	return
}
