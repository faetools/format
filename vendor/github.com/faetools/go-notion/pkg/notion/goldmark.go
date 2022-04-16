package notion

import (
	"fmt"
	"log"

	"github.com/yuin/goldmark/ast"
	_ast "github.com/yuin/goldmark/extension/ast"
)

// Nodes returns the goldmark nodes of the notion blocks.
func (bs Blocks) Nodes() []ast.Node {
	ns := []ast.Node{}

	var list *ast.List

	for _, b := range bs {
		switch b.Type {
		case BlockTypeBulletedListItem:
			switch {
			case list == nil:
				list = ast.NewList('-')
			case list.IsOrdered():
				ns = append(ns, list)
				list = ast.NewList('-')
			}

			list.AppendChild(list, b.Node())

		case BlockTypeNumberedListItem:
			switch {
			case list == nil:
				list = ast.NewList('.')
			case !list.IsOrdered():
				ns = append(ns, list)
				list = ast.NewList('.')
			}

			list.AppendChild(list, b.Node())
		default:
			ns = append(ns, b.Node())
		}
	}

	if list != nil {
		ns = append(ns, list)
	}

	return ns
}

// Node returns the goldmark node of this notion block.
func (b Block) Node() ast.Node {
	switch b.Type {
	case BlockTypeBookmark:
		return b.Bookmark.Node()
	case BlockTypeBulletedListItem:
		return b.BulletedListItem.Node()
	// case BlockTypeCallout:
	// 	return b.Callout.Node()
	// case BlockTypeChildDatabase:
	// case BlockTypeChildPage:
	// case BlockTypeColumn:
	// 	return b.Column.Node()
	// case BlockTypeColumnList:
	// 	return b.ColumnList.Node()
	// case BlockTypeDivider:
	// 	return b.Divider.Node()
	// case BlockTypeEmbed:
	// 	return b.Embed.Node()
	// case BlockTypeEquation:
	// 	return b.Equation.Node()
	// case BlockTypeFile:
	// 	return b.File.Node()
	// case BlockTypeHeading1:
	// 	return b.Heading1.Node()
	// case BlockTypeHeading2:
	// 	return b.Heading2.Node()
	// case BlockTypeHeading3:
	// 	return b.Heading3.Node()
	case BlockTypeImage:
		link := ast.NewLink()
		link.Destination = []byte(b.Image.URL())
		return ast.NewImage(link)
	case BlockTypeLinkPreview:
		return b.LinkPreview.Node()
	case BlockTypeLinkToPage:
		return b.LinkToPage.Node()
	// case BlockTypeNumberedListItem:
	// 	return b.NumberedListItem.Node()
	case BlockTypeParagraph:
		return b.Paragraph.Node()
	// case BlockTypePdf:
	// 	return b.Pdf.Node()
	// case BlockTypeQuote:
	// 	return b.Quote.Node()
	// case BlockTypeSyncedBlock:
	// 	return b.SyncedBlock.Node()
	// case BlockTypeTable:
	// 	return b.Table.Node()
	// case BlockTypeTableOfContents:
	// 	return b.TableOfContents.Node()
	// case BlockTypeTableRow:
	// 	return b.TableRow.Node()
	// case BlockTypeTemplate:
	// 	return b.Template.Node()
	// case BlockTypeToDo:
	// 	return b.ToDo.Node()
	// case BlockTypeToggle:
	// 	return b.Toggle.Node()
	// case BlockTypeUnsupported:
	// 	return b.Unsupported.Node()
	// case BlockTypeVideo:
	// 	return b.Video.Node()
	default:
		log.Fatalf("Node() for block type %q not implemented yet", b.Type)
		return nil
	}
}

// Node returns the goldmark node of this notion paragraph.
func (p Paragraph) Node() ast.Node {
	if len(p.Children) > 0 {
		log.Fatalf("Node() for paragraph with children not implemented yet")
	}

	var n ast.Node = ast.NewParagraph()

	p.RichText.appendTo(n)
	n = p.Color.wrapNode(n)

	return n
}

func (c Color) wrapNode(n ast.Node) ast.Node {
	if c == ColorDefault {
		return n
	}

	return wrapNode(c.Node(), n)
}

// Node returns the goldmark node of this notion color element.
func (c Color) Node() ast.Node { return &GoldmarkColor{Color: c} }

// Node returns the goldmark node of this notion equation element.
func (eq Equation) Node() ast.Node {
	return wrapNode(&GoldmarkEquation{}, stringNode(eq.Expression))
}

func (rts RichTexts) appendTo(n ast.Node) {
	for _, rt := range rts {
		n.AppendChild(n, rt.Node())
	}
}

