package builder

import (
	"reflect"
	"strings"
)

type Parameter struct {
	Type    string
	Name    string
	Package string
	Path    string
}

type R struct {
	Query         interface{}
	Request       interface{}
	Response      interface{}
	RequestParam  *Parameter
	ResponseParam *Parameter
}

type RouteDef struct {
	Method     string
	Handler    string
	Definition *R
	Param      string
}

type Route struct {
	URL     string
	FullURL string
	Methods map[string]*RouteDef
}

type AST struct {
	Level         string
	Node          *Route
	Tree          map[string]*AST
	HasDefinition bool
	Package       string
}

func (r *R) processDefinition() {
	if r.Request != nil {
		r.processRequest()
	}
	if r.Response != nil {
		r.processResponse()
	}
}

func (r *R) processRequest() {
	request := reflect.TypeOf(r.Request)
	if request.Kind() == reflect.Pointer {
		request = request.Elem()
	}
	requestType := request.Name()
	pkgPath := request.PkgPath()
	r.RequestParam = &Parameter{
		Type:    requestType,
		Name:    strings.ToLower(requestType[0:1]) + requestType[1:],
		Package: getLastComponent(pkgPath),
		Path:    pkgPath,
	}
}

func (r *R) processResponse() {
	response := reflect.TypeOf(r.Response)
	if response.Kind() == reflect.Pointer {
		response = response.Elem()
	}
	responseType := response.Name()
	pkgPath := response.PkgPath()
	r.RequestParam = &Parameter{
		Type:    responseType,
		Name:    strings.ToLower(responseType[0:1]) + responseType[1:],
		Package: getLastComponent(pkgPath),
		Path:    pkgPath,
	}
}
