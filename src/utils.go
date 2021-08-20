package main

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"regexp"
	"strings"

	"github.com/frustra/bbcode"
	"github.com/gin-gonic/gin"
)

var tagRegex = regexp.MustCompile(`<.*?>`)
var bbExclude = regexp.MustCompile(`<bbexclude>.*?</bbexclude>`)

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
	s = bbExclude.ReplaceAllString(s, "")
	return tagRegex.ReplaceAllString(s, "")
}

var youtubeRE = regexp.MustCompile(`[a-zA-Z\-_0-9]{11}`)

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
		var youtubeID string
		if youtubeRE.MatchString(text) {
			youtubeID = youtubeRE.FindString(text)
		} else {
			youtubeID = text
		}
		iframe.Attrs["src"] = fmt.Sprintf("https://www.youtube.com/embed/%s", youtubeID)
		iframe.Attrs["allowfullscreen"] = ""
		iframe.AppendChild(nil)
		// todo this needs improving :/
		iframe.Attrs["style"] = "height:max(225px,100%); width:min(100%,400px);"

		out := bbcode.NewHTMLTag("")
		out.Name = "figure"
		// here too
		out.Attrs["class"] = "image is-disable-16by9"
		out.AppendChild(iframe)

		node.Children = make([]*bbcode.BBCodeNode, 0)
		return out, true
	})

	// video tag
	compiler.SetTag("video", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		videoFrame := bbcode.NewHTMLTag("")
		videoFrame.Name = "video"
		//videoFrame.Attrs["style"] = "height:max(225px,100%); width:min(100%,400px);"
		videoFrame.Attrs["controls"] = "true"
		value := node.GetOpeningTag().Value
		var src string
		if value == "" {
			text := bbcode.CompileText(node)
			if len(text) > 0 {
				src = strings.ReplaceAll(text, "javascript:", "")
			}
		} else {
			src = strings.ReplaceAll(value, "javascript:", "")
		}
		source := bbcode.NewHTMLTag("")
		source.Name = "source"
		source.Attrs["src"] = src
		source.Attrs["type"] = "video/mp4"
		videoFrame.AppendChild(source)

		exclude := bbcode.NewHTMLTag("")
		exclude.Name = "bbexclude"

		warn := bbcode.NewHTMLTag("")
		warn.Name = "strong"
		warn.Attrs["class"] = "has-text-danger is-underlined"
		warn.AppendChild(bbcode.NewHTMLTag("Your browser does not support the video tag."))
		exclude.AppendChild(warn)
		videoFrame.AppendChild(exclude)

		out := bbcode.NewHTMLTag("")
		out.Name = "figure"

		out.Attrs["class"] = "image is-disable-16by9"
		out.AppendChild(videoFrame)

		node.Children = make([]*bbcode.BBCodeNode, 0)
		return out, true
	})
}

var _repeatedEnter = regexp.MustCompile(`\n\n+`)
var _repeatedSpace = regexp.MustCompile(`  +`)

func cleanRepeatedEnter(input string) string {
	return _repeatedEnter.ReplaceAllString(input, "\n")
}
func cleanRepeatedSpace(input string) string {
	return _repeatedSpace.ReplaceAllString(input, ` `)
}

func hashTo32(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
