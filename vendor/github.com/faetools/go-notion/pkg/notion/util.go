package notion

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// filterBy returns only blocks that match the filter.
func (bs Blocks) filterBy(tp BlockType) Blocks {
	out := Blocks{}

	for _, b := range bs {
		if b.Type == tp {
			out = append(out, b)
		}
	}

	return out
}

// ChildPages returns all child pages.
func (bs Blocks) ChildPages() Blocks { return bs.filterBy(BlockTypeChildPage) }

// ChildDatabases returns all child databases.
func (bs Blocks) ChildDatabases() Blocks { return bs.filterBy(BlockTypeChildDatabase) }

// NewParagraph constructs a new block that is a valid paragraph.
func NewParagraph(content string) Block {
	return Block{
		Object: "block",
		Id:     UUID(uuid.NewString()),
		Type:   BlockTypeParagraph,
		Paragraph: &Paragraph{
			RichText: NewRichTexts(content),
			Color:    ColorDefault,
			Children: Blocks{},
		},
	}
}

// NewRichText creates a RichText object with the desired content.
func NewRichText(content string) RichText {
	return RichText{
		Type:        RichTextTypeText,
		PlainText:   content,
		Text:        &Text{Content: content},
		Annotations: Annotations{Color: ColorDefault},
	}
}

// NewRichTexts creates a RichTexts object with the desired content.
func NewRichTexts(content string) RichTexts {
	return RichTexts{NewRichText(content)}
}

// TitleProperty represents a Title property.
var TitleProperty = PropertyMeta{
	Id:    "title", // must be this
	Name:  "Title",
	Type:  PropertyTypeTitle,
	Title: &map[string]interface{}{},
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s: %s - %s", e.Status, http.StatusText(e.Status), e.Code, e.Message)
}

// Title returns the page title.
func (p Page) Title() string { return p.GetProperties().Title() }

// TitleWithEmoji returns the page title, prepended by an emoji if present.
func (p Page) TitleWithEmoji() string {
	if p.Icon != nil && p.Icon.Type == IconTypeEmoji {
		return fmt.Sprintf("%s %s", *p.Icon.Emoji, p.Title())
	}

	return p.Title()
}

// PropertyValueMap is a map of all property values.
type PropertyValueMap map[string]PropertyValue

// Title returns the title of the page.
func (props PropertyValueMap) Title() string {
	return props.title().Title.Raw()
}

func (props PropertyValueMap) title() PropertyValue {
	for _, prop := range props {
		if prop.Title != nil {
			return prop
		}
	}

	log.Fatal("could not find title property which is mandatory")

	return PropertyValue{}
}

// GetProperties returns a map of all properties.
func (p *Page) GetProperties() PropertyValueMap {
	if props, ok := p.Properties.(PropertyValueMap); ok {
		return props
	}

	props := PropertyValueMap{}

	if err := p.unmarshalProperties(&props); err != nil {
		log.Fatal(err)
	}

	p.Properties = props

	return props
}

