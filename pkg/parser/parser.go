package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	treesitter "github.com/tree-sitter/go-tree-sitter"
	treesitter_py "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// SchemaNode represents an abstract schema element extracted from treesitter output
type SchemaNode struct {
	Identifier string       // Name or identifier of the schema element
	Value      any          // Value associated with the identifier (could be string, number, bool, or nested SchemaNode)
	NodeType   string       // Type of the node (e.g., "variable_declaration", "object", "array", etc.)
	Children   []SchemaNode // Child nodes for nested structures
}

func dump(n *treesitter.Node, src []byte, indent int) {
	fmt.Printf("%s%s → %q\n",
		strings.Repeat("  ", indent),
		n.Kind(),        // instead of n.Type()
		n.Utf8Text(src), // instead of n.Content(src)
	)
	for i := uint(0); i < n.ChildCount(); i++ {
		dump(n.Child(i), src, indent+1)
	}
}

func openFiles(filename string) []byte {
	cwd, _ := os.Getwd()
	fn := filepath.Join(cwd, filename)
	file, _ := os.ReadFile(fn)
	return file
}

func returnTextWithAgreeString(text []byte) []string {
	// Convert byte slice to string for easier handling
	content := string(text)

	// Split the content by lines
	lines := strings.Split(content, "\n")

	// Create a slice to hold lines containing "[agree"
	agreeLines := []string{}

	// Iterate through each line checking for "[agree"
	for _, line := range lines {
		if strings.Contains(line, "[agree") {
			agreeLines = append(agreeLines, line)
		}
	}

	return agreeLines
}

func ParseBytes() {
	code := openFiles("test-data/pydantic.py")

	parser := treesitter.NewParser()
	defer parser.Close()

	language := treesitter.NewLanguage(treesitter_py.Language())
	parser.SetLanguage(language)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query, err := treesitter.NewQuery(language, `

		`)
	if err != nil {
		panic(err)
	}

	defer query.Close()

	qc := treesitter.NewQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), code)

	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		fmt.Printf(
			"Capture %d: %s\n",
			index,
			match.Captures[index].Node.Utf8Text(code),
		)
	}
}
