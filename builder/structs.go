package builder

import (
	"fmt"
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

func (r *R) processRequest() {
	if r.requestName != "" || r.responseName != "" {
		return
	}

	if r.Request != nil {
		request := reflect.TypeOf(r.Request)
		if request.Kind() == reflect.Pointer {
			request = request.Elem()
		}
		r.requestType = request.Name()
		r.requestName = strings.ToLower(r.requestType[0:1]) + r.requestType[1:]
		r.requestTypePath = request.PkgPath()
		fmt.Println(r.requestTypePath)
	}

	if r.Response != nil {
		response := reflect.TypeOf(r.Response)
		if response.Kind() == reflect.Pointer {
			response = response.Elem()
		}
		r.responseType = response.Name()
		r.responseName = strings.ToLower(r.responseType[0:1]) + r.responseType[1:]
		r.responseTypePath = response.PkgPath()
		fmt.Println(r.responseTypePath)
	}
}
