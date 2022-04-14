package writers_test

import (
	"testing"

	"github.com/faetools/format/writers"
)

var bqTestTable = testTable{
	{
		`Empty`,
		``,
		``,
	},
	{
		`One Line`,
		`Foo`,
		`> Foo`,
	},
	{
		`Two Lines`,
		`Foo
Bar`,
		`> Foo
> Bar`,
	},
	{
		`Empty Lines`,
		`
Foo


Bar

`,
		`>
> Foo
>
>
> Bar
>
>`,
	},
	{
		`Nested Quote`,
		`> This is already a quote
Foo


> Another Quote

`,
		`>> This is already a quote
> Foo
>
>
>> Another Quote
>
>`,
	},
}

func TestBlockquotewriter(t *testing.T) {
	t.Parallel()

	testWriter(t, bqTestTable, writers.NewBlockquoteWriter)
}
