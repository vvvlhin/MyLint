package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				if node.Type == nil || node.Type.Params == nil {
					return true
				}
				if node.Type.Params.NumFields() > 3 {
					pass.Reportf(node.Type.Params.Pos(), "fn %s has more than 3 params", node.Name.Name)
				}
				if isSnakeCase(node.Name.Name) {
					newName := toCamelCase(node.Name.Name)
					// pass.Reportf(node.Name.Pos(), "fn %s in snake case", node.Name.Name)
					pass.Report(analysis.Diagnostic{
						Pos:     node.Name.Pos(),
						End:     node.Name.End(),
						Message: fmt.Sprintf("fn %s in snake case", node.Name.Name),
						SuggestedFixes: []analysis.SuggestedFix{
							{

								TextEdits: []analysis.TextEdit{
									{

										Pos:     node.Name.Pos(),
										End:     node.Name.End(),
										NewText: []byte(newName),
									},
								},
							},
						},
					})
				}
			case *ast.ValueSpec:
				for _, indent := range node.Names {
					checkVarName(pass, indent)
				}
			case *ast.AssignStmt:
				for _, expr := range node.Lhs {
					if indent, ok := expr.(*ast.Ident); ok {
						checkVarName(pass, indent)
					}
				}
			case *ast.RangeStmt:
				if indent, ok := node.Key.(*ast.Ident); ok {
					checkVarName(pass, indent)
				}
				if indent, ok := node.Value.(*ast.Ident); ok {
					checkVarName(pass, indent)
				}
			}
			return true
		})
	}
	return nil, nil
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i == 0 {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func checkVarName(pass *analysis.Pass, indent *ast.Ident) {
	if isSnakeCase(indent.Name) {
		pass.Reportf(indent.Pos(), "var %s in snake case", indent.Name)
	}
}

func isSnakeCase(s string) bool {
	if s == "_" {
		return false
	}
	return strings.Contains(s, "_")
}
