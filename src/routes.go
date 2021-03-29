package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func formatUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Host = domainBase.Host
		c.Request.URL.Scheme = domainBase.Scheme
	}
}

func pageLogo(c *gin.Context) string {
	logo, _ := url.Parse(c.Request.URL.String())
	logo.Path = "/static/DecaFans-big.png"
	return logo.String()
}

func initializeRoutes() {
	router.Use(setUserStatus())
	router.Use(formatUrl())

	router.NoRoute(func(c *gin.Context) {
		abortWithMessage(c, http.StatusNotFound)
	})

	router.GET("/", showIndex)

	router.Static("/static", "./resources")
	router.StaticFile("/favicon.png", "./resources/decafansLogoSmall.png")
	router.StaticFile("/favicon.ico", "./resources/DecaFans-favicon.ico")
	router.StaticFile("/robots.txt", "./resources/robots.txt")

	router.GET("/official", func(c *gin.Context) {
		render(c, gin.H{},
			"Official DecaGear1 news page",
			"Official news from Megadodo about the DecaGear1 headset.",
			"",
			c.Request.URL,
			"official.html")
	})
	router.GET("/about", func(c *gin.Context) {
		render(c, gin.H{},
			"About",
			"DecaFans is a site maintained and run by fans of the DecaGear1 headset to share the latest news and leaks.",
			"",
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

		leaksRoutes.GET("/create", ensureLoggedIn(), canPost(), showArticleCreationPage)
		leaksRoutes.POST("/create", ensureLoggedIn(), canPost(), createArticle)
		leaksRoutes.POST("/update/:leak_id", ensureLoggedIn(), canPost(), updateArticle)
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

	apiRoot := router.Group("/api")
	{
		apiRoot.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/api/v1")
		})
		apiV1 := apiRoot.Group("/v1")
		{
			apiV1.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, map[string]string{"hello": "world"})
			})
			apiV1.POST("/archive/:uid", ensureLoggedIn(), archiveLeak)
		}
	}
}
