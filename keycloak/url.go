package keycloak

import (
	"fmt"
	"net/url"
	"strings"
)

// Transforms a map of string<>string entries to a url-encoded query string.
// The result does not start with a question mark.
func urlEncode(params map[string]string) string {
	parts := make([]string, len(params))

	i := 0
	for k, v := range params {
		parts[i] = fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		i++
	}

	return strings.Join(parts, "&")
}
