package main

import "github.com/gin-gonic/gin"

//todo have admin actions update cache
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
