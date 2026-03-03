package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "MyLint",
	Doc:  "checks log messages for style and security issues",
	Run:  run,
}
