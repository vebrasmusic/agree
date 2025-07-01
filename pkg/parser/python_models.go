package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ts "github.com/tree-sitter/go-tree-sitter"
	py "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// Field represents a single field with a name and type.
type Field struct {
	Name string
	Type string
}

// Model represents a parsed model with its fields.
type Model struct {
	Name   string
	Fields map[string]Field
}

// agreeBlock represents a single agree section inside a file.
type agreeBlock struct {
	Nickname string
	Type     string
	Code     string
}

// ParsePythonFiles walks the given directory, reads all files and parses agree
// blocks tagged as either "sqlalchemy" or "pydantic".
func ParsePythonFiles(dir string) (map[string]Model, map[string]Model, error) {
	sqlModels := make(map[string]Model)
	pydModels := make(map[string]Model)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		blocks := extractAgreeBlocks(string(src))
		for _, b := range blocks {
			switch b.Type {
			case "sqlalchemy":
				m, err := parsePythonModel([]byte(b.Code), "sqlalchemy")
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				sqlModels[b.Nickname] = m
			case "pydantic":
				m, err := parsePythonModel([]byte(b.Code), "pydantic")
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				pydModels[b.Nickname] = m
			}
		}
		return nil
	})
	return sqlModels, pydModels, err
}

// extractAgreeBlocks scans the given text for agree comment blocks.
func extractAgreeBlocks(src string) []agreeBlock {
	lines := strings.Split(src, "\n")
	var blocks []agreeBlock
	var current *agreeBlock
	for _, line := range lines {
		if current == nil {
			if idx := strings.Index(line, "[agree:"); idx != -1 {
				rest := line[idx+len("[agree:"):]
				end := strings.Index(rest, "]")
				if end == -1 {
					continue
				}
				header := rest[:end]
				parts := strings.SplitN(header, ":", 2)
				if len(parts) != 2 {
					continue
				}
				current = &agreeBlock{Nickname: strings.TrimSpace(parts[0]), Type: strings.TrimSpace(parts[1])}
				current.Code = ""
			}
			continue
		}
		if strings.Contains(line, "[agree:end]") {
			blocks = append(blocks, *current)
			current = nil
			continue
		}
		current.Code += line + "\n"
	}
	return blocks
}

// parsePythonModel parses one Python class definition from src. modelType must
// be either "sqlalchemy" or "pydantic".
func parsePythonModel(src []byte, modelType string) (Model, error) {
	parser := ts.NewParser()
	parser.SetLanguage(ts.NewLanguage(py.Language()))
	tree := parser.Parse(src, nil)
	defer tree.Close()

	root := tree.RootNode()
	for i := uint(0); i < root.NamedChildCount(); i++ {
		n := root.NamedChild(i)
		if n.Kind() != "class_definition" {
			continue
		}
		nameNode := n.ChildByFieldName("name")
		if nameNode == nil {
			continue
		}
		className := nameNode.Utf8Text(src)
		body := n.ChildByFieldName("body")
		fields := make(map[string]Field)
		if body != nil {
			for j := uint(0); j < body.NamedChildCount(); j++ {
				stmt := body.NamedChild(j)
				if stmt.Kind() != "expression_statement" {
					continue
				}
				assign := stmt.NamedChild(0)
				if assign == nil || assign.Kind() != "assignment" {
					continue
				}
				left := assign.ChildByFieldName("left")
				if left == nil || left.Kind() != "identifier" {
					continue
				}
				fieldName := left.Utf8Text(src)
				if strings.HasPrefix(fieldName, "__") {
					continue
				}
				var fieldType string
				if t := assign.ChildByFieldName("type"); t != nil {
					fieldType = t.Utf8Text(src)
				}
				if modelType == "sqlalchemy" {
					r := assign.ChildByFieldName("right")
					if r == nil || r.Kind() != "call" {
						continue
					}
					fn := r.ChildByFieldName("function")
					if fn == nil || fn.Utf8Text(src) != "Column" {
						continue
					}
					args := r.ChildByFieldName("arguments")
					if args != nil && args.NamedChildCount() > 0 {
						first := args.NamedChild(0)
						fieldType = first.Utf8Text(src)
					}
				}
				fields[fieldName] = Field{Name: fieldName, Type: normalizeType(fieldType)}
			}
		}
		return Model{Name: className, Fields: fields}, nil
	}
	return Model{}, fmt.Errorf("no class definition found")
}

// normalizeType normalizes simple python/sqlalchemy type names.
func normalizeType(t string) string {
	t = strings.TrimSpace(strings.ToLower(t))
	if idx := strings.Index(t, "|"); idx > 0 {
		t = strings.TrimSpace(t[:idx])
	}
	switch t {
	case "integer", "int":
		return "int"
	case "string", "str":
		return "str"
	case "float":
		return "float"
	case "boolean", "bool":
		return "bool"
	case "emailstr":
		return "emailstr"
	}
	return t
}

// CompareModels compares SQLAlchemy models with Pydantic models and returns a report.
func CompareModels(sqlModels, pydModels map[string]Model) string {
	var sb strings.Builder
	for nick, sqlModel := range sqlModels {
		pydModel, ok := pydModels[nick]
		if !ok {
			continue
		}
		missingSQL := []string{}
		missingPyd := []string{}
		typeMismatch := []string{}
		for fname, f := range pydModel.Fields {
			sf, ok := sqlModel.Fields[fname]
			if !ok {
				missingSQL = append(missingSQL, fname)
				continue
			}
			if sf.Type != f.Type {
				typeMismatch = append(typeMismatch, fmt.Sprintf("%s (%s != %s)", fname, sf.Type, f.Type))
			}
		}
		for fname := range sqlModel.Fields {
			if _, ok := pydModel.Fields[fname]; !ok {
				missingPyd = append(missingPyd, fname)
			}
		}
		if len(missingSQL)+len(missingPyd)+len(typeMismatch) > 0 {
			sb.WriteString(fmt.Sprintf("Model %s:\n", nick))
			if len(missingSQL) > 0 {
				sb.WriteString("  Missing in SQLAlchemy: " + strings.Join(missingSQL, ", ") + "\n")
			}
			if len(missingPyd) > 0 {
				sb.WriteString("  Missing in Pydantic: " + strings.Join(missingPyd, ", ") + "\n")
			}
			if len(typeMismatch) > 0 {
				sb.WriteString("  Type mismatches: " + strings.Join(typeMismatch, ", ") + "\n")
			}
		}
	}
	if sb.Len() == 0 {
		return "No mismatches found"
	}
	return sb.String()
}
