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
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) POST(def *R) *AST {
	method := "POST"
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) PUT(def *R) *AST {
	method := "PUT"
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) PATCH(def *R) *AST {
	method := "PATCH"
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) DELETE(def *R) *AST {
	method := "DELETE"
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) OPTIONS(def *R) *AST {
	method := "OPTIONS"
	def.processRequest()
	a.Node.Methods[method] = &RouteDef{
		Method:     method,
		Handler:    a.getHandler(method),
		Definition: def,
		Param:      getURLParam(a.Node.URL),
	}
	a.HasDefinition = true
	return a
}

func (a *AST) getHandler(method string) string {
	name := ""
	components := strings.Split(a.Node.URL, "/")
	for _, c := range components {
		name += strings.Title(cleanupRoute(c))
	}
	return fmt.Sprintf("%s%s", name, strings.Title(strings.ToLower(method)))
}
