package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initializeRoutes() {
	router.Use(setUserStatus())

	router.NoRoute(func(c *gin.Context) {
		abortWithMessage(c, http.StatusNotFound)
	})

	router.GET("/", showIndex)

	router.Static("/static", "./resources")
	router.StaticFile("/favicon.png", "./resources/cropped-deca_transparent_logo_clean_square-32x32.png")
	router.StaticFile("/favicon.ico", "./resources/DecaFans-favicon.ico")

	router.GET("/official", func(c *gin.Context) {
		render(c, gin.H{},
			"Official news page",
			"Official news from deca.",
			" ",
			c.Request.URL,
			"official.html")
	})
	router.GET("/about", func(c *gin.Context) {
		render(c, gin.H{},
			"About DecaFans",
			"DecaFans is a site for fans of the DecaGear headset to share the latest news.",
			" ",
			c.Request.URL,
			"about.html")
	})

	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", ensureNotLoggedIn(), performLogin)
		userRoutes.GET("/login/callback", isLoggingIn(), loginCallback)

		userRoutes.GET("/logout", ensureLoggedIn(), logout)

		//view profile
		userRoutes.GET("/profile/:profile_id", userProfile)
	}

	// Group article related routes together
	leaksRoutes := router.Group("/leaks")
	{
		leaksRoutes.GET("/leak/:article_id", getArticle)
		leaksRoutes.GET("/", leakListFirst)
		leaksRoutes.GET("/list/:page", leakList)

		leaksRoutes.GET("/create", ensureLoggedIn(), showArticleCreationPage)
		leaksRoutes.POST("/create", ensureLoggedIn(), createArticle)
	}

	admin := router.Group("/admin")
	admin.Use(minAuthLevel(2))
	{
		admin.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "/admin/dashboard")
		})
		admin.GET("/dashboard", adminBoard)

		adminApi := admin.Group("/api")
		{
			adminApi.POST("clearCache/user/:uid", clearUserCache)
			adminApi.POST("togglePosting/:uid", togglePostingPerms)
			adminApi.POST("updateRank/:uid", UpdateUserRank)
		}
	}
}
