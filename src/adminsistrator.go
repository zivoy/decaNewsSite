package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type adminLog struct {
	Time        int64                  `json:"time"`
	Uid         string                 `json:"uid"`
	LogLevel    int                    `json:"log_level"`
	Description string                 `json:"description"`
	Information map[string]interface{} `json:"information,omitempty"`
}

func addLog(actionLevel int, actee string, description string, additionalInformation ...map[string]interface{}) {
	location := dataBase.NewRef("admin/logs")
	currTime := time.Now().Unix()
	var details map[string]interface{}
	if len(additionalInformation) > 0 {
		details = additionalInformation[0]
	}
	_, err := location.Push(ctx, &adminLog{
		Time:        currTime,
		Uid:         actee,
		LogLevel:    actionLevel,
		Description: description,
		Information: details,
	})
	if err != nil && debug {
		log.Println(fmt.Errorf("log error: %w", err))
	}
}

func adminBoard(c *gin.Context) {
	render(c, gin.H{
		"autoClearCache": clearingCache,
	},
		"Administrator Dashboard",
		"place to administrate all the things.",
		"",
		c.Request.URL,
		"adminDashboard.html")

}

func clearUserCache(c *gin.Context) {
	uid := c.Param("uid")
	requesterUser, _ := c.Get("user")
	requester := requesterUser.(user).UID

	if requester == "" {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `requester` id"})
		return
	}

	addLog(0, requester, "Refresh Cache", map[string]interface{}{"user_affected": uid})

	userCache.delete(uid)
	c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

func togglePostingPerms(c *gin.Context) {
	uid := c.Param("uid")
	requesterUser, _ := c.Get("user")
	requester := requesterUser.(user).UID

	if requester == "" {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `requester` id"})
		return
	}

	if userExists(uid) {
		user := getUser(uid)
		user.PostingPrivilege = !user.PostingPrivilege
		addUser(uid, user)
		c.JSON(http.StatusOK, map[string]interface{}{"success": true, "perms": user.PostingPrivilege})
		addLog(1, requester, "Toggle Posting Privilege", map[string]interface{}{
			"user_affected": uid,
			"value_set":     user.PostingPrivilege,
		})
		return
	}
	c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "message": "user not found"})
}

func generateAuthButtons(viewerAuth int, viewedAuth int) template.HTML {
	returnString := ""
	for k := 0; k < len(authorities); k++ {
		v := authorities[k]
		returnString += "<button class=\"button "
		if viewedAuth == k {
			returnString += "is-static is-selected "
		}
		returnString += "is-size-7-mobile\""
		if viewerAuth < k {
			returnString += "disabled"
		}
		returnString += "onclick=\"setRank(" + strconv.Itoa(k) + ")\" " +
			"id=\"rank_" + strconv.Itoa(k) + "\">" + v + "</button>"
	}
	return unescape(returnString)
}

// takes `rank` in post request
func UpdateUserRank(c *gin.Context) {
	uid := c.Param("uid")
	requesterUser, _ := c.Get("user")
	requester := requesterUser.(user).UID
	rawRank := c.PostForm("rank")

	if requester == "" {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `requester` id"})
		return
	}

	rank, err := strconv.Atoi(rawRank)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `rank` as integer"})
		return
	}

	if userExists(uid) {
		user := getUser(uid)
		user.AuthLevel = rank
		addUser(uid, user)
		c.JSON(http.StatusOK, map[string]interface{}{"success": true, "auth_level": user.AuthLevel})
		addLog(1, requester, "Update User Authorization", map[string]interface{}{
			"user_affected": uid,
			"value_set":     user.AuthLevel,
		})
		return
	}
	c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "message": "user not found"})
}
