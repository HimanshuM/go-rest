package builder

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type handlerGroup struct {
	packagesMap  map[string]importDef
	level        string
	levelServer  string
	pkg          string
	packagePath  string
	leaf         *AST
	routePkgName string
	routePkgPath string
}

type HandlersContent struct {
	Package        string
	Imports        string
	LevelServer    string
	HasLevelServer bool
	Methods        []string
	Level          string
	RoutesPackage  string
	InitName       string
	SubPaths       []string
}

type HandlerContent struct {
	LevelServer string
	Handler     string
	Params      string
	ReturnType  string
	Param       string
	Response    *RequestContent
}

func writeHandlerFile(level string, leaf *AST, path *PathSpec) error {
	hg := &handlerGroup{
		packagesMap:  map[string]importDef{},
		level:        level,
		pkg:          path.HandlerPackageName,
		packagePath:  path.HandlerPackagePath,
		leaf:         leaf,
		routePkgPath: path.RoutePackagePath,
		routePkgName: path.RoutePackageName,
	}
	if leaf.HasDefinition {
		addPackageToMap("github.com/gin-gonic/gin", hg.packagesMap, 0)
		addPackageToMap(hg.routePkgPath, hg.packagesMap, 0)
	}
	filepath := path.HandlerFilePath + ".go"
	if level == "" {
		filepath = fmt.Sprintf("%s/%s", path.HandlerFilePath, filepath)
	}
	hnd, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer hnd.Close()
	return hg.writeHandlerContent(hnd)
}

func (hg *handlerGroup) writeHandlerContent(hnd *os.File) error {
	tpl, err := template.ParseFiles("tpl/handlers.tpl")
	if err != nil {
		return err
	}
	hg.levelServer = Title(hg.level) + "Server"
	content, err := hg.buildHandlersContent()
	if err != nil {
		return err
	}
	return tpl.Execute(hnd, content)
}

func (hg *handlerGroup) buildHandlersContent() (*HandlersContent, error) {
	importsArr := []importDef{}
	methods, err := hg.buildHandlerMethods()
	if err != nil {
		return nil, err
	}
	initName, subPaths := hg.buildInitInvokes()
	for _, pkg := range hg.packagesMap {
		importsArr = append(importsArr, pkg)
	}
	imports, err := imports(importsArr...)
	if err != nil {
		return nil, err
	}
	cnt := &HandlersContent{
		Package:        hg.pkg,
		Imports:        imports,
		LevelServer:    hg.levelServer,
		HasLevelServer: hg.leaf.HasDefinition,
		Methods:        methods,
		RoutesPackage:  hg.routePkgName,
		Level:          Title(hg.level),
		InitName:       initName,
		SubPaths:       subPaths,
	}
	return cnt, nil
}

func (hg *handlerGroup) buildHandlerMethods() (methods []string, err error) {
	i := 0
	methods = make([]string, len(hg.leaf.Node.Methods))
	for _, method := range hg.leaf.Node.Methods {
		methods[i], err = hg.buildHandlerMethod(method, hg.leaf.Node.URL)
		if err != nil {
			return nil, err
		}
		i++
	}
	return
}

func (hg *handlerGroup) buildHandlerMethod(methodDef *RouteDef, url string) (string, error) {
	tpl, err := template.ParseFiles("tpl/handler.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, hg.getHandlerMethod(methodDef, url)); err != nil {
		return "", err
	}
	return content.String(), nil
}

func (hg *handlerGroup) getHandlerMethod(methodDef *RouteDef, url string) *HandlerContent {
	requestAlias, responseAlias := "", ""
	if methodDef.Definition.RequestParam != nil {
		requestAlias = addPackageToMap(methodDef.Definition.RequestParam.Path, hg.packagesMap, 0)
	}
	if methodDef.Definition.ResponseParam != nil {
		responseAlias = addPackageToMap(methodDef.Definition.ResponseParam.Path, hg.packagesMap, 0)
	}
	ctn := &HandlerContent{
		LevelServer: hg.levelServer,
		Handler:     methodDef.Handler,
		Params:      hg.getRequestProto(methodDef, requestAlias),
		ReturnType:  hg.getReturnType(methodDef, responseAlias),
		Response:    hg.parseResponseObject(methodDef.Definition),
		Param:       methodDef.Param,
	}
	return ctn
}

func (hg *handlerGroup) getRequestProto(methodDef *RouteDef, alias string) string {
	params := []string{"g *gin.Context"}
	if methodDef.Param != "" {
		params = append(params, fmt.Sprintf("%s string", methodDef.Param))
	}
	if methodDef.Definition.Request != nil {
		params = append(params, methodDef.Definition.RequestParam.getObjectDeclaration())
	}
	return strings.Join(params, ", ")
}

func (hg *handlerGroup) getReturnType(methodDef *RouteDef, alias string) string {
	pre, post := "", ""
	params := make([]string, 0)
	if methodDef.Definition.ResponseParam != nil {
		pre, post = "(", ")"
		params = append(params, methodDef.Definition.ResponseParam.getTypeDeclaration())
	}
	params = append(params, "error")
	return fmt.Sprintf("%s%s%s", pre, strings.Join(params, ", "), post)
}

func (hg *handlerGroup) parseResponseObject(def *R) *RequestContent {
	if def.ResponseParam == nil {
		return nil
	}
	return &RequestContent{
		Name:            def.ResponseParam.Name,
		Type:            def.ResponseParam.Type,
		Alias:           addPackageToMap(def.ResponseParam.Path, hg.packagesMap, 0),
		TypeDeclaration: def.ResponseParam.getUnnamedObjectDeclaration(),
	}
}

func (hg *handlerGroup) buildInitInvokes() (string, []string) {
	initName := "Init"
	invokes := make([]string, len(hg.leaf.Tree))
	i := 0
	for _, child := range hg.leaf.Tree {
		pkgPath := hg.packagePath + "/" + hg.level
		pkgName := getLastComponent(pkgPath)
		if pkgName != "" {
			pkgName += "."
			addPackageToMap(pkgPath, hg.packagesMap, 0)
		} else {
			initName = "init"
		}
		invokes[i] = fmt.Sprintf("%sInit%s", pkgName, Title(cleanupRoute(child.Level)))
		i++
	}
	return initName, invokes
}
