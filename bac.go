package builder

import (
	"fmt"
	"strings"

	models "github.com/HimanshuM/go-rest-builder/v1/model"
)

func definitions() {
	RoutesPackage("github.com/HimanshuM/go-rest-builder/v1/routes")
	HandlersPackage("github.com/HimanshuM/go-rest-builder/v1/apis")
	root, _ := Path("/")
	apiV1, _ := root.Path("/api/v1", models.Authenticate)
	apiV2, _ := root.Path("/api/v2")
	schools, _ := apiV1.Path("/schools")
	schools.GET(&R{}).POST(&R{})
	schoolID, _ := schools.Path("/{id}")
	schoolID.GET(&R{}).PATCH(&R{}).DELETE(&R{})
	stds, _ := apiV1.Path("/standards")
	stds.GET(&R{Response: []*models.StandardResponse{}}).POST(&R{Request: []*models.StandardRequest{}, Response: []*models.StandardResponse{}})
	stdID, _ := stds.Path("/{id}")
	stdID.GET(&R{Response: &models.StandardResponse{}}).PATCH(&R{Response: &models.StandardResponse{}}).DELETE(&R{})
	sections, _ := apiV1.Path("/sections")
	sections.GET(&R{}).POST(&R{})
	secID, _ := sections.Path("/{id}")
	secID.GET(&R{}).PATCH(&R{}).DELETE(&R{})
	v2Admin, _ := apiV2.Path("/admin")
	v2Admin.GET(&R{})
}

func stringifyHandlers(node *Route) string {
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

func printLeaf(leaf *AST, indent int) {
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
	if err := Generate(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}
