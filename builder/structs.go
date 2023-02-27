package builder

type Parameter struct {
	Type    string
	Name    string
	Package string
	Path    string
	IsArray bool
}

type R struct {
	Query         interface{}
	Request       interface{}
	Response      interface{}
	Error         interface{}
	RequestParam  *Parameter
	ResponseParam *Parameter
	ErrorParam    *Parameter
}

type RouteDef struct {
	Method     string
	Handler    string
	Definition *R
	Param      string
	Comment    string
}

type Route struct {
	URL     string
	FullURL string
	Methods map[string]*RouteDef
	Comment string
}

type AST struct {
	Level         string
	Node          *Route
	Tree          map[string]*AST
	HasDefinition bool
	Package       string
}

func (r *R) processDefinition() {
	if r.Request != nil {
		r.processRequest()
	}
	if r.Response != nil {
		r.processResponse()
	}
}
