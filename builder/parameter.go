package builder

import (
	"fmt"
	"reflect"
	"strings"
)

func (r *R) processRequest() {
	request := reflect.TypeOf(r.Request)
	r.RequestParam = processType(request)
}

func (r *R) processResponse() {
	response := reflect.TypeOf(r.Response)
	r.ResponseParam = processType(response)
}

func processType(typeObj reflect.Type) *Parameter {
	typeObj, isArray := getType(typeObj)
	typeName := typeObj.Name()
	pkgPath := typeObj.PkgPath()
	return &Parameter{
		Type:    typeName,
		Name:    strings.ToLower(typeName[0:1]) + typeName[1:],
		Package: getLastComponent(pkgPath),
		Path:    pkgPath,
		IsArray: isArray,
	}
}

func getType(typeObj reflect.Type) (reflect.Type, bool) {
	isArray := false
	for {
		switch typeObj.Kind() {
		case reflect.Pointer:
			typeObj = typeObj.Elem()
		case reflect.Array, reflect.Slice:
			typeObj = typeObj.Elem()
			isArray = true
		default:
			return typeObj, isArray
		}
	}
}

func (p *Parameter) getObjectDeclaration() string {
	return fmt.Sprintf("%s %s", p.Name, p.getTypeDeclaration())
}

func (p *Parameter) getTypeDeclaration() string {
	decl := "*"
	if p.IsArray {
		decl = "[]*"
	}
	decl += fmt.Sprintf("%s.%s", p.Package, p.Type)
	return decl
}
