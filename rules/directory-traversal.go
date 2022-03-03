package rules

import (
	"go/ast"
	"regexp"

	"github.com/securego/gosec/v2"
)

type traversal struct {
	pattern          *regexp.Regexp
	gosec.MetaData
}

func (r *traversal) ID() string {
	return r.MetaData.ID
}

func (r *traversal) Match(n ast.Node, ctx *gosec.Context) (*gosec.Issue, error) {
	switch node := n.(type) {
	case *ast.CallExpr:
		return r.matchCallExpr(node, ctx)
	}
	return nil, nil
}

func (r *traversal) matchCallExpr(assign *ast.CallExpr, ctx *gosec.Context) (*gosec.Issue, error) {
	for _, i := range assign.Args {
		if basiclit, ok1 := i.(*ast.BasicLit); ok1 {
			if fun, ok2 := assign.Fun.(*ast.SelectorExpr); ok2 {
				if x, ok3 := fun.X.(*ast.Ident); ok3 {
					string := x.Name + "." + fun.Sel.Name + "(" + basiclit.Value + ")"
					if r.pattern.MatchString(string) {
						return gosec.NewIssue(ctx, assign, r.ID(), r.What, r.Severity, r.Confidence), nil
					}
				}
			}
		}
	}
	return nil, nil
}

// NewDirectoryTraversal attempts to find the use of http.Dir("/")
func NewDirectoryTraversal(id string, conf gosec.Config) (gosec.Rule, []ast.Node) {
	pattern := `http\.Dir\("\/"\)|http\.Dir\('\/'\)`
	if val, ok := conf["G101"]; ok {
		conf := val.(map[string]interface{})
		if configPattern, ok := conf["pattern"]; ok {
			if cfgPattern, ok := configPattern.(string); ok {
				pattern = cfgPattern
			}
		}
	}

	return &traversal{
		pattern:        regexp.MustCompile(pattern),
		MetaData: gosec.MetaData{
			ID:         id,
			What:       "Potential directory traversal",
			Confidence: gosec.Medium,
			Severity:   gosec.Medium,
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
