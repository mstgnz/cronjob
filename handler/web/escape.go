package web

import (
	"fmt"
	"html"
)

// esc renders a value for inclusion in the HTML fragments these handlers build by
// hand. The htmx partials are written straight to the response, so they bypass the
// auto escaping that html/template would normally provide, and any user supplied
// text (urls, names, titles, header values, marshalled json in data-* attributes)
// would otherwise be interpreted as markup.
func esc(value any) string {
	return html.EscapeString(fmt.Sprintf("%v", value))
}

// escBytes is the []byte form of esc, used for json.Marshal output embedded in attributes.
func escBytes(value []byte) string {
	return html.EscapeString(string(value))
}
