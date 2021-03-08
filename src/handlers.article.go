package main

import (
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
		"Leak list",
		"List of the latest DecaLeaks.",
		" ",
		c.Request.URL,
		"leakList.html")

}

func leakListFirst(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/leaks/list/1")
}

func showIndex(c *gin.Context) {
	render(c, gin.H{},
		"Home Page",
		"deca news page for getting all the latest DecaLeaks.",
		" ",
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
	}},
		"Create new",
		"Share a new DecaLeak",
		" ",
		c.Request.URL,
		"postLeak.html")
}
