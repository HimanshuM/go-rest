package builder

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type ServerContent struct {
	Package            string
	Imports            string
	LevelServer        string
	Methods            []string
	LevelServerHandler string
	Functions          []string
	Level              string
}

type MethodContent struct {
	Method     string
	Params     string
	ReturnType string
}

type RequestContent struct {
	Name string
	Type string
}

type HandlerDefContent struct {
	Handler            string
	Param              string
	Request            *RequestContent
	LevelServerHandler string
	Method             string
	Params             string
	Returns            string
	Response           *RequestContent
}

func writeServerFile(path, level, pkg, pkgPath string, leaf *AST) error {
	path += level + "_server.go"
	hnd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer hnd.Close()
	return writeServerContent(hnd, level, pkg, pkgPath, leaf)
}

func writeServerContent(hnd *os.File, level, pkg, pkgPath string, leaf *AST) error {
	tpl, err := template.ParseFiles("tpl/routes.tpl")
	if err != nil {
		return err
	}
	levelServer := strings.Title(level) + "Server"
	levelServerHandler := level + "ServerHandler"
	content, err := buildServerContent(leaf, pkg, levelServer, levelServerHandler)
	if err != nil {
		return err
	}
	return tpl.Execute(hnd, content)
}

func buildServerContent(leaf *AST, pkg, levelServer, levelServerHandler string) (*ServerContent, error) {
	importsArr := []string{"github.com/gin-gonic/gin"}
	packages, methods, functions, err := buildServerMethods(leaf, levelServer, levelServerHandler)
	if err != nil {
		return nil, err
	}
	for _, pkg := range packages {
		importsArr = append(importsArr, pkg)
	}
	imports, err := imports(importsArr...)
	if err != nil {
		return nil, err
	}
	cnt := &ServerContent{
		Package:            pkg,
		Imports:            imports,
		LevelServer:        levelServer,
		LevelServerHandler: levelServerHandler,
		Methods:            methods,
		Functions:          functions,
	}
	return cnt, nil
}

func buildServerMethods(leaf *AST, levelServer, levelServerHandler string) (pkgs, methods, functions []string, err error) {
	i := 0
	methods = make([]string, len(leaf.Node.Methods))
	functions = make([]string, len(leaf.Node.Methods))
	for _, method := range leaf.Node.Methods {
		methods[i], err = buildServerMethod(method, leaf.Node.URL)
		if err != nil {
			return nil, nil, nil, err
		}
		functions[i], err = buildServerFunction(method, leaf.Node.URL, levelServer, levelServerHandler)
		if err != nil {
			return nil, nil, nil, err
		}
		i++
	}
	return
}

func buildServerMethod(methodDef *RouteDef, url string) (string, error) {
	tpl, err := template.ParseFiles("tpl/method.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, getServerMethod(methodDef, url)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func getServerMethod(methodDef *RouteDef, url string) *MethodContent {
	return &MethodContent{
		Method:     methodDef.Handler,
		Params:     getRequestProto(methodDef),
		ReturnType: getReturnType(methodDef),
	}
}

func getRequestProto(methodDef *RouteDef) string {
	params := make([]string, 0)
	if methodDef.Param != "" {
		params = append(params, fmt.Sprintf("%s string", methodDef.Param))
	}
	if methodDef.Definition.Request != nil {
		methodDef.Definition.processRequest()
		params = append(params, fmt.Sprintf("%s *%s", methodDef.Definition.requestName, methodDef.Definition.requestType))
	}
	return strings.Join(params, ", ")
}

func getReturnType(methodDef *RouteDef) string {
	pre, post := "", ""
	params := make([]string, 0)
	if methodDef.Definition.Response != nil {
		pre, post = "(", ")"
		methodDef.Definition.processRequest()
		params = append(params, fmt.Sprintf("*%s", methodDef.Definition.responseType))
	}
	params = append(params, "error")
	return fmt.Sprintf("%s%s%s", pre, strings.Join(params, ", "), post)
}

func buildServerFunction(methodDef *RouteDef, url, levelServer, levelServerHandler string) (string, error) {
	tpl, err := template.ParseFiles("tpl/route.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, getServerFunction(methodDef, url, levelServer, levelServerHandler)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func getServerFunction(methodDef *RouteDef, url, levelServer, levelServerHandler string) *HandlerDefContent {
	return &HandlerDefContent{
		Handler:            methodDef.Handler,
		Param:              methodDef.Param,
		Request:            parseObjectType(methodDef.Definition, true),
		LevelServerHandler: levelServerHandler,
		Method:             methodDef.Handler,
		Returns:            getResponseParams(methodDef.Definition),
		Params:             getRequestParams(methodDef),
		Response:           parseObjectType(methodDef.Definition, false),
	}
}

func parseObjectType(def *R, request bool) *RequestContent {
	if request {
		if def.Request == nil {
			return nil
		}
		return &RequestContent{
			Name: def.requestName,
			Type: def.requestType,
		}
	} else {
		if def.Response == nil {
			return nil
		}
		return &RequestContent{
			Name: def.responseName,
			Type: def.responseType,
		}
	}
}

func getRequestParams(methodDef *RouteDef) string {
	params := make([]string, 0)
	if methodDef.Param != "" {
		params = append(params, methodDef.Param)
	}
	if methodDef.Definition.Request != nil {
		params = append(params, methodDef.Definition.requestName)
	}
	return strings.Join(params, ", ")
}

func getResponseParams(def *R) string {
	params := make([]string, 0)
	if def.Response != nil {
		def.processRequest()
		params = append(params, def.responseName)
	}
	params = append(params, "err")
	return strings.Join(params, ", ")
}