// Node returns the goldmark node of this notion rich text element.
func (t RichText) Node() ast.Node {
	var n ast.Node

	switch t.Type {
	case RichTextTypeText:
		n = t.Text.Node()
	case RichTextTypeEquation:
		n = t.Equation.Node()
	case RichTextTypeMention:
		n = t.Mention.Node()
	default:
		log.Fatalf("invalid RichText of type %q", t.Type)
	}

	n = t.Annotations.wrapNode(n)

	return n
}

// Node returns the goldmark node of this notion text element.
func (t Text) Node() ast.Node {
	if t.Link != nil {
		return link(t.Content, t.Link.Url)
	}

	return stringNode(t.Content)
}

// Node returns the goldmark node of this notion mention element.
func (m Mention) Node() ast.Node {
	switch m.Type {
	case MentionTypeDatabase:
		return linkToPage(&m.Database.Id)
	case MentionTypeDate:
		return m.Date.Node()
	case MentionTypeUser:
		return m.User.Node()
	case MentionTypePage:
		return linkToPage(&m.Page.Id)
	// MentionTypeLinkPreview MentionType = "link_preview"
	// MentionTypePage        MentionType = "page"
	// MentionTypeUser        MentionType = "user"

	default:
		log.Fatalf("Node() for mention type %q not implemented yet", m.Type)
		return nil
	}
}

// Node returns the goldmark node of this notion user element.
func (u User) Node() ast.Node {
	switch u.Type {
	case UserTypePerson:
		return &GoldmarkUserPerson{
			ID:        u.Id,
			Name:      u.Name,
			AvatarURL: u.AvatarUrl,
			Person:    *u.Person,
		}
	case UserTypeBot:
		return &GoldmarkUserBot{
			ID:        u.Id,
			Name:      u.Name,
			AvatarURL: u.AvatarUrl,
			Bot:       *u.Bot,
		}
	default:
		log.Fatalf("invalid User of type %q", u.Type)
		return nil
	}
}

// Node returns the goldmark node of this notion element.
func (b Bookmark) Node() ast.Node {
	n := &GoldmarkBookmark{URL: b.Url}
	b.Caption.appendTo(n)
	return n
}

// Node returns the goldmark node of this notion element.
func (lp LinkPreview) Node() ast.Node {
	return &GoldmarkLinkPreview{URL: lp.Url}
}

func link(content, dest string) *ast.Link {
	n := ast.NewLink()
	n.Destination = []byte(dest)

	n.AppendChild(n, stringNode(content))

	return n
}

func linkToPage(id *UUID) ast.Node {
	// NOTE title to be filled by hugo
	return link("", fmt.Sprintf("/%s", id))
}

// Node returns the goldmark node of this notion element.
func (l LinkToPage) Node() ast.Node {
	switch l.Type {
	case LinkToPageTypeDatabaseId:
		return linkToPage(l.DatabaseId)
	case LinkToPageTypePageId:
		return linkToPage(l.PageId)
	default:
		log.Fatalf("invalid LinkToPage of type %q", l.Type)
		return nil
	}
}

// Node returns the goldmark node of this notion element.
// NOTE: not very sophisticated at the moment, we will change that
func (d Date) Node() ast.Node {
	s := d.Start.String()

	if d.End != nil {
		s += "-"
		s += d.End.String()
	}

	return stringNode(s)
}

// WrapNode wraps the given node in a node according to the annotations.
func (a Annotations) wrapNode(n ast.Node) ast.Node {
	if a.Code {
		n = wrapNode(ast.NewCodeSpan(), n)
	}

	if a.Bold {
		n = wrapNode(ast.NewEmphasis(2), n)
	}

	if a.Italic {
		n = wrapNode(ast.NewEmphasis(1), n)
	}

	if a.Strikethrough {
		n = wrapNode(_ast.NewStrikethrough(), n)
	}

	if a.Underline {
		n = wrapNode(ast.NewEmphasis(3), n)
	}

	n = a.Color.wrapNode(n)

	return n
}

// Node returns the node of the file.
// func (f File) node() ast.Node {
// 	switch f.Type {
// 	case FileTypeFile:
// 		return f.File.Node()
// 	case FileTypeExternal:
// 		return f.External.Node()
// 	default:
// 		log.Fatalf("invalid File of type %q", f.Type)
// 		return nil
// 	}
// }

// Node returns the node of the notion file.
// func (f NotionFile) Node() ast.Node {
// 	return &GoldmarkNotionFile{NotionFile: f}
// }

// Node returns the node of the notion file.
func (f ExternalFile) Node() ast.Node {
	return &GoldmarkExternalFile{URL: f.Url}
}

// Node returns the node of the list item.
func (l BulletedListItem) Node() ast.Node {
	n := ast.NewListItem(0)

	l.RichText.appendTo(n)

	for _, b := range l.Children.Nodes() {
		fmt.Println(b.Kind())
		n.AppendChild(n, b)
	}

	return l.Color.wrapNode(n)
}
