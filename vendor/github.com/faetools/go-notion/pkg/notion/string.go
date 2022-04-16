package notion

import (
	"fmt"
	"strings"
)

func (b Block) Title() string {
	switch b.Type {
	case BlockTypeChildDatabase:
		return b.ChildDatabase.Title
	case BlockTypeChildPage:
		return b.ChildPage.Title
	default:
		return fmt.Sprintf("<no title defined for block type %q>", b.Type)
	}
}

func (ts RichTexts) Raw() string {
	s := make([]string, len(ts))
	for i, t := range ts {
		s[i] = t.Raw()
	}

	return strings.Join(s, "\n")
}

func (id *UUID) String() string {
	if id == nil {
		return "<no uuid>"
	}

	return string(*id)
}

// func (t Text) String() string { return t.Content }

func (t RichText) Raw() string {
	switch t.Type {
	case RichTextTypeText:
		return t.Text.Content
	case RichTextTypeMention:
		return fmt.Sprintf("%#v", *t.Mention)
	case RichTextTypeEquation:
		return fmt.Sprintf("%#v", *t.Equation)
	default:
		return fmt.Sprintf("invalid RichText type %q", t.Type)
	}
}
