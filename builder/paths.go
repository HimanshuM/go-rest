package builder

import (
	"errors"
	"os"
)

type routesPathSpec struct {
	// Path to the routes files
	FilePath string
	// Package name for routes
	PackageName string
	// Complete routes package
	PackagePath string
	// Relative directory path
	DirPath string
}

type handlersPathSpec struct {
	// Path to the handlers files
	FilePath string
	// Package name for handlers
	PackageName string
	// Complete handlers package
	PackagePath string
}

func buildLevel(leaf *AST, path *routesPathSpec) (err error) {
	if len(leaf.Tree) > 0 {
		err = os.Mkdir(path.DirPath, 0755)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				return
			}
		}
	}
	level := cleanupRoute(leaf.Level)
	if level != "" {
		if err = writeRoutesFile(level, leaf, path); err != nil {
			return
		}
	}
	if leaf.HasDefinition {
		if err = writeServerFile(level, leaf, path); err != nil {
			return
		}
	}
	for _, node := range leaf.Tree {
		dirPath := path.DirPath + level
		newPath := &routesPathSpec{
			FilePath:    dirPath,
			PackageName: level,
			PackagePath: path.PackagePath + level,
			DirPath:     dirPath + "/",
		}
		if err = buildLevel(node, newPath); err != nil {
			return
		}
	}
	return
}

func Generate() error {
	routePackage := getLastComponent(routesPkgPath)
	// handlerPackage := getLastComponent(handlersPkgPath)
	return buildLevel(root, &routesPathSpec{
		FilePath:    routePackage,
		PackageName: routePackage,
		PackagePath: routesPkgPath,
		DirPath:     routePackage + "/",
	})
}
