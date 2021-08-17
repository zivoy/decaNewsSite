package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
)

//todo have reload also use refresh token
func performLogin(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "discord")
	c.Request.URL.RawQuery = q.Encode()
	setCookie(c, "login", "in process of logging in", 60*5, false)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func loginCallback(c *gin.Context) {

	deleteCookie(c, "login")
	q := c.Request.URL.Query()
	q.Add("provider", "discord")
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		abortWithMessage(c, http.StatusInternalServerError, err)
		return
	}
	loggInUser(c, user)
	c.Redirect(http.StatusTemporaryRedirect, "/")
	//render(c, gin.H{
	//	"title":       "logged in success!",
	//	"description": "login success",
	//	"url":         fmt.Sprint(c.Request.URL),
	//	"image":       "",
	//	"info":        string(getUserByToken(token)),
	//}, "login-successful.html")

}

func logout(c *gin.Context) {
	token, _ := c.Cookie("token")
	_ = deletePath(dataBase, sessionPathString(token))
	sessionsCache.delete(token)
	// Clear the cookie
	_ = gothic.Logout(c.Writer, c.Request)
	deleteCookie(c, "token")
	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func userProfile(c *gin.Context) {
	uid := c.Param("profile_id")

	if userExists(uid) {
		user := getUser(uid)
		// Call the HTML method of the Context to render a template
		userLeaks, err := getAllUsersArticles(uid)
		if err != nil {
			abortWithMessage(c, http.StatusInternalServerError, err)
		}
		user.RefreshToken = ""
		title := fmt.Sprintf("%s's profile", user.Username)
		if user.Username[len(user.Username)-1] == 's' {
			title = fmt.Sprintf("%s' profile", user.Username)
		}
		render(c, gin.H{"payload": map[string]interface{}{
			"user":      user,
			"leaksMade": userLeaks,
		}},
			title,
			fmt.Sprintf("%s#%s - %s", user.Username, user.UserDiscriminator, authorityLevel(user.AuthLevel)),
			user.AvatarUrl,
			c.Request.URL,
			"profile.html")
	} else {
		// If the profile is not found, abort with an error
		abortWithMessage(c, http.StatusNotFound)
	}
}
