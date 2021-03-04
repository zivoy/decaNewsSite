package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

//todo add admin logs

func adminBoard(c *gin.Context) {
	render(c, gin.H{
		"autoClearCache": clearingCache,
	},
		"Administrator Dashboard",
		"place to administrate all the things.",
		" ",
		c.Request.URL,
		"adminDashboard.html")

}

// takes `requester` in post request
func clearUserCache(c *gin.Context) {
	uid := c.Param("uid")
	requester := c.PostForm("requester")

	if requester == "" {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `requester` id"})
		return
	}

	deleteCache(userCache, uid)
	c.JSON(http.StatusOK, map[string]interface{}{"success": true})
}

// takes `requester` in post request
func togglePostingPerms(c *gin.Context) {
	uid := c.Param("uid")
	requester := c.PostForm("requester")

	if requester == "" {
		c.JSON(http.StatusOK, map[string]interface{}{"success": false, "message": "need `requester` id"})
		return
	}

	if userExists(uid) {
		user := getUser(uid)
		user.PostingPrivilege = !user.PostingPrivilege
		addUser(uid, user)
		c.JSON(http.StatusOK, map[string]interface{}{"success": true, "perms": user.PostingPrivilege})
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

// takes `rank` and `requester` in post request
func UpdateUserRank(c *gin.Context) {
	uid := c.Param("uid")
	requester := c.PostForm("requester")
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
		return
	}
	c.JSON(http.StatusNotFound, map[string]interface{}{"success": false, "message": "user not found"})
}
