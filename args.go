package gotr

import (
	"fmt"
	"strings"
)

type Args struct {
	Identifier string         // The identifier registered in the translator
	Localizer  string         // JSON path or text
	Vars       map[string]any // Arguments to be replaced in the template
	Count      *int           // Count for pluralization; nil means singular, 0 means none
}

// CountOf returns a pointer to n for use in Args.Count.
func CountOf(n int) *int {
	return &n
}

func (t Args) apply(originalString string) string {
	for k, v := range t.Vars {
		originalString = strings.ReplaceAll(originalString, fmt.Sprintf("{{.%s}}", k), fmt.Sprintf("%v", v))
	}

	return originalString
}
