package markdown_test

import (
	"testing"

	"github.com/faetools/format/markdown"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type (
	testTable    []testTableRow
	testTableRow struct{ name, input, want string }
)

var renderTestTable = testTable{
	{
		`Empty`,
		``,
		`
`,
	},
	{
		`Alternative Heading 1`,
		`Test File
===============`,
		`# Test File
`,
	},
	{
		`Alternative Heading 2`,
		`Heading level 2
---------------`,
		`## Heading level 2
`,
	},
	{
		`Paragraphs`,
		`
I really like using Markdown.
Very much so!

I think I'll use it to format all of my documents from now on.`,
		`I really like using Markdown.
Very much so!

I think I'll use it to format all of my documents from now on.
`,
	},
	{
		`No blank lines`,
		`Without blank lines, this might not look right.
# Heading
Don't do this!`,
		`Without blank lines, this might not look right.

# Heading

Don't do this!
`,
	},
	{
		`spaces in front of paragraphs`,
		`  This can result in unexpected formatting problems.

  Don't add tabs or spaces in front of paragraphs.`,
		`This can result in unexpected formatting problems.

Don't add tabs or spaces in front of paragraphs.
`,
	},
	{
		`HTML Break`,
		`First line with the HTML tag after.<br>
And the next line.`,
		`First line with the HTML tag after.<br>
And the next line.
`,
	},
	{
		`Bold Text`,
		`I just love **bold text**.

I just love __bold text__.

Love**is**bold`,
		`I just love **bold text**.

I just love **bold text**.

Love**is**bold
`,
	},
	{
		`More Emphasis`,
		`Italicized text is the _cat's meow_.

This text is ***really important***.

This text is ___really important___.

This text is __*really important*__.

This text is **_really important_**.

This is really***very***important text.`,
		`Italicized text is the *cat's meow*.

This text is ***really important***.

This text is ***really important***.

This text is ***really important***.

This text is ***really important***.

This is really***very***important text.
`,
	},
	{
		`Simple Quote`,
		`> Dorothy followed her through many of the beautiful rooms in her castle.`,
		`> Dorothy followed her through many of the beautiful rooms in her castle.
`,
	},
	{
		`Blockquote`,
		`Here is a blockquote:

> Dorothy followed her through many of the beautiful rooms in her castle.
>
> The Witch bade her clean the pots and kettles and sweep the floor and keep the fire fed with wood

This concludes the blockquote.`,
		`Here is a blockquote:

> Dorothy followed her through many of the beautiful rooms in her castle.
>
> The Witch bade her clean the pots and kettles and sweep the floor and keep the fire fed with wood

This concludes the blockquote.
`,
	},
	{
		`Nested Blockquotes`,
		`> Dorothy followed her through many of the beautiful rooms in her castle.
>
>> The Witch bade her clean the pots and kettles and sweep the floor and keep the fire fed with wood.`,
		`> Dorothy followed her through many of the beautiful rooms in her castle.
>
>> The Witch bade her clean the pots and kettles and sweep the floor and keep the fire fed with wood.
`,
	},
	{
		`Blockquote with formatted elements`,
		`> #### The quarterly results look great!
>
> - Revenue was off the chart.
> - Profits were higher than ever.
>
>  _Everything_ is going according to __plan__.`,
		`> #### The quarterly results look great!
>
> - Revenue was off the chart.
> - Profits were higher than ever.
>
> *Everything* is going according to **plan**.
`,
	},
	{
		`Ordered List with All Ones`,
		`1. First item
1. Second item
1. Third item
1. Fourth item`,
		`1. First item
2. Second item
3. Third item
4. Fourth item
`,
	},
	{
		`Ordered List with Wrong Numbers`,
		`1. First item
8. Second item
3. Third item
5. Fourth item`,
		`1. First item
2. Second item
3. Third item
4. Fourth item
`,
	},

	{
		`Ordered List with Blank Item`,
		`1. First item
8. Second item
3.
5. Fourth item`,
		`1. First item
2. Second item
3. ` + // Prevent this from space from getting removed.
			`
4. Fourth item
`,
	},
	{
		`Ordered List with Indent`,
		`6. First item
7. Second item
8. Third item
    1. First indented item
    5. Second indented item
9.  Fourth item`,
		`1. First item
2. Second item
3. Third item
	1. First indented item
	2. Second indented item
4. Fourth item
`,
	},
	{
		`Unordered List with Indent`,
		`- First item
- Second item
- Third item
    - First indented item
    - Second indented item
- Fourth item`,
		`- First item
- Second item
- Third item
	- First indented item
	- Second indented item
- Fourth item
`,
	},
	{
		`Unordered List with Blank Item`,
		`- First item
- Second item
-
- Fourth item`,
		`- First item
- Second item
- ` + // Prevent this from space from getting removed.
			`
- Fourth item
`,
	},
	{
		`Unordered List with Plus Sign`,
		`+ First item
+ Second item
+ Third item
+ Fourth item`,
		`- First item
- Second item
- Third item
- Fourth item
`,
	},
	{
		`Unordered List with Asterisk`,
		`* First item
* Second item
* Third item
* Fourth item`,
		`- First item
- Second item
- Third item
- Fourth item
`,
	},
	{
		`Unordered List Items with Numbers`,
		`- 1968\. A great year!
- I think 1969 was second best.`,
		`- 1968\. A great year!
- I think 1969 was second best.
`,
	},
	{
		`Nested Paragraph in List`,
		`* This is the first list item.
* Here's the second list item.

    I need to add another paragraph below the second list item.
    And it goes on for more than one line.

* And here's the third list item.`,
		`- This is the first list item.
- Here's the second list item.

	I need to add another paragraph below the second list item.
	And it goes on for more than one line.

- And here's the third list item.
`,
	},
	{
		`Nested Blockquote in List`,
		`* This is the first list item.
* Here's the second list item.

    > A blockquote would look great below the second list item.
    > Indeed, it would!

* And here's the third list item.`,
		`- This is the first list item.
- Here's the second list item.

	> A blockquote would look great below the second list item.
	> Indeed, it would!

- And here's the third list item.
`,
	},
	{
		`Code Block`,
		`    <html>
      <head>
        <title>Test</title>
      </head>`,
		"```\n<html>\n  <head>\n    <title>Test</title>\n  </head>\n```\n",
	},
	{
		`Formatted Code Block`,
		"```go\npackage foo\nvar v int =     3\n```\n",
		"```go\npackage foo\n\nvar v int = 3\n```\n",
	},
	{
		// Code blocks are normally indented four spaces or one tab. When they’re in a list, indent them eight spaces or two tabs.
		`Code Block in List`,
		`1. Open the file.
2. Find the following code block on line 21:

	        <html>
	          <head>
	            <title>Test</title>
	          </head>

3. Update the title to match the name of your website.`,
		`1. Open the file.
2. Find the following code block on line 21:

		<html>
		  <head>
		    <title>Test</title>
		  </head>

3. Update the title to match the name of your website.
`,
	},
	{
		`Code Block with Tabs in List`,
		"1. Open the file.\n2. Find the following code block on line 21:\n\n\t        <html>\n\t\t\t  <head>\n\t            <title>Test</title>\n\t          </head>\n\n3. Update the title to match the name of your website.",
		`1. Open the file.
2. Find the following code block on line 21:

		<html>
		  <head>
		    <title>Test</title>
		  </head>

3. Update the title to match the name of your website.
`,
	},
	{
		`Image`,
		`![The San Juan Mountains are beautiful!](/assets/images/san-juan-mountains.jpg "San Juan Mountains")`,
		`![The San Juan Mountains are beautiful!](/assets/images/san-juan-mountains.jpg "San Juan Mountains")
`,
	},
	{
		`Thematic Breaks / Horizontal Rules`,
		`the following are breaks

***

---

_________________`,
		`the following are breaks

---

---

---
`,
	},
	{
		`Thematic Breaks / Horizontal Rules #2`,
		`Try to put a blank line before...

---

...and after a horizontal rule.`,
		`Try to put a blank line before...

---

...and after a horizontal rule.
`,
	},
	{
		`Link`,
		`My favorite search engine is [Duck Duck Go](https://duckduckgo.com).
`,
		`My favorite search engine is [Duck Duck Go](https://duckduckgo.com).
`,
	},
	{
		`Auto Link`,
		`<https://www.markdownguide.org>`,
		`<https://www.markdownguide.org>
`,
	},
	{
		`Auto E-Mail`,
		`<fake@example.com>`,
		`<fake@example.com>
`,
	},
	{
		`Formatted Links`,
		`I love supporting the __[EFF](https://eff.org)__.

This is the _[Markdown Guide](https://www.markdownguide.org)_.

See the section on [` + "`code`" + `](#code).`,
		`I love supporting the **[EFF](https://eff.org)**.

This is the *[Markdown Guide](https://www.markdownguide.org)*.

See the section on [` + "`code`" + `](#code).
`,
	},
	{
		`Reference-Style Links`,
		`In a hole in the ground there lived a hobbit. Not a nasty, dirty, wet hole, filled with the ends
of worms and an oozy smell, nor yet a dry, bare, sandy hole with nothing in it to sit down on or to
eat: it was a [hobbit-hole][1], and that means comfort.

[1]: <https://en.wikipedia.org/wiki/Hobbit#Lifestyle> "Hobbit lifestyles"`,
		`In a hole in the ground there lived a hobbit. Not a nasty, dirty, wet hole, filled with the ends
of worms and an oozy smell, nor yet a dry, bare, sandy hole with nothing in it to sit down on or to
eat: it was a [hobbit-hole](https://en.wikipedia.org/wiki/Hobbit#Lifestyle "Hobbit lifestyles"), and that means comfort.
`,
	},
	{
		`Linked Image`,
		`[![An old rock in the desert](/assets/images/shiprock.jpg "Shiprock, New Mexico by Beau Rogers")](https://www.flickr.com/photos/beaurogers/31833779864/in/photolist-Qv3rFw-34mt9F-a9Cmfy-5Ha3Zi-9msKdv-o3hgjr-hWpUte-4WMsJ1-KUQ8N-deshUb-vssBD-6CQci6-8AFCiD-zsJWT-nNfsgB-dPDwZJ-bn9JGn-5HtSXY-6CUhAL-a4UTXB-ugPum-KUPSo-fBLNm-6CUmpy-4WMsc9-8a7D3T-83KJev-6CQ2bK-nNusHJ-a78rQH-nw3NvT-7aq2qf-8wwBso-3nNceh-ugSKP-4mh4kh-bbeeqH-a7biME-q3PtTf-brFpgb-cg38zw-bXMZc-nJPELD-f58Lmo-bXMYG-bz8AAi-bxNtNT-bXMYi-bXMY6-bXMYv)`,
		`[![An old rock in the desert](/assets/images/shiprock.jpg "Shiprock, New Mexico by Beau Rogers")](https://www.flickr.com/photos/beaurogers/31833779864/in/photolist-Qv3rFw-34mt9F-a9Cmfy-5Ha3Zi-9msKdv-o3hgjr-hWpUte-4WMsJ1-KUQ8N-deshUb-vssBD-6CQci6-8AFCiD-zsJWT-nNfsgB-dPDwZJ-bn9JGn-5HtSXY-6CUhAL-a4UTXB-ugPum-KUPSo-fBLNm-6CUmpy-4WMsc9-8a7D3T-83KJev-6CQ2bK-nNusHJ-a78rQH-nw3NvT-7aq2qf-8wwBso-3nNceh-ugSKP-4mh4kh-bbeeqH-a7biME-q3PtTf-brFpgb-cg38zw-bXMZc-nJPELD-f58Lmo-bXMYG-bz8AAi-bxNtNT-bXMYi-bXMY6-bXMYv)
`,
	},
	{
		`Escaping Characters`,
		`\* Without the backslash, this would be a bullet in an unordered list.`,
		`\* Without the backslash, this would be a bullet in an unordered list.
`,
	},
	{
		`HTML`,
		`This **word** is bold. This <em>word</em> is italic.`,
		`This **word** is bold. This <em>word</em> is italic.
`,
	},
	{
		`Image in List`,
		`1. Open the file containing the Linux mascot.
2. Marvel at its beauty.

    ![Tux, the Linux mascot](/assets/images/tux.png)

1. Close the file.`,
		`1. Open the file containing the Linux mascot.
2. Marvel at its beauty.

	![Tux, the Linux mascot](/assets/images/tux.png)

3. Close the file.
`,
	},
	{
		`Image surrounded by paragraphs`,
		`some text

![Tux, the Linux mascot](/assets/images/tux.png)

some more text
`,
		`some text

![Tux, the Linux mascot](/assets/images/tux.png)

some more text
`,
	},
	{
		`Unordered List in Ordered List`,
		`1. First item
2. Second item
3. Third item
    - Indented item
    - Indented item
4. Fourth item`,
		`1. First item
2. Second item
3. Third item
	- Indented item
	- Indented item
4. Fourth item
`,
	},
	{
		`Code and Backticks need escaping`,
		"At the command prompt, type `nano`.\n\n``Use `code` in your Markdown file.``",
		"At the command prompt, type `nano`.\n\n``Use `code` in your Markdown file.``\n",
	},
	{
		`HTML Block`,
		`There is about to be an HTML block.


<html>
		<p>Foo</p>
</html>


This is the end of the HTML block.
`,
		"There is about to be an HTML block.\n\n<html>\n\t\t<p>Foo</p>\n</html>\n\nThis is the end of the HTML block.\n",
	},
	{
		`One Line HTML Comment`,
		`What follows is a comment.

<!-- asf -->

This is the end of the HTML comment.
`,
		`What follows is a comment.

<!-- asf -->

This is the end of the HTML comment.
`,
	},
	{
		`Three Line HTML Comment`,
		`What follows is a comment.

<!--
asf
-->

This is the end of the HTML comment.
`,
		`What follows is a comment.

<!--
asf
-->

This is the end of the HTML comment.
`,
	},
	{
		`HTML Long Comment`,
		`What follows is a comment.

<!--

## Purpose
A short, one to two sentences max description of the service. What is the importance of this service? What's its core functionality?

Last part of the comment.
--->

This is the end of the HTML comment.
`,
		`What follows is a comment.

<!--

## Purpose
A short, one to two sentences max description of the service. What is the importance of this service? What's its core functionality?

Last part of the comment.
--->

This is the end of the HTML comment.
`,
	},
	{
		`YAML Frontmatter`,
		`---
draft: true
---

some text
`,
		`---
draft: true
---

some text
`,
	},

	{
		`YAML Frontmatter #2`,
		`---
title: Solving Problems, One at a Time
notionID: 37973bbb-0dc7-4c47-95c6-4fcc1c238c72
---

some text
`,
		`---
title: Solving Problems, One at a Time
notionID: 37973bbb-0dc7-4c47-95c6-4fcc1c238c72
---

some text
`,
	},
}

func TestRendering(t *testing.T) {
	t.Parallel()

	for _, tt := range renderTestTable {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := markdown.Format([]byte(tt.input))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(res), "does not match wanted output")

			secondResult, err := markdown.Format(res)
			assert.NoError(t, err)
			assert.Equal(t, string(res), string(secondResult), "is not idempotent")
		})
	}
}

