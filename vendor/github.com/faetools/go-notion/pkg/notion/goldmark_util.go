package notion

import "github.com/yuin/goldmark/ast"

func wrapNode(parent, child ast.Node) ast.Node {
	parent.AppendChild(parent, child)
	return parent
}

func stringNode(s string) *ast.String {
	return ast.NewString([]byte(s))
}
