package builder

import (
	"errors"
	"fmt"
	"os"
)

type PathSpec struct {
	// Path to the routes files
	RouteFilePath string
	// Package name for routes
	RoutePackageName string
	// Complete routes package
	RoutePackagePath string
	// Relative directory path
	RouteDirPath string
	// Path to the handlers files
	HandlerFilePath string
	// Package name for handlers
	HandlerPackageName string
	// Complete handlers package
	HandlerPackagePath string
	// Relative directory path
	HandlerDirPath string
}

func getPackageInfo(level, name, path string) (string, string) {
	if level != "" {
		path += "/" + level
		name = level
	}
	return name, path
}

func mkDir(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	return nil
}

func buildRoutes(node *AST, path *PathSpec) (err error) {
	if len(node.Tree) > 0 {
		if err = mkDir(path.RouteDirPath); err != nil {
			return
		}
		if err = mkDir(path.HandlerDirPath); err != nil {
			return
		}
	}
	level := cleanupRoute(node.Level)
	if level != "" {
		if err = writeRoutesFile(level, node, path); err != nil {
			return
		}
	}
	if node.HasDefinition {
		if err = writeServerFile(level, node, path); err != nil {
			return
		}
		if err = writeHandlerFile(level, node, path); err != nil {
			return
		}
	}
	routePkg, routePkgPath := getPackageInfo(level, path.RoutePackageName, path.RoutePackagePath)
	handlerPkg, handlerPkgPath := getPackageInfo(level, path.HandlerPackageName, path.HandlerPackagePath)

	for _, child := range node.Tree {
		childLevel := cleanupRoute(child.Level)
		routeDirPath := fmt.Sprintf("%s/%s", path.RouteFilePath, childLevel)
		handlerDirPath := fmt.Sprintf("%s/%s", path.HandlerFilePath, childLevel)
		newPath := &PathSpec{
			RouteFilePath:      routeDirPath,
			RoutePackageName:   routePkg,
			RoutePackagePath:   routePkgPath,
			RouteDirPath:       routeDirPath + "/",
			HandlerFilePath:    handlerDirPath,
			HandlerPackageName: handlerPkg,
			HandlerPackagePath: handlerPkgPath,
			HandlerDirPath:     handlerDirPath + "/",
		}
		if err = buildRoutes(child, newPath); err != nil {
			return
		}
	}
	return
}

func Generate() error {
	routePackage := getLastComponent(routesPkgPath)
	handlerPackage := getLastComponent(handlersPkgPath)
	return buildRoutes(root, &PathSpec{
		RouteFilePath:      routePackage,
		RoutePackageName:   routePackage,
		RoutePackagePath:   routesPkgPath,
		RouteDirPath:       routePackage + "/",
		HandlerFilePath:    handlerPackage,
		HandlerPackageName: handlerPackage,
		HandlerPackagePath: handlersPkgPath,
		HandlerDirPath:     handlerPackage + "/",
	})
}