// unmarshalProperties decodes the properties into the given interface.
func (p Page) unmarshalProperties(v interface{}) error {
	b, err := json.Marshal(p.Properties)
	if err != nil {
		return fmt.Errorf("marshalling properties: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("unmarshalling %s into %T: %w", string(b), v, err)
	}

	return nil
}

// unmarshalProperties decodes the properties into the given interface.
func (db Database) unmarshalProperties(v interface{}) error {
	b, err := json.Marshal(db.Properties)
	if err != nil {
		return fmt.Errorf("marshalling properties: %w", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("unmarshalling %s into %T: %w", string(b), v, err)
	}

	return nil
}

// PropertyMetaMap is a map of all properties.
type PropertyMetaMap map[string]PropertyMeta

// GetProperties returns a map of all properties.
func (db *Database) GetProperties() PropertyMetaMap {
	if props, ok := db.Properties.(PropertyMetaMap); ok {
		return props
	}

	props := PropertyMetaMap{}

	if err := db.unmarshalProperties(&props); err != nil {
		log.Fatal(err)
	}

	db.Properties = props

	return props
}

// Title returns the title property.
func (m PropertyMetaMap) Title() PropertyMeta {
	for _, v := range m {
		if v.Type == PropertyTypeTitle {
			return v
		}
	}

	return PropertyMeta{}
}

func (m PropertyMetaMap) findSameRelation(rel *RelationConfiguration) PropertyMeta {
	for _, v := range m {
		if v.Type == PropertyTypeRelation &&
			v.Relation.DatabaseId == rel.DatabaseId &&
			v.Relation.SyncedPropertyName == rel.SyncedPropertyName {
			return v
		}
	}

	return PropertyMeta{}
}

// Merge merges a PropertyMeta with another.
// NOTE: For now, this is just the name.
func (m PropertyMeta) Merge(other PropertyMeta) PropertyMeta {
	m.Name = other.Name
	return m
}

// Merge merges a PropertyMetaMap with another.
func (m PropertyMetaMap) Merge(other PropertyMetaMap) (PropertyMetaMap, bool) {
	changed := false

	for k, v := range other {
		// we get only title key from notion
		k := strings.Title(k)

		if _, ok := m[k]; ok {
			// we are not yet implementing the functionality to change a property
			// that already has the same key/name
			// instead, create new properties and/or change property names
			// this way, we prevent accidental deletion of data
			continue
		}

		// we may update the name of the property
		updateName := func(original PropertyMeta) {
			k = original.Name // we *cannot* assign to a new key, use the original one (same as Name)
			v = original.Merge(v)
		}

		switch v.Type {
		case PropertyTypeTitle:
			updateName(m.Title())
		case PropertyTypeRelation:
			if rel := m.findSameRelation(v.Relation); rel.Name != "" {
				updateName(rel)
			}
		}

		// if updateName was not called, we are creating a new property
		// otherwise, we update an existing one with the new name
		m[k] = v
		changed = true
	}

	return m, changed
}

// Get returns the block with the given ID or nil, if it wasn't found.
func (bs Blocks) Get(id UUID) *Block {
	for _, b := range bs {
		if b.Id == id {
			return &b
		}
	}

	return nil
}

// GetRelations returns UUIDs of all objects this value is related to.
func (v PropertyValue) GetRelations() []UUID {
	if v.Relation == nil || len(*v.Relation) == 0 {
		return nil
	}

	rels := make([]UUID, len(*v.Relation))

	for i, rel := range *v.Relation {
		rels[i] = rel.Id
	}

	return rels
}

// GetMultiSelection returns names of all objects that was selected.
func (v PropertyValue) GetMultiSelection() []string {
	if v.MultiSelect == nil || len(*v.MultiSelect) == 0 {
		return nil
	}

	selections := make([]string, len(*v.MultiSelect))

	for i, sel := range *v.MultiSelect {
		selections[i] = sel.Name
	}

	return selections
}

// GetSelection returns name of the object that was selected.
func (v PropertyValue) GetSelection() string {
	if v.Select == nil {
		return ""
	}

	return v.Select.Name
}

// GetRawText returns the raw text of the property.
func (v PropertyValue) GetRawText() string {
	if v.RichText == nil {
		return ""
	}

	return v.RichText.Raw()
}

// GetBool returns the checkbox value.
func (v PropertyValue) GetBool() bool {
	return v.Checkbox != nil && *v.Checkbox
}

// GetNumber returns the number.
func (v PropertyValue) GetNumber() float32 {
	if v.Number == nil {
		return 0
	}

	return *v.Number
}

// URL return the URL of the file
func (f File) URL() string {
	switch f.Type {
	case FileTypeExternal:
		return f.External.Url
	case FileTypeFile:
		return f.File.Url
	default:
		log.Fatalf("invalid File of type %q", f.Type)
		return ""
	}
}

// IsExternal returns wheather or not the file is external.
func (f File) IsExternal() bool {
	return f.External != nil
}

// IsInternal returns wheather or not the file is internal.
func (f File) IsInternal() bool {
	return f.File != nil
}
