package markdown_test

import (
	_ "embed" // test file
	"encoding/json"
	"testing"

	"github.com/faetools/format/markdown"
	"github.com/faetools/go-notion/pkg/notion"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
)

//go:embed fab826ba-93ed-495e-b299-c9752fd2485b.json
var aboutBlocks []byte

func TestRender(t *testing.T) {
	t.Parallel()

	var blocks notion.Blocks
	assert.NoError(t, json.Unmarshal(aboutBlocks, &blocks))
	assert.Len(t, blocks, 11)

	doc := ast.NewDocument()

	for _, b := range blocks.Nodes() {
		doc.AppendChild(doc, b)
	}

	b, err := markdown.Render(nil, nil, doc, nil)
	assert.NoError(t, err)
	assert.Equal(t, want, string(b))
}

const want = `Hi, I'm Mark Rosemaker[^1].

[^1]: Or at least, that is the direct translation of my birth name Marco Rösler.

![](https://s3.us-west-2.amazonaws.com/secure.notion-static.com/97341fe7-1335-4d5f-96a7-1641dc094e6f/about.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=AKIAT73L2G45EIPT3X45%2F20220412%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20220412T152239Z&X-Amz-Expires=3600&X-Amz-Signature=833102e70cbab662c06ca3b82fc4b28b934c968bb5107a720e821596f196b753&X-Amz-SignedHeaders=host&x-id=GetObject)

Since 2015, I've been teaching German online with [Authentic German Learning](https://www.authenticgermanlearning.com/) and since 2016, I am doing it as a digital nomad[^2].

[^2]: A digital nomad is someone who works online and can, therefore, live a nomadic lifestyle and does so.

Besides that, I have many other interests. Furthermore, I now feel ready to teach all the other things I learned in the last couple of years.

So stay tuned if you want to find out more about

- how to start and run an online business
	- social media
	- SEO
- philosophy
	- rationality
	- ethics
- programming
	- website creation with Hugo / HTML / CSS
	- creating a chat bot with Go (golang)
- lifestyle
	- my journey through Europe
	- traveling on a budget
`
