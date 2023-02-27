package builder

import (
	"reflect"
	"strings"
)

type R struct {
	Query            interface{}
	Request          interface{}
	Response         interface{}
	requestType      string
	requestName      string
	requestTypeName  string
	requestTypePath  string
	responseType     string
	responseName     string
	responseTypeName string
	responseTypePath string
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
	r.requestType = request.Name()
	r.requestName = strings.ToLower(r.requestType[0:1]) + r.requestType[1:]
	r.requestTypePath = request.PkgPath()
	r.requestTypeName = getLastComponent(r.requestTypePath)
}

func (r *R) processResponse() {
	response := reflect.TypeOf(r.Response)
	if response.Kind() == reflect.Pointer {
		response = response.Elem()
	}
	r.responseType = response.Name()
	r.responseName = strings.ToLower(r.responseType[0:1]) + r.responseType[1:]
	r.responseTypePath = response.PkgPath()
	r.responseTypeName = getLastComponent(r.responseTypePath)
}
