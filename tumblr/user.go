package tumblr

import (
	"net/url"
)

type User struct {
}

//Use this method to retrieve the user's account information that matches the OAuth credentials submitted with the request.
func (a *API) GetUserInfo() (user User, err error) {
	return
}

//Use this method to retrieve the dashboard that matches the OAuth credentials submitted with the request.
func (a *API) GetUserDashboard() (posts []Post, err error) {
	return
}

//Use this method to retrieve the liked posts that match the OAuth credentials submitted with the request.
func (a *API) GetUserLikes(params url.Values) (likedPosts []Post, likedCount int, err error) {
	return
}

//Use this method to retrieve the blogs followed by the user whose OAuth credentials are submitted with the request.
func (a *API) GetUserFollowing(params url.Values) (totalBlogs int, blogs []Blog, err error) {
	return
}

//Follow a blog
func (a *API) PostUserFollow(blogUrl string) (err error) {
	return
}

func (a *API) PostUserUnfollow(blogUrl string) (err error) {
	err = a.PostUserFollowDelete(blogUrl)
	return
}

//Unfollow a blog
func (a *API) PostUserFollowDelete(blogUrl string) (err error) {

	return
}

//Like a post
func (a *API) PostUserLikePost(postID int, reblogKey string) (err error) {
	return
}

//Unlike a post
func (a *API) PostUserUnlikePost(postID int, reblogKey string) (err error) {
	err = a.PostUserLikePostDelete(postID, reblogKey)
	return
}

//Unlike a post
func (a *API) PostUserLikePostDelete(postID int, reblogKey string) (err error) {
	return
}
