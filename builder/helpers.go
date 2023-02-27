package builder

import "strings"

func cleanupRoute(route string) string {
	route = strings.Trim(route, "/")
	route = strings.ReplaceAll(route, "{", "")
	return strings.ReplaceAll(route, "}", "")
}

func getURLParam(url string) string {
	opening := strings.Index(url, "{")
	if opening < 0 {
		return ""
	}
	closing := strings.Index(url, "}")
	return url[opening+1 : closing]
}
