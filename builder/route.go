package builder

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var routesPkgPath, handlersPkgPath string

func RoutesPackage(pkg string) {
	routesPkgPath = strings.Trim(pkg, "/")
}

func HandlersPackage(pkg string) {
	handlersPkgPath = strings.Trim(pkg, "/")
}

func Path(route string, middlewares ...gin.HandlerFunc) (*AST, error) {
	leaf, err := getPath(route, root)
	if err != nil {
		return nil, err
	}
	leaf.Middlewares = append(leaf.Middlewares, middlewares...)
	return leaf, nil
}

func (a *AST) Path(route string, middlewares ...gin.HandlerFunc) (*AST, error) {
	leaf, err := getPath(route, a)
	if err != nil {
		return nil, err
	}
	leaf.Middlewares = append(leaf.Middlewares, middlewares...)
	return leaf, nil
}

func (a *AST) append(route string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(a.Node.FullURL, "/"), strings.TrimLeft(route, "/"))
}

func (a *AST) GET(def *R) *AST {
	method := "GET"
	def.processDefinition()
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
	def.processDefinition()
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
	def.processDefinition()
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
	def.processDefinition()
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
	def.processDefinition()
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
	def.processDefinition()
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
		name += Title(cleanupRoute(c))
	}
	return fmt.Sprintf("%s%s", name, Title(strings.ToLower(method)))
}
