package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type article struct {
	ID          string `json:"-"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	LeakTime    int64  `json:"time"`
	ImageUrl    string `json:"image_url"`
	SourceLink  string `json:"source_url"`
	ReporterUid string `json:"reporter_uid"`
	EditedBy    string `json:"edited_by,omitempty"`
	Title       string `json:"title"`
	DateCreate  int64  `json:"created_time"`
	DateEdit    int64  `json:"edited_time,omitempty"`
}

const (
	articleLocation         = "leaks"
	adminBasePath           = "admin"
	archivedArticleLocation = adminBasePath + "/archived_leaks"
	allowedLinkLocation     = adminBasePath + "/allowed_links"
)

// you need this auth level to post with no link
const linkLessAuthLevel = 1

func getAllArticles(low, high int) ([]article, error) {
	var err error
	returned := articleListCache.get("articles", func(_ string) interface{} {
		ref := dataBase.NewRef(articleLocation)
		var data map[string]article
		if err = ref.Get(ctx, &data); err != nil {
			return nil
		}

		articleList := make([]article, 0)
		for k, v := range data {
			v.ID = k
			articleList = append(articleList, v)
			articleCache.add(k, v)
		}

		sort.Slice(articleList, func(i, j int) bool {
			return articleList[i].LeakTime > articleList[j].LeakTime
		})
		return articleList
	})
	if returned == nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}

	articleList := returned.([]article)
	if high == -1 {
		high = len(articleList)
	}
	return articleList[low:high], nil
}

func getAllUsersArticles(uid string) ([]article, error) {
	ref := dataBase.NewRef(articleLocation)
	var data map[string]article
	if err := ref.OrderByChild("reporter_uid").EqualTo(uid).Get(ctx, &data); err != nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}

	articleList := make([]article, 0)
	for k, v := range data {
		v.ID = k
		articleList = append(articleList, v)
		articleCache.add(k, v)
	}

	sort.Slice(articleList, func(i, j int) bool {
		return articleList[i].LeakTime > articleList[j].LeakTime
	})
	return articleList, nil
}

func getArticleByID(id string) (article, error) {
	if articleExists(id) {
		leak := articleCache.get(id, func(string) interface{} {
			articleData, err := readEntry(dataBase, articleCache.path(id))
			if err != nil && debug {
				log.Println(err)
			}

			// was edited
			var edited string
			if editor, ok := articleData["edited_by"]; ok {
				edited = editor.(string)
			} else {
				edited = ""
			}
			var editedWhen int64
			if editTime, ok := articleData["edited_time"]; ok {
				editedWhen = int64(editTime.(float64))
			} else {
				editedWhen = 0
			}

			return article{
				ID:          id,
				Description: articleData["description"].(string),
				Summary:     articleData["summary"].(string),
				LeakTime:    int64(articleData["time"].(float64)),
				ImageUrl:    articleData["image_url"].(string),
				SourceLink:  articleData["source_url"].(string),
				ReporterUid: articleData["reporter_uid"].(string),
				EditedBy:    edited,
				DateEdit:    editedWhen,
				Title:       articleData["title"].(string),
				DateCreate:  int64(articleData["created_time"].(float64)),
			}
		})
		return leak.(article), nil
	}
	return article{}, errors.New("article not found")
}

func articleExists(id string) bool {
	exists := articleCache.has(id)
	if !exists {
		exists = pathExists(dataBase, articleCache.path(id))
	}
	return exists
}

func compileBBCode(in string) string {
	return BBCompiler.Compile(in)
}

func getAllowedLinks() []string {
	items := allowedLinkCache.get("links", func(string) interface{} {
		ref := dataBase.NewRef(allowedLinkLocation)
		var data []string
		if err := ref.Get(ctx, &data); err != nil && debug {
			log.Println(err)
			return nil
		}
		return data
	})
	return items.([]string)
}

func allowedLinksForUserContext(c *gin.Context) []string {
	usr, exists := c.Get("user")
	if exists && usr != nil {
		if usr.(user).AuthLevel < linkLessAuthLevel {
			return getAllowedLinks()
		}
	}
	r := []string{"^(https?|ftp):\\/\\/[^\\s/$.?#].[^\\s]*$"}
	return r
}

func createNewLeak(title, description, rawTime, imageUrl, sourceUrl string, reporter user) (article, error) {
	leak, code := leakSanitization(title, description, rawTime, imageUrl, sourceUrl,
		reporter, user{UID: ""}, time.Now().Unix(), 0)

	switch code {
	case 1:
		addLog(2, reporter.UID, "Unauthorised to Skip Source Link", map[string]interface{}{"leak_metadata": leak})
		return article{}, errors.New("missing source url")
	case 2:
		addLog(2, reporter.UID, "Tried to Post an Invalid Link", map[string]interface{}{"leak_metadata": leak})
		return article{}, errors.New("invalid url")
	case 3:
		addLog(2, reporter.UID, "No leak Body", map[string]interface{}{"leak_metadata": leak})
		return article{}, errors.New("missing body")
	case 4:
		return article{}, errors.New("invalid time")
	}

	key, err := pushEntry(dataBase, articleLocation, leak)
	leak.ID = key
	if title == "" {
		leak.Title = fmt.Sprintf("DecaLeak %d", hashTo32(key))
		err2 := setEntry(dataBase, fmt.Sprintf("%s/title", articleCache.path(leak.ID)), leak.Title)

		if err2 != nil {
			err = err2
		}
	}

	if err != nil {
		addLog(2, reporter.UID, "Failed to Create Leak", map[string]interface{}{"article": leak.ID,
			"leak_metadata": leak})
		return article{}, err
	}

	addLog(2, reporter.UID, "Created Leak", map[string]interface{}{"article": leak.ID})

	clearData := cacheAction{
		CacheListId: articleListCache.id,
		ItemId:      "articles",
		ActionType:  clearList,
	}
	sendAction(clearData)

	return leak, nil
}

/*
cases:
	0 - success
	1 - missing source
	2 - invalid source
	3 - no body
	4 - invalid time
*/
func leakSanitization(title, description, rawTime, imageUrl, sourceUrl string, reporter, editedBy user,
	created, edited int64) (article, int) {
	leakTime, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		return article{}, 4
	}

	sourceUrl = strings.ReplaceAll(sourceUrl, "javascript:", "")

	description = strings.ReplaceAll(description, "\r\n", "\n")
	description = strings.Trim(description, "\n ")

	summery := BBCompiler.Compile(description)
	summery = strings.ReplaceAll(summery, "<br>", "\n")
	summery = stripHtmlRegex(summery)
	summery = cleanRepeatedEnter(cleanRepeatedSpace(summery))
	summery = clip(strings.Trim(summery, "\n "), 200)

	title = strings.Trim(title, " ")
	title = clip(title, 60)

	leak := article{
		Description: description,
		Summary:     summery,
		LeakTime:    leakTime,
		ImageUrl:    imageUrl,
		SourceLink:  sourceUrl,
		ReporterUid: reporter.UID,
		EditedBy:    editedBy.UID,
		Title:       title,
		DateCreate:  created,
		DateEdit:    edited,
	}

	checkPerms := reporter
	if editedBy.UID != "" {
		checkPerms = editedBy
	}
	if checkPerms.AuthLevel < linkLessAuthLevel && sourceUrl == "" {
		return article{}, 1
	}

	if _, err := url.ParseRequestURI(sourceUrl); err != nil && checkPerms.AuthLevel < linkLessAuthLevel {
		return article{}, 2
	}

	if checkPerms.AuthLevel < linkLessAuthLevel {
		var regex *regexp.Regexp
		valid := false
		for _, s := range getAllowedLinks() {
			regex = regexp.MustCompile(s)
			if regex.MatchString(sourceUrl) {
				valid = true
			}
		}
		if !valid {
			return article{}, 2
		}
	}

	if description == "" {
		return article{}, 3
	}
	return leak, 0
}

func createArticle(c *gin.Context) {
	description := c.PostForm("description")
	leakTime := c.PostForm("time")
	imageUrl := c.PostForm("image_url")
	sourceUrl := c.PostForm("source_url")
	//reporter := getUser(c.PostForm("reporter_uid"))
	title := c.PostForm("title")
	reporterUser, _ := c.Get("user")
	reporter := reporterUser.(user)

	if a, err := createNewLeak(title, description, leakTime, imageUrl, sourceUrl, reporter); err == nil {
		// success
		leakLocation := url.URL{
			Scheme: domainBase.Scheme,
			Host:   domainBase.Host,
			Path:   fmt.Sprintf("/leaks/leak/%s", a.ID),
		}
		render(c, gin.H{"status": "success",
			"payload": map[string]interface{}{
				"leakTitle":     a.Title,
				"leakUrl":       leakLocation.String(), //c.Request.URL.Scheme, c.Request.URL.Host,
				"leak":          a,
				"allowed_links": allowedLinksForUserContext(c),
			}, "publishSuccess": true, "linkLessAuthLevel": linkLessAuthLevel},
			"Create new",
			"Share a new DecaLeak",
			"",
			c.Request.URL,
			"postLeak.html", http.StatusCreated)
	} else {
		// error
		if debug {
			log.Println(err)
		}
		render(c, gin.H{"status": "error",
			"payload": map[string]interface{}{
				"description":   description,
				"time":          leakTime,
				"image_url":     imageUrl,
				"source_url":    sourceUrl,
				"reporter_uid":  reporter.UID,
				"title":         title,
				"allowed_links": allowedLinksForUserContext(c),
				"error":         err,
			}, "errorPublishing": true, "linkLessAuthLevel": linkLessAuthLevel},
			"Create new",
			"Share a new DecaLeak",
			"",
			c.Request.URL,
			"postLeak.html")
	}
}
