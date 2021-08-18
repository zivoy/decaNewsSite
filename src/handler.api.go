package main

import (
	"log"
	"net/http"
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
