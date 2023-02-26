package builder

import "strings"

func cleanupRoute(route string) string {
	route = strings.Trim(route, "/")
	route = strings.ReplaceAll(route, "{", "")
	return strings.ReplaceAll(route, "}", "")
}
