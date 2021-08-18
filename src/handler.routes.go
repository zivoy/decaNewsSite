package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func noRouteFunc(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.JSON(404, getApiError("Page not found"))
	} else {
		abortWithMessage(c, http.StatusNotFound)
	}
}

func healthFunc(c *gin.Context) {
	render(c, gin.H{"pageTitle": "Server is healthy",
		"pageSubtitle": "Server is alive!",
		"explanation":  "There are no errors",
	}, "Server Health", "Server is alive!", "", c.Request.URL, "health.html")
}

func readinessFunc(c *gin.Context) {
	title := "Server Readiness"
	template := "health.html"
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
		"official.html")
}

func aboutPageFunc(c *gin.Context) {
	render(c, gin.H{},
		"About",
		"DecaFans is a site maintained and run by fans of the DecaGear1 headset to share the latest news and leaks.",
		"",
		c.Request.URL,
		"about.html")
}

func adminRedirectFunc(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/admin/dashboard")
}
