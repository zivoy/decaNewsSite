package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, apiError{"error": true, "message": "Page not found"})
		} else {
			abortWithMessage(c, http.StatusNotFound)
		}
	})

	router.GET("/", showIndex)

	router.GET("/health", func(c *gin.Context) {
		render(c, gin.H{"pageTitle": "Server is healthy",
			"pageSubtitle": "Server is alive!",
			"explanation":  "There are no errors",
		}, "Server Health", "Server is alive!", "", c.Request.URL, "health.html")
	})

	router.GET("/readiness", func(c *gin.Context) {
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
	})

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
			leekApi := apiV1.Group("/leaks")
			{
				leekApi.GET("/get", func(c *gin.Context) {
					var low, high int64
					var err error
					if l, ok := c.GetQuery("low"); !ok {
						low = 0
					} else {
						low, err = strconv.ParseInt(l, 10, 32)
						if err != nil {
							log.Println(err)
							c.JSON(http.StatusBadRequest, apiError{"error": true, "message": l + " in low is invalid"})
							return
						}
					}

					if h, ok := c.GetQuery("high"); !ok {
						high = -1
					} else {
						high, err = strconv.ParseInt(h, 10, 32)
						if err != nil {
							log.Println(err)
							c.JSON(http.StatusBadRequest, apiError{"error": true, "message": h + " in high is invalid"})
							return
						}
					}

					art, err := getAllArticles(int(low), int(high))
					if err != nil {
						log.Println(err)
						c.JSON(http.StatusInternalServerError, apiError{"error": true, "message": "problem fetching articles"})
						return
					}
					c.JSON(http.StatusOK, art)
				})
				leekApi.GET("/amount", func(c *gin.Context) {
					art, err := getAllArticles(0, -1)
					if err != nil {
						log.Println(err)
						c.JSON(http.StatusInternalServerError, apiError{"error": true, "message": "problem fetching articles"})
						return
					}
					c.JSON(http.StatusOK, map[string]int{"hello": len(art)})
				})
			}
			apiV1.POST("/archive/:uid", ensureLoggedIn(), archiveLeak)
		}
	}
}
