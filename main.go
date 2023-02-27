package main

import (
	"fmt"
	"strings"

	"github.com/HimanshuM/go-rest-builder/builder"
	"github.com/HimanshuM/go-rest-builder/models"
)

func definitions() {
	builder.Package("github.com/HimanshuM/go-rest-builder/routes")
	root, _ := builder.Path("/")
	apiV1, _ := root.Path("/api/v1")
	apiV2, _ := root.Path("/api/v2")
	schools, _ := apiV1.Path("/schools")
	schools.GET(&builder.R{}).POST(&builder.R{})
	schoolID, _ := schools.Path("/{id}")
	schoolID.GET(&builder.R{}).PATCH(&builder.R{}).DELETE(&builder.R{})
	stds, _ := apiV1.Path("/standards")
	stds.GET(&builder.R{Response: &models.StandardResponse{}}).POST(&builder.R{Request: &models.StandardRequest{}, Response: &models.StandardResponse{}})
	stdID, _ := stds.Path("/{id}")
	stdID.GET(&builder.R{}).PATCH(&builder.R{}).DELETE(&builder.R{})
	sections, _ := apiV1.Path("/sections")
	sections.GET(&builder.R{}).POST(&builder.R{})
	secID, _ := sections.Path("/{id}")
	secID.GET(&builder.R{}).PATCH(&builder.R{}).DELETE(&builder.R{})
	v2Admin, _ := apiV2.Path("/admin")
	v2Admin.GET(&builder.R{})
}

func stringifyHandlers(node *builder.Route) string {
	handlers := make([]string, len(node.Methods))
	if len(node.Methods) == 0 {
		return ""
	}
	i := 0
	for method, def := range node.Methods {
		handlers[i] = fmt.Sprintf("%s: %s", method, def.Handler)
		i++
	}
	return ", Handlers: " + strings.Join(handlers, "\t")
}

func printLeaf(leaf *builder.AST, indent int) {
	spaces := strings.Repeat(" ", indent)
	fmt.Printf("%sLEAF: %s (%t)", spaces, leaf.Level, leaf.HasDefinition)
	if leaf.Node != nil {
		handlers := stringifyHandlers(leaf.Node)
		fmt.Printf(" NODE: %s%s", leaf.Node.FullURL, handlers)
	}
	fmt.Println()
	if leaf.Tree == nil {
		return
	}
	for _, node := range leaf.Tree {
		printLeaf(node, indent+4)
	}
}

func main() {
	definitions()
	if err := builder.BuildPaths(); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
}
