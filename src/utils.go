package main

import (
	"fmt"
	"github.com/frustra/bbcode"
	"github.com/gin-gonic/gin"
	"html/template"
	"regexp"
	"strings"
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

func initBBCode(compiler *bbcode.Compiler) {
	// sanitise the url tag from JS
	compiler.SetTag("url", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "a"
		value := node.GetOpeningTag().Value
		if value == "" {
			text := bbcode.CompileText(node)
			if len(text) > 0 {
				text = strings.ReplaceAll(text, "javascript:", "")
				out.Attrs["href"] = bbcode.ValidURL(text)
			}
		} else {
			value = strings.ReplaceAll(value, "javascript:", "")
			out.Attrs["href"] = bbcode.ValidURL(value)
		}
		return out, true
	})

	// youtube tag
	compiler.SetTag("youtube", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		iframe := bbcode.NewHTMLTag("")
		iframe.Name = "iframe"
		iframe.Attrs["frameborder"] = "0"
		text := bbcode.CompileText(node)
		iframe.Attrs["src"] = fmt.Sprintf("https://www.youtube.com/embed/%s", text)
		iframe.Attrs["allowfullscreen"] = ""
		iframe.AppendChild(nil)
		iframe.Attrs["style"] = "height:max(225px,100%); width:min(100%,400px);"

		out := bbcode.NewHTMLTag("")
		out.Name = "figure"
		out.Attrs["class"] = "image is-disable-16by9"
		out.AppendChild(iframe)

		node.Children = make([]*bbcode.BBCodeNode, 0)
		return out, true
	})
}
