package tumblr

import (
	"net/url"
)

type User struct {
}

//Use this method to retrieve the user's account information that matches the OAuth credentials submitted with the request.
func (a *Api) GetUserInfo() (user User, err error) {
	return
}

//Use this method to retrieve the dashboard that matches the OAuth credentials submitted with the request.
func (a *Api) GetUserDashboard() (posts []Post, err error) {
	return
}

//Use this method to retrieve the liked posts that match the OAuth credentials submitted with the request.
func (a *Api) GetUserLikes(params url.Values) (likedPosts []Post, likedCount int, err error) {
	return
}

//Use this method to retrieve the blogs followed by the user whose OAuth credentials are submitted with the request.
func (a *Api) GetUserFollowing(params url.Values) (totalBlogs int, blogs []Blog, err error) {
	return
}

//Follow a blog
func (a *Api) PostUserFollow(blogUrl string) (err error) {
	return
}

func (a *Api) PostUserUnfollow(blogUrl string) (err error) {
	err = a.PostUserFollowDelete(blogUrl)
	return
}

//Unfollow a blog
func (a *Api) PostUserFollowDelete(blogUrl string) (err error) {

	return
}

//Like a post
func (a *Api) PostUserLikePost(postID int, reblogKey string) (err error) {
	return
}

//Unlike a post
func (a *Api) PostUserUnlikePost(postID int, reblogKey string) (err error) {
	err = a.PostUserLikePostDelete(postID, reblogKey)
	return
}

//Unlike a post
func (a *Api) PostUserLikePostDelete(postID int, reblogKey string) (err error) {
	return
}
