package writers_test

import (
	"io"
	"testing"

	"github.com/faetools/format/writers"
)

var iwTestTable = testTable{
	{
		`Empty`,
		``,
		``,
	},
	{
		`One Line`,
		`Foo`,
		`		Foo`,
	},
	{
		`Two Lines`,
		`Foo
Bar`,
		`		Foo
		Bar`,
	},
	{
		`Empty Lines`,
		`
Foo


Bar

`,
		`
		Foo


		Bar

`,
	},
	{
		`Nested Indent`,
		`	This is already indented
Foo


	Another indent

`,
		`			This is already indented
		Foo


			Another indent

`,
	},
}

func TestIndentWriter(t *testing.T) {
	t.Parallel()

	testWriter(t, iwTestTable, func(b io.Writer) writers.Writer {
		return writers.NewIndentWriter(b, 2)
	})
}
