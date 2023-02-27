package builder

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

func Title(str string) string {
	return cases.Title(language.English).String(str)
}

func getLastComponent(str string) string {
	components := strings.Split(str, "/")
	return components[len(components)-1]
}
