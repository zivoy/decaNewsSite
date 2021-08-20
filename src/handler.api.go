package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

type apiError map[string]interface{}

func getApiError(message string) apiError {
	return apiError{"error": true, "message": message}
}

//

func apiRootFunc(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/api/v1")
}

// apiv1

func apiV1RootFunc(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"hello": "world"})
}

func apiV1ImageFunc(c *gin.Context) {
	var width64, height64 int64
	var err error
	var vErr string

	w, wOk := c.GetQuery("w")
	wi, widthOk := c.GetQuery("width")
	h, hOk := c.GetQuery("h")
	hi, heightOk := c.GetQuery("height")

	if wOk || widthOk {
		if wOk {
			width64, err = strconv.ParseInt(w, 10, 32)
			vErr = w
		} else {
			width64, err = strconv.ParseInt(wi, 10, 32)
			vErr = wi
		}

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError(vErr+" in width is invalid"))
			return
		}
	}

	if hOk || heightOk {
		if hOk {
			height64, err = strconv.ParseInt(h, 10, 32)
			vErr = h
		} else {
			height64, err = strconv.ParseInt(hi, 10, 32)
			vErr = hi
		}

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError(vErr+" in height is invalid"))
			return
		}
	}

	width := int(width64)
	height := int(height64)
	var path string

	// url was provided
	if imageUrl, ok := c.GetQuery("url"); ok {
		imageUrl, _ = url.QueryUnescape(imageUrl)
		path, err = getImage(makeAbsUrl(c, imageUrl), width, height)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError("'"+imageUrl+"' is invalid"))
			return
		}
	}

	// leak id was provided
	if id, ok := c.GetQuery("id"); ok {
		path, err = getImageFromId(id, width, height)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError("'"+id+"' is an invalid id"))
			return
		}
		if path == "" {
			c.JSON(http.StatusBadRequest, getApiError("'"+id+"' has no image"))
			return
		}
	}

	if path != "" {
		c.FileFromFS(path, http.FS(cacheFS))
		return
	}

	c.JSON(422, getApiError("No url or leak id provided, please fill either the `url` or `id` params"))
}

/// leaks api

func leaksApiGetFunc(c *gin.Context) {
	var low, high int64
	var err error
	if l, ok := c.GetQuery("low"); !ok {
		low = 0
	} else {
		low, err = strconv.ParseInt(l, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError(l+" in low is invalid"))
			return
		}
	}

	if h, ok := c.GetQuery("high"); !ok {
		high = -1
	} else {
		high, err = strconv.ParseInt(h, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, getApiError(h+" in high is invalid"))
			return
		}
	}

	art, err := getAllArticles(int(low), int(high))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, getApiError("problem fetching articles"))
		return
	}
	c.JSON(http.StatusOK, art)
}

func leaksApiAmountFunc(c *gin.Context) {
	art, err := getAllArticles(0, -1)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, getApiError("problem fetching articles"))
		return
	}
	c.JSON(http.StatusOK, map[string]int{"leaks": len(art)})
}
