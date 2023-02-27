package builder

import (
	"errors"
	"os"
)

func buildLevel(leaf *AST, pkg, dir, pkgPath string) (err error) {
	level := cleanupRoute(leaf.Level)
	if level != "" {
		if err = writeRoutesFile(dir, level, pkg, pkgPath, leaf); err != nil {
			return
		}
	}
	if leaf.HasDefinition {
		if err = writeServerFile(dir, level, pkg, pkgPath, leaf); err != nil {
			return
		}
	}
	dir += level
	if len(leaf.Tree) > 0 {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				return
			}
		}
	}
	if level != "" {
		pkgPath += "/" + level
	}
	if level == "" {
		level = pkg
	}
	for _, node := range leaf.Tree {
		if err = buildLevel(node, level, dir+"/", pkgPath); err != nil {
			return
		}
	}
	return
}

func Generate() error {
	rootPackage := getLastComponent(pkgPath)
	return buildLevel(root, rootPackage, rootPackage, pkgPath)
}
