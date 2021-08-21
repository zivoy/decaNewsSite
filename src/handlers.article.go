package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func leakList(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		abortWithMessage(c, http.StatusBadRequest)
		return
	}
	page = int(math.Max(float64(page), 1))
	articles, err := getAllArticles(0, -1)
	if err != nil {
		abortWithMessage(c, http.StatusInternalServerError, err)
		return
	}

	// Call the render function with the name of the template to render
	render(c, gin.H{"payload": articles, "amount": len(articles), "page": page},
		"List of Leaks",
		"List of the latest DecaLeaks.",
		"",
		c.Request.URL,
		"leakList.gohtml")

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
		"index.gohtml")

}

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	articleID := c.Param("article_id")
	// Check if the article exists
	if article, err := getArticleByID(articleID); err == nil {
		// Call the HTML method of the Context to render a template
		render(c, gin.H{"payload": article, "allowed_links": allowedLinksForUserContext(c)},
			article.Title,
			strings.Trim(strings.ReplaceAll(article.Summary, "\n", " "), " "),
			imagePath(c, parseUrlValues(urlValues{"id": article.ID})),
			c.Request.URL, "leak.gohtml")
	} else {
		// If the article is not found, abort with an error
		abortWithMessage(c, http.StatusNotFound, err)
	}
}

func showArticleCreationPage(c *gin.Context) {
	render(c, gin.H{"payload": map[string][]string{
		"allowed_links": allowedLinksForUserContext(c),
	}, "linkLessAuthLevel": linkLessAuthLevel},
		"Create new",
		"Share a new DecaLeak",
		"",
		c.Request.URL,
		"postLeak.gohtml")
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
		log.Println(err)
	}
	err = setEntry(dataBase, fmt.Sprintf(archivedArticleLocation+"/%s", uid), leak)
	if err != nil && debug {
		log.Println(err)
	}
	err = setEntry(dataBase, articleCache.path(uid), nil)
	if err != nil && debug {
		log.Println(err)
	}

	addLog(1, requester, "Archiving Leak", map[string]interface{}{"leak_id": uid})

}

func updateArticle(c *gin.Context) {
	leakID := c.Param("leak_id")
	var leak article
	var err error
	if leak, err = getArticleByID(leakID); err != nil {
		abortWithMessage(c, http.StatusNotFound, err)
		return
	}

	description := c.PostForm("description")
	leakTime := c.PostForm("time")
	imageUrl := c.PostForm("image_url")
	sourceUrl := c.PostForm("source_url")
	title := c.PostForm("title")
	tags := c.PostForm("tags")

	updaterUser, _ := c.Get("user")
	updater := updaterUser.(user)

	if description == "" {
		description = leak.Description
	}
	if leakTime == "" {
		leakTime = strconv.Itoa(int(leak.LeakTime))
	}
	if title == "" {
		title = fmt.Sprintf("DecaLeak %d", hashTo32(leakID))
	}

	// nothing was changed
	if !(leak.Description != description || leak.ImageUrl != imageUrl || strconv.Itoa(int(leak.LeakTime)) != leakTime ||
		leak.SourceLink != sourceUrl || leak.Title != title || !compareTagList(getTagsFromString(tags), leak.Tags)) {
		return
	}

	newLeak, code := leakSanitization(title, description, leakTime, imageUrl, sourceUrl, tags,
		getUser(leak.ReporterUid), updater, leak.DateCreate, time.Now().Unix())
	newLeak.ID = leak.ID
	switch code {
	case 1:
	case 2:
	case 3:
	case 4:
		log.Println("failed updating with code ", code)
		abortWithMessage(c, 406)
		return
	}

	err = setEntry(dataBase, articleCache.path(leak.ID), newLeak)
	if err != nil {
		if debug {
			log.Println(err)
		}
		abortWithMessage(c, 500, err)
		return
	}

	addLog(2, updater.UID, "Updated Leak",
		map[string]interface{}{"article": leak.ID, "before": leak, "after": newLeak})
	articleCache.delete(leak.ID)

	clearData := cacheAction{
		CacheListId: articleListCache.id,
		ItemId:      "articles",
		ActionType:  clearList,
	}
	sendAction(clearData)
	articleListCache.clear()

	c.JSON(200, map[string]string{"success": "true"})
}
