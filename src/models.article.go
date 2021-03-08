package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type article struct {
	ID          string `json:"-"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	LeakTime    int64  `json:"time"`
	ImageUrl    string `json:"image_url"`
	DiscordLink string `json:"source_url"`
	ReporterUid string `json:"reporter_uid"`
}

func articlePathString(uid string) string {
	return fmt.Sprintf("leaks/%s", uid)
}

// this function is really bad todo make this be done on the front end
// todo implement splitting of pages
func getAllArticles(page int) ([]article, error) {
	ref := dataBase.NewRef("leaks")
	var data map[string]article
	if err := ref.Get(ctx, &data); err != nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}

	articleList := make([]article, 0)
	for k, v := range data {
		v.ID = k
		articleList = append(articleList, v)
		addCache(articleCache, k, v)
	}

	sort.Slice(articleList, func(i, j int) bool {
		return articleList[i].LeakTime > articleList[j].LeakTime
	})
	return articleList, nil
}

func getAllUsersArticles(uid string) ([]article, error) {
	ref := dataBase.NewRef("leaks")
	var data map[string]article
	if err := ref.OrderByChild("reporter_uid").EqualTo(uid).Get(ctx, &data); err != nil {
		return nil, fmt.Errorf("error reading from database: %v", err)
	}

	articleList := make([]article, 0)
	for k, v := range data {
		v.ID = k
		articleList = append(articleList, v)
		addCache(articleCache, k, v)
	}

	sort.Slice(articleList, func(i, j int) bool {
		return articleList[i].LeakTime > articleList[j].LeakTime
	})
	return articleList, nil
}

func getArticleByID(id string) (article, error) {
	if articleExists(id) {
		leak := getCache(articleCache, id, func(string) interface{} {
			articleData, err := readEntry(dataBase, articlePathString(id))
			if err != nil && debug {
				fmt.Println(err)
			}
			return article{
				ID:          id,
				Description: articleData["description"].(string),
				Summary:     articleData["summary"].(string),
				LeakTime:    int64(articleData["time"].(float64)),
				ImageUrl:    articleData["image_url"].(string),
				DiscordLink: articleData["source_url"].(string),
				ReporterUid: articleData["reporter_uid"].(string),
			}
		})
		return leak.(article), nil
	}
	return article{}, errors.New("article not found")
}

func articleExists(id string) bool {
	_, exists := articleCache[id]
	if !exists {
		exists = pathExists(dataBase, articlePathString(id))
	}
	return exists
}

func getAllowedLink() []string {
	ref := dataBase.NewRef("admin/allowed_links")
	var data []string
	if err := ref.Get(ctx, &data); err != nil && debug {
		fmt.Println(err)
		return nil
	}
	return data
}

func createNewLeak(description string, rawTime string, imageUrl string, discordUrl string, reporterUid string) (article, error) {
	time, err := strconv.Atoi(rawTime)
	if err != nil {
		return article{}, errors.New("invalid time")
	}

	leak := article{
		Description: description,
		Summary:     clip(description, 200),
		LeakTime:    int64(time),
		ImageUrl:    imageUrl,
		DiscordLink: discordUrl,
		ReporterUid: reporterUid,
	}

	key, err := pushEntry(dataBase, "leaks", leak)
	if err != nil {
		return article{}, err
	}

	leak.ID = key
	return leak, nil
}

func createArticle(c *gin.Context) {
	description := c.PostForm("description")
	time := c.PostForm("time")
	imageUrl := c.PostForm("image_url")
	discordUrl := c.PostForm("discord_url")
	reporterUid := c.PostForm("reporter_uid")

	if a, err := createNewLeak(description, time, imageUrl, discordUrl, reporterUid); err == nil {
		// success
		leakLocation := url.URL{
			Scheme: domainBase.Scheme,
			Host:   domainBase.Host,
			Path:   fmt.Sprintf("/leaks/leak/%s", a.ID),
		}
		render(c, gin.H{"status": "success",
			"payload": map[string]interface{}{
				"leakId":        a.ID,
				"leakUrl":       leakLocation.String(), //c.Request.URL.Scheme, c.Request.URL.Host,
				"leak":          a,
				"allowed_links": getAllowedLink(),
			}, "publishSuccess": true},
			"Create new",
			"Share a new DecaLeak",
			" ",
			c.Request.URL,
			"postLeak.html", http.StatusCreated)
	} else {
		// error
		render(c, gin.H{"status": "error",
			"payload": map[string]interface{}{
				"description":   description,
				"time":          time,
				"image_url":     imageUrl,
				"discord_url":   discordUrl,
				"reporter_uid":  reporterUid,
				"allowed_links": getAllowedLink(),
			}, "errorPublishing": true},
			"Create new",
			"Share a new DecaLeak",
			" ",
			c.Request.URL,
			"postLeak.html")
	}
}
