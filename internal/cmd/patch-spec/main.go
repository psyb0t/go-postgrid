// patch-spec collapses anyOf schemas containing only identical references.
// The four mail creation responses repeat their schema twice, causing
// oapi-codegen to emit duplicate As, From and Merge methods.
// It also names the inline text branch of the shared search parameter,
// which otherwise produces methods referring to an undefined N0 type.
package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	sourcePath = "api/openapi.yaml"
	outputPath = "api/openapi.patched.yaml"
	fileMode   = 0o600
	pairSize   = 2
	searchName = "ListSearchText"
	searchRef  = "#/components/schemas/" + searchName
)

func main() {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		panic(err)
	}

	var document yaml.Node

	if err := yaml.Unmarshal(data, &document); err != nil {
		panic(err)
	}

	collapseReferences(&document)
	nameSearchText(&document)

	patched, err := yaml.Marshal(&document)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(outputPath, patched, fileMode); err != nil {
		panic(err)
	}
}

func collapseReferences(node *yaml.Node) {
	for _, child := range node.Content {
		collapseReferences(child)
	}

	if node.Kind != yaml.MappingNode || len(node.Content) != pairSize {
		return
	}

	if node.Content[0].Value != "anyOf" {
		return
	}

	branches := node.Content[1]
	if branches.Kind != yaml.SequenceNode || len(branches.Content) < pairSize {
		return
	}

	first := branches.Content[0]

	reference := schemaReference(first)
	if reference == "" {
		return
	}

	for _, branch := range branches.Content[1:] {
		if schemaReference(branch) != reference {
			return
		}
	}

	*node = *first
}

func schemaReference(node *yaml.Node) string {
	if node.Kind != yaml.MappingNode || len(node.Content) != pairSize {
		return ""
	}

	if node.Content[0].Value != "$ref" {
		return ""
	}

	return node.Content[1].Value
}

func nameSearchText(document *yaml.Node) {
	components := mappingValue(document.Content[0], "components")
	parameters := mappingValue(components, "parameters")
	search := mappingValue(parameters, "ListParams.search")
	schema := mappingValue(search, "schema")
	branches := mappingValue(schema, "anyOf")

	text := branches.Content[0]
	if schemaReference(text) == searchRef {
		return
	}

	if len(text.Content) != pairSize ||
		mappingValue(text, "type").Value != "string" {
		panic("expected an inline string branch in ListParams.search")
	}

	schemas := mappingValue(components, "schemas")
	for i := 0; i < len(schemas.Content); i += pairSize {
		if schemas.Content[i].Value == searchName {
			panic("search schema name already exists")
		}
	}

	schemas.Content = append(schemas.Content, scalar(searchName), text)
	branches.Content[0] = &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			scalar("$ref"), scalar(searchRef),
		},
	}
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(node.Content); i += pairSize {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	panic("missing YAML mapping key: " + key)
}

func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value}
}
