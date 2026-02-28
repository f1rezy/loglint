package logcheck

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var repeatedPunct = regexp.MustCompile(`[!?*]+`).MatchString
var isLetter = regexp.MustCompile(`[^\x00-\x7F]`).MatchString

var sensitiveWords = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private_key",
	"jwt",
}

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages style",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			method := sel.Sel.Name
			if !isLogMethod(method) {
				return true
			}

			if len(call.Args) == 0 {
				return true
			}

			msg := extractMessage(call.Args[0])
			if msg == "" {
				return true
			}

			checkRules(pass, call, call.Args[0], msg)

			return true
		})
	}
	return nil, nil
}

func isLogMethod(name string) bool {
	return name == "Info" || name == "Error" || name == "Warn" || name == "Debug"
}

func extractMessage(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			return strings.Trim(v.Value, `"`)
		}

	case *ast.BinaryExpr:
		if v.Op == token.ADD {
			return extractMessage(v.X) + extractMessage(v.Y)
		}

	case *ast.Ident:
		return v.Name
	}
	return ""
}

func checkRules(pass *analysis.Pass, call *ast.CallExpr, expr ast.Expr, msg string) {
	if len(msg) > 0 && unicode.IsUpper(rune(msg[0])) {
		pass.Reportf(call.Pos(), "log message should start with lowercase letter")
	}

	if isLetter(msg) {
		pass.Reportf(call.Pos(), "log message should contain only english characters")
	}

	if containsForbiddenSymbols(msg) {
		pass.Reportf(call.Pos(), "log message contains forbidden symbols or emoji")
	}

	if repeatedPunct(msg) {
		pass.Reportf(call.Pos(), "log message contains excessive punctuation (!, ?, *)")
	}

	checkSensitive(pass, call, expr)
}

func containsForbiddenSymbols(msg string) bool {
	for _, r := range msg {
		if unicode.Is(unicode.So, r) {
			return true
		}

		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func checkSensitive(pass *analysis.Pass, call *ast.CallExpr, expr ast.Expr) {
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return
	}

	right := extractMessage(bin.Y)

	if containsSensitiveWord(right) {
		pass.Reportf(call.Pos(), "log message may contain sensitive data")
	}
}

func containsSensitiveWord(s string) bool {
	s = strings.ToLower(s)
	for _, w := range sensitiveWords {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}
