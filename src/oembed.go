package main

import (
	"encoding/xml"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var leakRE = regexp.MustCompile("/leaks/leak/(.+)")

type OEmbed struct {
	XMLName      xml.Name `json:"-" xml:"oembed"`
	Type         string   `json:"type" xml:"type"`
	Version      string   `json:"version" xml:"version"`
	Title        string   `json:"title" xml:"title"`
	AuthorName   string   `json:"author_name" xml:"author_name"`
	AuthorUrl    string   `json:"author_url" xml:"author_url"`
	Html         string   `json:"html" xml:"html"`
	ProviderName string   `json:"provider_name" xml:"provider_name"`
	ProviderUrl  string   `json:"provider_url" xml:"provider_url"`
	//cache_age
}

func oembedEndpoint(c *gin.Context) {
	url, uOk := c.GetQuery("url")
	maxWidth, wOk := c.GetQuery("maxwidth")
	maxHeight, hOk := c.GetQuery("maxheight")
	format, fOk := c.GetQuery("format")

	if !uOk {
		c.JSON(http.StatusBadRequest, getApiError("url is missing"))
		return
	}

	if !fOk {
		format = "json"
	}
	format = strings.ToLower(format)

	var maxH, maxW int
	var err error
	if wOk {
		maxW, err = strconv.Atoi(maxWidth)
		if err != nil {
			c.JSON(http.StatusBadRequest, getApiError("\""+maxWidth+"\" is invalid"))
			return
		}
	}
	if hOk {
		maxH, err = strconv.Atoi(maxHeight)
		if err != nil {
			c.JSON(http.StatusBadRequest, getApiError("\""+maxHeight+"\" is invalid"))
			return
		}
	}

	var articleId string
	if leakRE.MatchString(url) {
		articleId = leakRE.FindStringSubmatch(url)[1]
		if _, err = getArticleByID(articleId); err != nil {
			c.JSON(http.StatusNotFound, getApiError(articleId+" is not a valid article"))
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, getApiError(url+" is invalid"))
		return
	}

	oembed := makeOembed(articleId, maxW, maxH)
	if format == "json" {
		c.JSON(http.StatusOK, oembed)
	} else if format == "xml" {
		c.XML(http.StatusOK, oembed)
	} else {
		// not valid format
		c.JSON(http.StatusNotImplemented, getApiError("\""+format+"\" is not supported"))
	}
}

func makeOembed(articleId string, maxWidth, maxHeight int) OEmbed {
	article, _ := getArticleByID(articleId)
	width := int(math.Max(float64(maxWidth), 300))
	height := int(math.Max(float64(maxHeight), 400))

	reporter := getUser(article.ReporterUid)
	return OEmbed{
		Type:         "rich",
		Version:      "1.0",
		Title:        article.Title,
		AuthorName:   reporter.Username + "#" + reporter.UserDiscriminator,
		AuthorUrl:    "https://decafans.com/u/profile/" + reporter.UID,
		Html:         fmt.Sprintf("<img src=\"https://decafans.com/api/v1/image?url=/static/DecaFans-banner.png&height=%d&width=%d\" title=\"hover\">hippo</p>", height, width),
		ProviderName: "DecaFans",
		ProviderUrl:  "https://decafans.com",
	}
}
