package builder

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type routeGroup struct {
	packagesMap map[string]string
}

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
	Name  string
	Type  string
	Alias string
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
	rg := &routeGroup{
		packagesMap: map[string]string{},
	}
	path += level + "_server.go"
	hnd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer hnd.Close()
	return rg.writeServerContent(hnd, level, pkg, pkgPath, leaf)
}

func (rg *routeGroup) writeServerContent(hnd *os.File, level, pkg, pkgPath string, leaf *AST) error {
	tpl, err := template.ParseFiles("tpl/routes.tpl")
	if err != nil {
		return err
	}
	levelServer := Title(level) + "Server"
	levelServerHandler := level + "ServerHandler"
	content, err := rg.buildServerContent(leaf, pkg, levelServer, levelServerHandler)
	if err != nil {
		return err
	}
	return tpl.Execute(hnd, content)
}

func (rg *routeGroup) buildServerContent(leaf *AST, pkg, levelServer, levelServerHandler string) (*ServerContent, error) {
	importsArr := []string{"github.com/gin-gonic/gin"}
	methods, functions, err := rg.buildServerMethods(leaf, levelServer, levelServerHandler)
	if err != nil {
		return nil, err
	}
	for alias, pkg := range rg.packagesMap {
		pkgName := getLastComponent(pkg)
		if pkgName == alias {
			importsArr = append(importsArr, pkg)
		} else {
			importsArr = append(importsArr, fmt.Sprintf("%s %s", alias, pkg))
		}
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

func (rg *routeGroup) buildServerMethods(leaf *AST, levelServer, levelServerHandler string) (methods, functions []string, err error) {
	i := 0
	methods = make([]string, len(leaf.Node.Methods))
	functions = make([]string, len(leaf.Node.Methods))
	for _, method := range leaf.Node.Methods {
		methods[i], err = rg.buildServerMethod(method, leaf.Node.URL)
		if err != nil {
			return nil, nil, err
		}
		functions[i], err = rg.buildServerFunction(method, leaf.Node.URL, levelServer, levelServerHandler)
		if err != nil {
			return nil, nil, err
		}
		i++
	}
	return
}

func (rg *routeGroup) buildServerMethod(methodDef *RouteDef, url string) (string, error) {
	tpl, err := template.ParseFiles("tpl/method.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, rg.getServerMethod(methodDef, url)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func (rg *routeGroup) getServerMethod(methodDef *RouteDef, url string) *MethodContent {
	requestAlias, responseAlias := "", ""
	if methodDef.Definition.requestTypePath != "" {
		requestAlias = addPackageToMap(methodDef.Definition.requestTypePath, rg.packagesMap, 0)
	}
	if methodDef.Definition.responseTypePath != "" {
		responseAlias = addPackageToMap(methodDef.Definition.responseTypePath, rg.packagesMap, 0)
	}
	ctn := &MethodContent{
		Method:     methodDef.Handler,
		Params:     rg.getRequestProto(methodDef, requestAlias),
		ReturnType: rg.getReturnType(methodDef, responseAlias),
	}
	return ctn
}

func (rg *routeGroup) getRequestProto(methodDef *RouteDef, alias string) string {
	params := make([]string, 0)
	if methodDef.Param != "" {
		params = append(params, fmt.Sprintf("%s string", methodDef.Param))
	}
	if methodDef.Definition.Request != nil {
		params = append(params, fmt.Sprintf("%s *%s.%s", methodDef.Definition.requestName, alias, methodDef.Definition.requestType))
	}
	return strings.Join(params, ", ")
}

func (rg *routeGroup) getReturnType(methodDef *RouteDef, alias string) string {
	pre, post := "", ""
	params := make([]string, 0)
	if methodDef.Definition.Response != nil {
		pre, post = "(", ")"
		params = append(params, fmt.Sprintf("*%s.%s", alias, methodDef.Definition.responseType))
	}
	params = append(params, "error")
	return fmt.Sprintf("%s%s%s", pre, strings.Join(params, ", "), post)
}

func (rg *routeGroup) buildServerFunction(methodDef *RouteDef, url, levelServer, levelServerHandler string) (string, error) {
	tpl, err := template.ParseFiles("tpl/route.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, rg.getServerFunction(methodDef, url, levelServer, levelServerHandler)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func (rg *routeGroup) getServerFunction(methodDef *RouteDef, url, levelServer, levelServerHandler string) *HandlerDefContent {
	return &HandlerDefContent{
		Handler:            methodDef.Handler,
		Param:              methodDef.Param,
		Request:            rg.parseObjectType(methodDef.Definition, true),
		LevelServerHandler: levelServerHandler,
		Method:             methodDef.Handler,
		Returns:            rg.getResponseParams(methodDef.Definition),
		Params:             rg.getRequestParams(methodDef),
		Response:           rg.parseObjectType(methodDef.Definition, false),
	}
}

func (rg *routeGroup) parseObjectType(def *R, request bool) *RequestContent {
	if request {
		if def.Request == nil {
			return nil
		}
		return &RequestContent{
			Name:  def.requestName,
			Type:  def.requestType,
			Alias: addPackageToMap(def.requestTypePath, rg.packagesMap, 0),
		}
	} else {
		if def.Response == nil {
			return nil
		}
		return &RequestContent{
			Name:  def.responseName,
			Type:  def.responseType,
			Alias: addPackageToMap(def.responseTypePath, rg.packagesMap, 0),
		}
	}
}

func (rg *routeGroup) getRequestParams(methodDef *RouteDef) string {
	params := make([]string, 0)
	if methodDef.Param != "" {
		params = append(params, methodDef.Param)
	}
	if methodDef.Definition.Request != nil {
		params = append(params, methodDef.Definition.requestName)
	}
	return strings.Join(params, ", ")
}

func (rg *routeGroup) getResponseParams(def *R) string {
	params := make([]string, 0)
	if def.Response != nil {
		params = append(params, def.responseName)
	}
	params = append(params, "err")
	return strings.Join(params, ", ")
}
