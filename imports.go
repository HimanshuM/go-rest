package gorest

import (
	"bytes"
	"sort"
	"strconv"
	"text/template"
)

type importDef struct {
	Alias string
	Path  string
	Name  string
}

func imports(pkgs ...importDef) (string, error) {
	if len(pkgs) == 0 {
		return "", nil
	}
	if len(pkgs) == 1 {
		return writeImportContent(pkgs[0])
	}
	return writeImportsContent(sortImports(pkgs))
}

func writeImportContent(pkg importDef) (string, error) {
	tpl, err := template.ParseFiles("tpl/import.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, pkg); err != nil {
		return "", err
	}
	return content.String(), nil
}

func writeImportsContent(pkgs []importDef) (string, error) {
	tpl, err := template.ParseFiles("tpl/imports.tpl")
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, pkgs); err != nil {
		return "", err
	}
	return content.String(), nil
}

func addPackageToMap(pkg string, packages map[string]importDef, suffix int) string {
	pkgName := getLastComponent(pkg)
	origName := pkgName
	if suffix > 0 {
		pkgName += strconv.Itoa(suffix)
	}
	if existingPkg, present := packages[pkgName]; present {
		if existingPkg.Path != pkg {
			return addPackageToMap(pkg, packages, suffix+1)
		}
	} else {
		packages[pkgName] = importDef{Path: pkg, Alias: pkgName, Name: origName}
	}
	return pkgName
}

func sortImports(pkgs []importDef) []importDef {
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Path < pkgs[j].Path
	})
	return pkgs
}
