package gorest

import (
	"fmt"
	"strings"
)

var root = &AST{
	Level: "/",
	Node: &Route{
		URL:     "/",
		FullURL: "/",
	},
	Tree: map[string]*AST{},
}

func getPath(url string, root *AST) (*AST, error) {
	if url == "/" {
		return root, nil
	}
	var leaf *AST
	components := strings.Split(url, "/")
	for i, component := range components {
		if component == "" {
			if i == 0 {
				leaf = root
			} else if len(components) > 2 {
				return nil, fmt.Errorf("invalid route %s", url)
			}
		} else {
			if leaf.Tree == nil {
				leaf.Tree = map[string]*AST{}
			}
			if _, present := leaf.Tree[component]; !present {
				leaf.Tree[component] = &AST{
					Level: component,
					Node: &Route{
						URL:     url,
						FullURL: leaf.append(url),
						Methods: map[string]*RouteDef{},
					},
					Tree: map[string]*AST{},
				}
			}
			leaf = leaf.Tree[component]
		}
	}
	if leaf == nil {
		return nil, fmt.Errorf("unprocessable route %s", url)
	}
	return leaf, nil
}
