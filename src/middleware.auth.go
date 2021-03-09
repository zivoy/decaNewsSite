package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// This middleware ensures that a request will be aborted with an error
// if the user is not logged in
func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's an error or if the token is empty
		// the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			//if token, err := c.Cookie("token"); err != nil || token == "" {
			abortWithMessage(c, http.StatusUnauthorized)
		}
	}
}

// This middleware ensures that a request will be aborted with an error
// if the user is already logged in
func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's no error or if the token is not empty
		// the user is already logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			// lol this will throw an error telling you to sign in XD
			abortWithMessage(c, http.StatusUnauthorized)
		}
	}
}

// This middleware sets whether the user is logged in or not
func setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil && (token != "" && isValidSession(token)) {
			c.Set("is_logged_in", true)
			user, _ := getUserByToken(token)
			c.Set("user", user)
		} else {
			c.Set("is_logged_in", false)
			c.Set("user", nil)
		}
	}
}

func isLoggingIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("login"); !(err == nil || token != "") {
			abortWithMessage(c, http.StatusUnauthorized)
		}
	}
}

func minAuthLevel(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if logged, ok := c.Get("is_logged_in"); ok && logged.(bool) {
			usera, _ := c.Get("user")
			if usera.(user).AuthLevel < level {
				abortWithMessage(c, http.StatusForbidden)
			}
		} else {
			abortWithMessage(c, http.StatusUnauthorized)
		}
	}
}

func canPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		usera, _ := c.Get("user")
		if !usera.(user).PostingPrivilege {
			render(c, gin.H{
				"payload":       map[string]string{"error": "no posting privilege!"},
				"errorTitle":    "Invalid Permission",
				"errorSubtitle": "You do not have posting permission",
				"explanation":   "You can ask an admin to grant you posting privileges",
			},
				"Cannot Post",
				"No posting privilege.",
				pageLogo(c),
				c.Request.URL,
				"error.html", http.StatusOK)
			c.Abort()
		}
	}
}

func abortWithMessage(c *gin.Context, code int, erre ...error) {
	var err error
	err = nil
	if len(erre) > 0 {
		err = erre[0]
	}
	switch code {
	case http.StatusNotFound:
		render(c, gin.H{
			"payload":       map[string]string{"error": "Page not found!", "url": c.Request.URL.String()},
			"errorTitle":    "404!",
			"errorSubtitle": "Page not found",
			"explanation":   unescape(fmt.Sprintf("no such page <code>%s</code> was found", c.Request.URL)),
		},
			"404 - page not found",
			"Page not found.",
			pageLogo(c),
			c.Request.URL,
			"error.html", code)
	case http.StatusUnauthorized:
		render(c, gin.H{
			"payload":       map[string]string{"error": "Not signed in."},
			"errorTitle":    "401!",
			"errorSubtitle": "Unauthorised access.",
			"explanation": unescape("You need to sign in to access this.</p><p>" +
				"<a class=\"button is-light\" href=\"/u/login\"><span class=\"icon-text\">" +
				"<span class=\"icon has-text-dark\"><ion-icon class=\"ion-ionic\" name=\"logo-discord\">" +
				"</ion-icon></span><span>Log in</span></span></a>"),
		},
			"401 - unauthorised access",
			"You are not signed in.",
			pageLogo(c),
			c.Request.URL,
			"error.html", code)
	case http.StatusForbidden:
		render(c, gin.H{
			"payload":       map[string]string{"error": "insufficient permissions"},
			"errorTitle":    "403!",
			"errorSubtitle": "insufficient privilege.",
			"explanation":   "You dont have the needed authority level to access this page.",
		},
			"403 - insufficient privileges",
			"Insufficient privileges to access page.",
			pageLogo(c),
			c.Request.URL,
			"error.html", code)
	case http.StatusBadRequest:
		render(c, gin.H{
			"payload":       map[string]string{"error": "page does not exist"},
			"errorTitle":    "400!",
			"errorSubtitle": "Bad Request.",
			"explanation":   "The page you are trying to access is not valid.",
		},
			"400 - bad request",
			"Url is invalid.",
			pageLogo(c),
			c.Request.URL,
			"error.html", code)
	default:
		if err != nil {
			_ = c.AbortWithError(code, err)
			return
		}
		c.AbortWithStatus(code)
		return
	}
	if err != nil {
		_ = c.Error(err)
	}
	c.Abort()
}
