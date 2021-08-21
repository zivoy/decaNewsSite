package main

import "strings"

// bulmaColors bitmask
var bulmaColors = map[BulmaColor]string{
	Black:   "is-black",
	Dark:    "is-dark",
	Light:   "is-light",
	White:   "is-white",
	Primary: "is-primary",
	Link:    "is-link",
	Info:    "is-info",
	Success: "is-success",
	Warning: "is-warning",
	Danger:  "is-danger",
	Lighter: "is-light",
}

const (
	Lighter = 1 << iota
	Black
	Dark
	Light
	White
	Primary
	Link
	Info
	Success
	Warning
	Danger
)

type BulmaColor int

func getColorClassString(mask BulmaColor) string {
	if mask == 0 {
		return ""
	}
	result := make([]string, 0)
	for k, v := range bulmaColors {
		if (mask & k) != 0 {
			result = append(result, v)
		}
	}
	return strings.Join(result, " ")
}