func TestRendering_Terminal(t *testing.T) {
	t.Parallel()

	table := testTable{
		{
			"Bold Text",
			"I just love **bold text**.",
			"I just love \x1b[1mbold text\x1b[0m.\n",
		},
		{
			"Italic Text",
			"I just love *italic text*.",
			"I just love \x1b[3mitalic text\x1b[0m.\n",
		},
		{
			"Strikethrough Text", // To be done later.
			"I ~like~ love strikethrough.",
			"I ~like~ love strikethrough.\n",
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := markdown.Format([]byte(tt.input), markdown.WithTerminal())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(res), "does not match wanted output")

			secondResult, err := markdown.Format(res)
			assert.NoError(t, err)
			assert.Equal(t, string(res), string(secondResult), "is not idempotent")
		})
	}
}

var customKind = ast.NewNodeKind("customKind")

type custom struct {
	*ast.String
}

func (c custom) Kind() ast.NodeKind { return customKind }

type customRenderer struct{}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *customRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(customKind, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			w.WriteString("</custom>")
			return ast.WalkContinue, nil
		}

		w.WriteString("<custom>")
		w.Write(n.(*custom).Value)

		return ast.WalkContinue, nil
	})
}

func wrapNode(parent, child ast.Node) ast.Node {
	parent.AppendChild(parent, child)
	return parent
}

func TestCustomKind(t *testing.T) {
	t.Parallel()

	n := wrapNode(ast.NewParagraph(), &custom{String: ast.NewString([]byte("foo"))})
	n = wrapNode(ast.NewParagraph(), n)
	n = wrapNode(ast.NewDocument(), n)

	b, err := markdown.Render(nil, nil, n, renderer.WithNodeRenderers(
		util.Prioritized(&customRenderer{}, 1),
	))
	assert.NoError(t, err)
	assert.Equal(t, "<custom>foo</custom>\n", string(b))
}
