package builder

type R struct {
	Query    interface{}
	Request  interface{}
	Response interface{}
}

type RouteDef struct {
	Method     string
	Handler    string
	Definition *R
}

type Route struct {
	URL     string
	FullURL string
	Methods map[string]*RouteDef
}

type AST struct {
	Level         string
	Node          *Route
	Tree          map[string]*AST
	HasDefinition bool
	Package       string
}
