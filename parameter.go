package builder

import (
	"fmt"
	"reflect"
	"strings"
)

func (r *R) processDefinition() {
	if r.Request != nil {
		r.processRequest()
	}
	if r.Response != nil {
		r.processResponse()
	}
	if r.Error != nil {
		r.processError()
	}
}

func (r *R) processRequest() {
	request := reflect.TypeOf(r.Request)
	r.RequestParam = processType(request, false)
}

func (r *R) processResponse() {
	response := reflect.TypeOf(r.Response)
	r.ResponseParam = processType(response, false)
}

func (r *R) processError() {
	errType := reflect.TypeOf(r.Error)
	if r.ErrorParam = processType(errType, true); r.ErrorParam == nil {
		panic("Error: error parameter cannot be an array")
	}
}

func processType(typeObj reflect.Type, noArray bool) *Parameter {
	typeObj, isArray := getType(typeObj)
	if noArray && isArray {
		return nil
	}
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

func (p *Parameter) getUnnamedObjectDeclaration() string {
	decl := "&"
	if p.IsArray {
		decl = "[]*"
	}
	decl += fmt.Sprintf("%s.%s", p.Package, p.Type)
	return decl
}
