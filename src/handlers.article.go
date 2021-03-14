package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

func leakList(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		abortWithMessage(c, http.StatusBadRequest)
		return
	}
	page = int(math.Min(float64(page), 1))
	articles, err := getAllArticles(page)
	if err != nil {
		abortWithMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Call the render function with the name of the template to render
	render(c, gin.H{"payload": articles},
		"List of Leaks",
		"List of the latest DecaLeaks.",
		"",
		c.Request.URL,
		"leakList.html")

}

func leakListFirst(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/leaks/list/1")
}

func showIndex(c *gin.Context) {
	render(c, gin.H{},
		"Home Page",
		"DecaFans is a news page run and maintained by fans of the DecaGear1 headset to help getting all the latest information.",
		"",
		c.Request.URL,
		"index.html")

}

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	articleID := c.Param("article_id")
	// Check if the article exists
	if article, err := getArticleByID(articleID); err == nil {
		// Call the HTML method of the Context to render a template
		render(c, gin.H{"payload": article}, "DecaLeak", article.Summary, article.ImageUrl, c.Request.URL,
			"leak.html")
	} else {
		// If the article is not found, abort with an error
		abortWithMessage(c, http.StatusNotFound, err)
	}
}

func showArticleCreationPage(c *gin.Context) {
	render(c, gin.H{"payload": map[string][]string{
		"allowed_links": getAllowedLink(),
	}, "linkLessAuthLevel": linkLessAuthLevel},
		"Create new",
		"Share a new DecaLeak",
		"",
		c.Request.URL,
		"postLeak.html")
}

func archiveLeak(c *gin.Context) {
	uid := c.Param("uid")
	requesterUser, _ := c.Get("user")
	requester := requesterUser.(user).UID

	if !articleExists(uid) {
		abortWithMessage(c, http.StatusBadRequest)
		return
	}

	leak, err := getArticleByID(uid)
	if err != nil && debug {
		fmt.Println(err)
	}
	err = setEntry(dataBase, fmt.Sprintf("leaks/%s", uid), nil)
	if err != nil && debug {
		fmt.Println(err)
	}
	err = setEntry(dataBase, fmt.Sprintf("admin/archived_leaks/%s", uid), leak)
	if err != nil && debug {
		fmt.Println(err)
	}

	addLog(1, requester, "Archiving Leak", map[string]interface{}{"leak_id": uid})

}
