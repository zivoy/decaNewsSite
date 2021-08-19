package main

func initializeRoutes() {
	router.Use(setUserStatus())
	router.Use(formatUrl())

	router.NoRoute(noRouteFunc)

	router.GET("/", showIndex)

	router.GET("/health", healthFunc)

	router.GET("/readiness", readinessFunc)

	router.Static("/static", "./resources")
	router.StaticFile("/favicon.png", "./resources/decafansLogoSmall.png")
	router.StaticFile("/favicon.ico", "./resources/DecaFans-favicon.ico")
	router.StaticFile("/robots.txt", "./resources/robots.txt")

	router.GET("/official", officialPageFunc)
	router.GET("/about", aboutPageFunc)

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
		admin.GET("/", adminRedirectFunc)
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
		apiRoot.GET("/", apiRootFunc)
		apiV1 := apiRoot.Group("/v1")
		{
			apiV1.GET("/", apiV1RootFunc)

			apiV1.GET("/image", apiV1ImageFunc)

			leekApi := apiV1.Group("/leaks")
			{
				leekApi.GET("/get", leaksApiGetFunc)
				leekApi.GET("/amount", leaksApiAmountFunc)
			}
			apiV1.POST("/archive/:uid", ensureLoggedIn(), archiveLeak)
		}
	}
}
