package builder

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type routeGroup struct {
	packagesMap        map[string]importDef
	level              string
	levelServer        string
	levelServerHandler string
	pkg                string
	packagePath        string
	leaf               *AST
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
	Name            string
	Type            string
	Alias           string
	TypeDeclaration string
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
	HTTPCode           string
}

func writeServerFile(level string, leaf *AST, path *routesPathSpec) error {
	rg := &routeGroup{
		packagesMap: map[string]importDef{},
		level:       level,
		pkg:         path.PackageName,
		packagePath: path.PackagePath,
		leaf:        leaf,
	}
	addPackageToMap("net/http", rg.packagesMap, 0)
	addPackageToMap("github.com/gin-gonic/gin", rg.packagesMap, 0)
	filepath := path.PackageName + level + "_server.go"
	hnd, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer hnd.Close()
	return rg.writeServerContent(hnd)
}

func (rg *routeGroup) writeServerContent(hnd *os.File) error {
	tpl, err := template.ParseFiles("tpl/routes.tpl")
	if err != nil {
		return err
	}
	rg.levelServer = Title(rg.level) + "Server"
	rg.levelServerHandler = rg.level + "ServerHandler"
	content, err := rg.buildServerContent()
	if err != nil {
		return err
	}
	return tpl.Execute(hnd, content)
}

func (rg *routeGroup) buildServerContent() (*ServerContent, error) {
	importsArr := []importDef{}
	methods, functions, err := rg.buildServerMethods()
	if err != nil {
		return nil, err
	}
	for _, pkg := range rg.packagesMap {
		importsArr = append(importsArr, pkg)
	}
	imports, err := imports(importsArr...)
	if err != nil {
		return nil, err
	}
	cnt := &ServerContent{
		Package:            rg.pkg,
		Imports:            imports,
		LevelServer:        rg.levelServer,
		LevelServerHandler: rg.levelServerHandler,
		Methods:            methods,
		Functions:          functions,
	}
	return cnt, nil
}

func (rg *routeGroup) buildServerMethods() (methods, functions []string, err error) {
	i := 0
	methods = make([]string, len(rg.leaf.Node.Methods))
	functions = make([]string, len(rg.leaf.Node.Methods))
	for _, method := range rg.leaf.Node.Methods {
		methods[i], err = rg.buildServerMethod(method, rg.leaf.Node.URL)
		if err != nil {
			return nil, nil, err
		}
		functions[i], err = rg.buildServerFunction(method, rg.leaf.Node.URL)
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
	if methodDef.Definition.RequestParam != nil {
		requestAlias = addPackageToMap(methodDef.Definition.RequestParam.Path, rg.packagesMap, 0)
		addPackageToMap("fmt", rg.packagesMap, 0)
		addPackageToMap("gopkg.in/go-playground/validator.v9", rg.packagesMap, 0)
	}
	if methodDef.Definition.ResponseParam != nil {
		responseAlias = addPackageToMap(methodDef.Definition.ResponseParam.Path, rg.packagesMap, 0)
	}
	ctn := &MethodContent{
		Method:     methodDef.Handler,
		Params:     rg.getRequestProto(methodDef, requestAlias),
		ReturnType: rg.getReturnType(methodDef, responseAlias),
	}
	return ctn
}

func (rg *routeGroup) getRequestProto(methodDef *RouteDef, alias string) string {
	params := []string{"g *gin.Context"}
	if methodDef.Param != "" {
		params = append(params, fmt.Sprintf("%s string", methodDef.Param))
	}
	if methodDef.Definition.Request != nil {
		params = append(params, methodDef.Definition.RequestParam.getObjectDeclaration())
	}
	return strings.Join(params, ", ")
}

func (rg *routeGroup) getReturnType(methodDef *RouteDef, alias string) string {
	pre, post := "", ""
	params := make([]string, 0)
	if methodDef.Definition.ResponseParam != nil {
		pre, post = "(", ")"
		params = append(params, methodDef.Definition.ResponseParam.getTypeDeclaration())
	}
	params = append(params, "error")
	return fmt.Sprintf("%s%s%s", pre, strings.Join(params, ", "), post)
}

func (rg *routeGroup) buildServerFunction(methodDef *RouteDef, url string) (string, error) {
	tpl, err := template.ParseFiles("tpl/route.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, rg.getServerFunction(methodDef, url)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func (rg *routeGroup) getServerFunction(methodDef *RouteDef, url string) *HandlerDefContent {
	return &HandlerDefContent{
		Handler:            methodDef.Handler,
		Param:              methodDef.Param,
		Request:            rg.parseRequestObject(methodDef.Definition),
		LevelServerHandler: rg.levelServerHandler,
		Method:             methodDef.Handler,
		Returns:            rg.getResponseParams(methodDef.Definition),
		Params:             rg.getRequestParams(methodDef),
		Response:           rg.parseResponseObject(methodDef.Definition),
		HTTPCode:           getHTTPCodeByMethod(methodDef.Method),
	}
}

func (rg *routeGroup) parseRequestObject(def *R) *RequestContent {
	if def.RequestParam == nil {
		return nil
	}
	return &RequestContent{
		Name:            def.RequestParam.Name,
		Type:            def.RequestParam.Type,
		Alias:           addPackageToMap(def.RequestParam.Path, rg.packagesMap, 0),
		TypeDeclaration: def.RequestParam.getTypeDeclaration(),
	}
}

func (rg *routeGroup) parseResponseObject(def *R) *RequestContent {
	if def.ResponseParam == nil {
		return nil
	}
	return &RequestContent{
		Name:            def.ResponseParam.Name,
		Type:            def.ResponseParam.Type,
		Alias:           addPackageToMap(def.ResponseParam.Path, rg.packagesMap, 0),
		TypeDeclaration: def.ResponseParam.getTypeDeclaration(),
	}
}

func (rg *routeGroup) getRequestParams(methodDef *RouteDef) string {
	params := []string{"g"}
	if methodDef.Param != "" {
		params = append(params, methodDef.Param)
	}
	if methodDef.Definition.Request != nil {
		params = append(params, methodDef.Definition.RequestParam.Name)
	}
	return strings.Join(params, ", ")
}

func (rg *routeGroup) getResponseParams(def *R) string {
	params := make([]string, 0)
	if def.ResponseParam != nil {
		params = append(params, def.ResponseParam.Name)
	}
	params = append(params, "err")
	return strings.Join(params, ", ")
}

func getHTTPCodeByMethod(method string) string {
	switch method {
	case "DELETE":
		return "204"
	case "POST":
		return "201"
	}
	return "200"
}
