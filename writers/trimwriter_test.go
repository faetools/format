package writers_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/faetools/format/writers"
)

var writersTestStrings = []string{
	`Nothing to trim`,

	``,

	"\n<html>\n  <head>\n    <title>Test</title>\n  </head>\n",

	"\n",

	`
One trimmed`,

	`
More trimmed


`,

	`

Trimmed


in between


`,
	`

Lorem ipsum dolor sit amet,
 consectetur adipiscing elit,
 sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
 Ut enim ad minim veniam,
 quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
 Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
 Excepteur sint occaecat cupidatat non proident,
 sunt in culpa qui officia deserunt mollit anim id est laborum.



`,
	`Excepteur sint occaecat cupidatat non proident,
 sunt in culpa qui officia deserunt mollit anim id est laborum.`,
}

func TestTrimwriter(t *testing.T) {
	t.Parallel()

	generalTests := make(testDefinitions, len(writersTestStrings))
	for i, s := range writersTestStrings {
		generalTests[i] = testDefinition{
			name:  fmt.Sprintf("Test %d", i),
			input: s,
		}
	}

	for _, cutset := range []string{"\n", "e", "asdf"} {
		cutset := cutset
		t.Run(fmt.Sprintf("for cutset %q", cutset), func(t *testing.T) {
			t.Parallel()

			t.Run("Trim", func(t *testing.T) {
				t.Parallel()

				testWriter(t, generalTests.asTable(func(s string) string {
					return strings.Trim(s, cutset)
				}), func(b io.Writer) writers.Writer {
					return writers.NewTrimWriter(b, cutset)
				})
			})

			t.Run("TrimLeft", func(t *testing.T) {
				t.Parallel()

				testWriter(t, generalTests.asTable(func(s string) string {
					return strings.TrimLeft(s, cutset)
				}), func(b io.Writer) writers.Writer {
					return writers.NewTrimLeftWriter(b, cutset)
				})
			})

			t.Run("TrimRight", func(t *testing.T) {
				t.Parallel()

				testWriter(t, generalTests.asTable(func(s string) string {
					return strings.TrimRight(s, cutset)
				}), func(b io.Writer) writers.Writer {
					return writers.NewTrimRightWriter(b, cutset)
				})
			})
		})
	}
}
