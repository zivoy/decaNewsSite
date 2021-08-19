package main

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"firebase.google.com/go/db"
	"github.com/Masterminds/sprig"
	"github.com/frustra/bbcode"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"

	"log"
	"net/http"
	"net/url"
	"os"
)

//todo make tests

var router *gin.Engine

var discordProvider *discord.Provider
var debug bool
var store *sessions.CookieStore

var dataBase *db.Client

var domainBase *url.URL
var BBCompiler bbcode.Compiler

var version string

var ServerComms bool   // for turning off the heartbeat routine
var HearRateAlive bool // to know if the server can reach the database

var authorities = map[int]string{
	0: "Browser",
	1: "Reporter",
	2: "Administrator",
	3: "Creator",
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dev, err := strconv.ParseBool(os.Getenv("DEV_MODE"))
	if err != nil {
		debug = false
	}
	version = os.Getenv("VERSION")
	if version == "" {
		version = "UNVERSIONED"
	}

	var confMap map[string]string
	if dev {
		// on dev mode make a .env file with your configs
		confMap, err = godotenv.Read(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	} else {
		var conf string
		location := os.Getenv("LOCATION")
		serverPassword := os.Getenv("SERVER_PASSWORD")
		filePassword := os.Getenv("FILE_PASSWORD")
		if location == "" || serverPassword == "" || filePassword == "" {
			log.Println(fmt.Sprintf("Missing credentials\n"+
				"\tURL:         %t\n"+
				"\tServer pass: %t\n"+
				"\tFile pass:   %t", location != "", serverPassword != "", filePassword != ""))
			failedSimplePage()
			return
		}
		conf, err = getConfiguration(location, serverPassword, filePassword)
		if err != nil {
			log.Println(err)
			failedSimplePage()
			return
		}
		confMap, err = godotenv.Unmarshal(conf)
		if err != nil {
			log.Println(err)
			failedSimplePage()
			return
		}
	}

	debug, err = strconv.ParseBool(confMap["DEBUG"])
	if err != nil {
		debug = false
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	discordRedirect := confMap["REDIRECT"]
	domainBase, err = url.Parse(discordRedirect)
	if err != nil {
		panic(err)
	}

	key := confMap["STORE_SECRET"]
	maxAge := 86400 * 30 // 30 days

	store = sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.Domain = domainBase.Host
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = !debug

	gothic.Store = store

	discordProvider = discord.New(confMap["DISCORD_KEY"], confMap["DISCORD_SECRET"],
		discordRedirect+"/u/login/callback", discord.ScopeIdentify)
	goth.UseProviders(discordProvider)

	router = gin.Default()

	functions := sprig.GenericFuncMap()
	functions["authLevelName"] = authorityLevel
	functions["getUser"] = getUser
	functions["makeButtonList"] = generateAuthButtons
	functions["unescape"] = unescape
	functions["compileBB"] = compileBBCode
	functions["hashString"] = hashTo32

	router.SetFuncMap(functions)

	router.LoadHTMLGlob("templates/*")

	// initDB
	err = initializeApp([]byte(confMap["GOOGLE_APPLICATION_CREDENTIALS"]), confMap["DATABASE_URL"])
	if err != nil {
		panic(err)
	}
	dataBase, err = initDB()
	if err != nil {
		panic(err)
	}

	// Initialize the routes
	initializeRoutes()

	setAutoClear(true)
	startServerComms()

	BBCompiler = bbcode.NewCompiler(true, true)
	initBBCode(&BBCompiler)

	srv := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	shutdown(srv)
}

func shutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Block
	<-interruptChan

	// shutdown looping routines
	stopClearingCache()
	ServerComms = false

	// close the firebase connection

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("--- Shutting down ---")
	os.Exit(0)
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
	authLevel, ok := authorities[auth]
	if !ok {
		authLevel = "Invalid"
	}
	return authLevel
	//fmt.Sprintf("%d: %s", auth, authorities[auth])
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, title string, description string, image string, url *url.URL, templateName string, status ...int) {
	data["title"] = fmt.Sprintf("%s - DecaFans", title)
	data["description"] = description
	data["url"] = url.String()
	data["image"] = image
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)
	userVals, _ := c.Get("user")
	data["user"] = userVals
	data["version"] = version

	data["logo"] = pageLogo(c)

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
