package notion

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GoldmarkRenderer struct {
	Color      renderer.NodeRendererFunc
	Equation   renderer.NodeRendererFunc
	UserPerson renderer.NodeRendererFunc
	NotionFile renderer.NodeRendererFunc
}

func noopRenderer(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func printNode(_ util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		fmt.Printf("%s (%T) not implemented yet\n", n.Kind(), n)
	}

	return ast.WalkContinue, nil
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *GoldmarkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	if r.Color == nil {
		reg.Register(GoldmarkKindColor, noopRenderer)
	} else {
		reg.Register(GoldmarkKindColor, r.Color)
	}

	if r.Equation == nil {
		reg.Register(GoldmarkKindEquation, noopRenderer)
	} else {
		reg.Register(GoldmarkKindEquation, r.Equation)
	}

	if r.UserPerson == nil {
		reg.Register(GoldmarkKindUserPerson, noopRenderer)
	} else {
		reg.Register(GoldmarkKindUserPerson, r.UserPerson)
	}

	// if r.NotionFile == nil {
	// 	reg.Register(GoldmarkKindNotionFile, noopRenderer)
	// } else {
	// 	reg.Register(GoldmarkKindNotionFile, r.NotionFile)
	// }

	reg.Register(GoldmarkKindExternalFile, printNode)

	reg.Register(GoldmarkKindBookmark, printNode)
	reg.Register(GoldmarkKindLinkPreview, printNode)
	reg.Register(GoldmarkKindUserBot, printNode)

	reg.Register(ast.KindString, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.Write(n.(*ast.String).Value)
		} else {
			s := n.NextSibling()
			if s != nil && s.Kind() == ast.KindList {
				w.WriteString("\n")
			}
		}

		return ast.WalkContinue, nil
	})
}
