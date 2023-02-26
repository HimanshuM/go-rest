package builder

import (
	"fmt"
	"strings"
)

var pkgPath string

func Package(pkg string) {
	pkgPath = strings.Trim(pkg, "/")
}

func Path(route string) (*AST, error) {
	return getPath(route, root)
}

func (a *AST) Path(route string) (*AST, error) {
	return getPath(route, a)
}

func (a *AST) append(route string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(a.Node.FullURL, "/"), strings.TrimLeft(route, "/"))
}

func (a *AST) GET(def *R) *AST {
	method := "GET"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) POST(def *R) *AST {
	method := "POST"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) PUT(def *R) *AST {
	method := "PUT"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) PATCH(def *R) *AST {
	method := "PATCH"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) DELETE(def *R) *AST {
	method := "DELETE"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) OPTIONS(def *R) *AST {
	method := "OPTIONS"
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
	}
	return a
}

func (a *AST) getHandler(method string) string {
	name := ""
	components := strings.Split(a.Node.URL, "/")
	for _, c := range components {
		c = strings.Replace(c, "{", "", -1)
		c = strings.Replace(c, "}", "", -1)
		name += strings.Title(c)
	}
	return fmt.Sprintf("%s%s", name, strings.Title(strings.ToLower(method)))
}
