package main

import (
	"firebase.google.com/go/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

//todo make tests

var router *gin.Engine

var discordProvider *discord.Provider
var debug bool
var store *sessions.CookieStore

var dataBase *db.Client

var domainBase *url.URL

var authorities = map[int]string{
	0: "Browser",
	1: "Reporter",
	2: "Administrator",
	3: "Creator",
}

func main() {
	var err error
	// go into production mode if there is an error with parsing the bool
	debug, err = strconv.ParseBool(os.Getenv("DEBUG"))

	if err != nil {
		debug = false
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	discordRedirect := os.Getenv("REDIRECT")
	domainBase, err = url.Parse(discordRedirect)
	if err != nil {
		panic(err)
	}

	key := os.Getenv("STORE_SECRET") // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30             // 30 days

	store = sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.Domain = domainBase.Host
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = !debug

	gothic.Store = store

	discordProvider = discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), discordRedirect+"/u/login/callback", discord.ScopeIdentify)
	goth.UseProviders(discordProvider)

	router = gin.Default()

	router.SetFuncMap(template.FuncMap{
		"authLevelName": authorityLevel,
		"getUser":       getUser,
	}) //todo create a function for name colours

	router.LoadHTMLGlob("templates/*")

	// initDB
	err = initializeApp([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		panic(err)
	}
	dataBase, err = initDB()
	if err != nil {
		panic(err)
	}

	// Initialize the routes
	initializeRoutes()

	initCacheClearing()

	err = router.Run(":5000")
	if err != nil {
		panic(err)
	}
}

func authorityLevel(auth int) string {
	//var valid []string
	//for k, v := range authorities {
	//	if auth < k {
	//		break
	//	}
	//	valid = append(valid, v)
	//}
	//return strings.Join(valid[:], ", ")
	return fmt.Sprintf("%d: %s", auth, authorities[auth])
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, title string, description string, image string, url *url.URL, templateName string, status ...int) {
	data["title"] = title
	data["description"] = description
	url.Host = domainBase.Host
	url.Scheme = domainBase.Scheme
	data["url"] = url.String()
	data["image"] = image
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)
	userVals, _ := c.Get("user")
	data["user"] = userVals

	stat := http.StatusOK
	if len(status) > 0 {
		stat = status[0]
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(stat, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(stat, data["payload"])
	default:
		// Respond with HTML
		c.HTML(stat, templateName, data)
	}

}
