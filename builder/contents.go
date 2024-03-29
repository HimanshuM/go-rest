package builder

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"text/template"
)

type serverGroup struct {
	packagesMap map[string]importDef
	path        string
	level       string
	pkg         string
	packagePath string
	leaf        *AST
}

type RoutesContent struct {
	Package     string
	Imports     string
	Server      string
	Route       string
	Lines       []string
	Level       string
	SetupCase   string
	Middlewares []string
}

// writeRoutesFile writes each routes file
/*
 * path: 	The file path
 * level: 	The route level that this set of routes will handle
 * pkg: 	Package name for this file
 * pkgPath: Complete package path
 */
func writeRoutesFile(level string, leaf *AST, path *PathSpec) error {
	sg := serverGroup{
		packagesMap: map[string]importDef{},
		path:        path.RouteFilePath,
		level:       level,
		pkg:         path.RoutePackageName,
		packagePath: path.RoutePackagePath,
		leaf:        leaf,
	}
	addPackageToMap("github.com/gin-gonic/gin", sg.packagesMap, 0)
	filepath := path.RouteFilePath + ".go"
	if level == "" {
		filepath = fmt.Sprintf("%s/%s", path.RouteFilePath, filepath)
	}
	hnd, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("error creating file: %v\n", err)
		return err
	}
	defer hnd.Close()
	return sg.writeRoutesContent(hnd)
}

func (sg *serverGroup) writeRoutesContent(hnd *os.File) error {
	tpl, err := template.ParseFiles("tpl/base_routes.tpl")
	if err != nil {
		return err
	}
	content, err := sg.buildRoutesContent()
	if err != nil {
		return err
	}
	return tpl.Execute(hnd, content)
}

func (sg *serverGroup) buildRoutesContent() (*RoutesContent, error) {
	if len(sg.leaf.Tree) > 0 && sg.level != "" {
		addPackageToMap(fmt.Sprintf("%s/%s", sg.packagePath, sg.level), sg.packagesMap, 0)
	}
	middlewares := sg.buildMiddlewares()
	importStrings := []importDef{}
	for _, pkg := range sg.packagesMap {
		importStrings = append(importStrings, pkg)
	}
	imports, err := imports(importStrings...)
	if err != nil {
		return nil, err
	}

	prefix := sg.level
	setupCase := "S"
	if sg.level == "" {
		prefix = "root"
		setupCase = "s"
	}
	server := ""
	if len(sg.leaf.Tree) > 0 {
		server = prefix + "Router"
	}
	ctn := &RoutesContent{
		Package:     sg.pkg,
		Imports:     imports,
		Server:      server,
		Route:       sg.leaf.Level,
		Level:       Title(sg.level),
		SetupCase:   setupCase,
		Middlewares: middlewares,
	}
	ctn.linesFromRoute(sg.leaf, sg.level, ctn.Server)
	return ctn, nil
}

func (sg *serverGroup) buildMiddlewares() []string {
	middlewares := make([]string, len(sg.leaf.Middlewares))
	for i, middleware := range sg.leaf.Middlewares {
		mw := runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name()
		mwName := getLastComponentBySeparator(mw, ".")
		mwPkgPath := strings.Replace(mw, "."+mwName, "", 1)
		alias := addPackageToMap(mwPkgPath, sg.packagesMap, 0)
		middlewares[i] = fmt.Sprintf("%s.%s", alias, mwName)
	}
	return middlewares
}

func (ctn *RoutesContent) linesFromRoute(leaf *AST, level, server string) {
	i := 0
	ctn.Lines = make([]string, len(leaf.Tree)+len(leaf.Node.Methods))
	if level != "" {
		level += "."
	}
	for _, node := range leaf.Tree {
		ctn.Lines[i] = fmt.Sprintf("%sSetup%sRoutes(%s)", level, Title(cleanupRoute(node.Level)), server)
		i++
	}
	for method, def := range leaf.Node.Methods {
		ctn.Lines[i] = fmt.Sprintf("server.%s(\"%s\", %s)", method, leaf.Node.URL, def.Handler)
		i++
	}
}
