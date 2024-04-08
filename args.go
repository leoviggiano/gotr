package gotr

import (
	"fmt"
	"strings"
)

type Args struct {
	Identifier string         // The identifier registered in the translator
	Localizer  string         // JSON path or text
	Args       map[string]any // Arguments to be replaced in the template
	Count      int            // Count of the item if applies
}

func (t Args) apply(originalString string) string {
	for k, v := range t.Args {
		originalString = strings.ReplaceAll(originalString, fmt.Sprintf("{{.%s}}", k), fmt.Sprintf("%v", v))
	}

	return originalString
}
