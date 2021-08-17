package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"log"
	"strings"
	"time"
)

type user struct {
	Username          string `json:"username"`
	UserDiscriminator string `json:"discriminator"`
	UID               string `json:"-"`
	AuthLevel         int    `json:"auth"`
	AvatarUrl         string `json:"avatar"`
	RefreshToken      string `json:"refresh_token,omitempty"`
	PostingPrivilege  bool   `json:"can_post"`
}

// todo add a last login/refresh and and implement the refresh token for users
// todo add display name
//discordProvider.RefreshToken()

type userSession struct {
	Cookie  string `json:"-"`
	UID     string `json:"uid"`
	Expires int64  `json:"expire"`
}

const (
	sessionLocation = "sessions"
	userLocation    = "users"
)

func userPathString(uid string) string {
	return fmt.Sprintf(userLocation+"/%s", uid)
}

func sessionPathString(uid string) string {
	return fmt.Sprintf(sessionLocation+"/%s", uid)
}

func getUser(uid string) user {
	userData := userCache.get(uid, func(uid string) interface{} {
		userData, err := readEntry(dataBase, userPathString(uid))
		if err != nil && debug {
			log.Println(err)
		}
		return user{
			Username:          userData["username"].(string),
			UID:               uid,
			UserDiscriminator: userData["discriminator"].(string),
			RefreshToken:      userData["refresh_token"].(string),
			AuthLevel:         int(userData["auth"].(float64)),
			AvatarUrl:         userData["avatar"].(string),
			PostingPrivilege:  userData["can_post"].(bool),
		}
	})
	return userData.(user)
}

func userExists(uid string) bool {
	exists := userCache.has(uid)
	if !exists {
		exists = pathExists(dataBase, userPathString(uid))
	}
	return exists
}

func addUser(uid string, user user) {
	err := setEntry(dataBase, userPathString(uid), user)
	if err != nil && debug {
		log.Println(err)
	}
	userCache.add(uid, user)
}

func getSession(token string) userSession {
	session := sessionsCache.get(token, func(token string) interface{} {
		sessionData, err := readEntry(dataBase, sessionPathString(token))
		if err != nil && debug {
			log.Println(err)
		}
		return userSession{
			Cookie:  token,
			UID:     sessionData["uid"].(string),
			Expires: int64(sessionData["expire"].(float64)),
		}
	})
	return session.(userSession)
}

func getUserByToken(token string) (user, error) {
	if !isValidSession(token) {
		return user{}, errors.New("invalid token")
	}
	session := getSession(token)
	return getUser(session.UID), nil
}

func isValidSession(token string) bool {
	exists := sessionsCache.has(token)
	if !exists {
		exists = pathExists(dataBase, sessionPathString(token))
	}
	if exists {
		session := getSession(token)
		if session.Expires-time.Now().Unix() > 0 {
			return true
		}
		_ = deletePath(dataBase, sessionPathString(token))
		sessionsCache.delete(token)
	}
	return false
}

func getCookie() string {
	return base64.URLEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

func loggInUser(c *gin.Context, userVals goth.User) {
	cookie := getCookie()
	for ok := isValidSession(cookie); ok; ok = isValidSession(cookie) {
		cookie = getCookie()
	}
	loggedIn := user{
		Username:          userVals.Name,
		UserDiscriminator: userVals.RawData["discriminator"].(string),
		UID:               userVals.UserID,
		AuthLevel:         0,
		AvatarUrl:         strings.Replace(userVals.AvatarURL, ".jpg", ".png", 1),
		RefreshToken:      userVals.RefreshToken,
		PostingPrivilege:  true,
	}
	if userExists(loggedIn.UID) {
		oldUser := getUser(loggedIn.UID)
		loggedIn.AuthLevel = oldUser.AuthLevel
	}
	addUser(loggedIn.UID, loggedIn)

	session := userSession{
		Cookie:  cookie,
		UID:     loggedIn.UID,
		Expires: time.Now().Unix() + int64(store.Options.MaxAge),
	}
	err := setEntry(dataBase, sessionPathString(cookie), session)
	if err != nil && debug {
		log.Println(err)
	}

	sessionsCache.add(cookie, session)
	setCookie(c, "token", cookie, store.Options.MaxAge)
}
