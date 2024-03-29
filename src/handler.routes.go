package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func formatUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Host = domainBase.Host
		c.Request.URL.Scheme = domainBase.Scheme
	}
}

func makeAbsUrl(c *gin.Context, s string) string {
	u, err := url.Parse(s)
	if err != nil {
		log.Println(err)
		return s
	}
	base, err := url.Parse(c.Request.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	return base.ResolveReference(u).String()
}

func pageLogo(c *gin.Context) string {
	logo, _ := url.Parse(c.Request.URL.String())
	logo.Path = "/static/DecaFans-big.png"
	return imagePath(c, parseUrlValues(urlValues{"url": logo.String()}))
}

type urlValues map[string]string

func parseUrlValues(values map[string]string) url.Values {
	r := url.Values{}
	for k, v := range values {
		r.Add(k, v)
	}
	return r
}

func imagePath(c *gin.Context, options ...url.Values) string {
	images, _ := url.Parse(c.Request.URL.String())
	images.Path = "/api/v1/image"
	images.RawQuery = ""

	if len(options) > 0 {
		images.RawQuery = options[0].Encode()
	}

	return images.String()
}

func noRouteFunc(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.JSON(404, getApiError("Page not found"))
	} else {
		abortWithMessage(c, http.StatusNotFound)
	}
}

func healthFunc(c *gin.Context) {
	title := "Server Health"
	template := "health.gohtml"
	if serverHealthy {
		render(c, gin.H{"pageTitle": "Server is healthy",
			"pageSubtitle": "Server is alive!",
			"explanation":  "There are no errors",
		}, title, "Server is alive!", "", c.Request.URL, template)
	} else {
		render(c, gin.H{"pageTitle": "Server is unhealthy",
			"pageSubtitle": "Server is dying",
			"explanation":  "The Cache folder got too big",
			"err":          true,
		}, title, "Server is dead!", "", c.Request.URL, template, http.StatusInternalServerError)
	}
}

func readinessFunc(c *gin.Context) {
	title := "Server Readiness"
	template := "health.gohtml"
	if HearRateAlive {
		render(c, gin.H{"pageTitle": "Server is ready",
			"pageSubtitle": "Server ready to serve",
			"explanation":  "There are no errors",
		}, title, "server is ready!", "", c.Request.URL, template)
	} else {
		render(c, gin.H{"pageTitle": "Server is not ready",
			"pageSubtitle": "Server is having issues",
			"explanation":  "There was an issue with the heartbeat to the database",
			"err":          true,
		}, title, "server is not ready!", "", c.Request.URL, template, http.StatusInternalServerError)
	}
}

func officialPageFunc(c *gin.Context) {
	render(c, gin.H{},
		"Official DecaGear1 news page",
		"Official news from Megadodo about the DecaGear1 headset.",
		"",
		c.Request.URL,
		"official.gohtml")
}

func aboutPageFunc(c *gin.Context) {
	render(c, gin.H{},
		"About",
		"DecaFans is a site maintained and run by fans of the DecaGear1 headset to share the latest news and leaks.",
		"",
		c.Request.URL,
		"about.gohtml")
}

func adminRedirectFunc(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/admin/dashboard")
}
