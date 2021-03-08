package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"regexp"
)

const tagString = `<.*?>`

var tagRegex = regexp.MustCompile(tagString)

func clip(input string, maxLength int) string {
	if len(input) > maxLength {
		s := []rune(input)
		input = string(s[:maxLength-1]) + "…"
	}
	return input
}

func deleteCookie(c *gin.Context, name string) {
	setCookie(c, name, "delat", -1)
}

func setCookie(c *gin.Context, name string, value string, maxAge int, secure ...bool) {
	secc := store.Options.Secure
	if len(secure) > 0 {
		secc = secure[0]
	}
	c.SetCookie(name, value, maxAge, store.Options.Path, store.Options.Domain, secc, store.Options.HttpOnly)
}

func unescape(s string) template.HTML {
	return template.HTML(s)
}

func stripHtmlRegex(s string) string {
	return tagRegex.ReplaceAllString(s, "")
}
