package writers_test

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/faetools/format/writers"
	"github.com/stretchr/testify/assert"
)

type (
	testDefinitions []testDefinition
	testDefinition  struct{ name, input string }

	testTable    []testTableRow
	testTableRow struct{ name, input, want string }
)

func (ds testDefinitions) asTable(f func(string) string) testTable {
	tt := make(testTable, len(ds))
	for i, d := range ds {
		tt[i] = testTableRow{
			name:  d.name,
			input: d.input,
			want:  f(d.input),
		}
	}

	return tt
}

type mode string

var (
	allWrite          mode = "Write"
	allWriteString    mode = "WriteString"
	allWriteByte      mode = "WriteByte"
	randomWrite       mode = "Write (randomly)"
	randomWriteString mode = "WriteString (randomly)"
	random            mode = "Completely Random"

	modes = []mode{
		allWrite,
		allWriteString,
		allWriteByte,
		randomWrite,
		randomWriteString,
		random,
	}
)

func testWrite(t *testing.T, w writers.Writer, s string, mode mode) (size int) {
	t.Helper()

	p := []byte(s)
	var err error

	switch mode {
	case allWrite:
		size, err = w.Write(p)
		assert.NoError(t, err)

		return size
	case allWriteString:
		size, err = w.WriteString(s)
		assert.NoError(t, err)

		return size
	case allWriteByte:
		for _, b := range p {
			assert.NoError(t, w.WriteByte(b))
		}

		return len(p)
	case random:
		for add := 0; len(p) > 0; {
			i := rand.Intn(len(p) + 1)
			switch rand.Intn(3) {
			case 0:
				add, err = w.Write(p[:i])
			case 1:
				add, err = w.WriteString(string(p[:i]))
			case 2:
				i = 1
				err = w.WriteByte(p[0])
				add = 1
			}

			assert.NoError(t, err)

			size += add
			p = p[i:]
		}

		return size
	default:
		for add := 0; len(p) > 0; {
			i := rand.Intn(len(p) + 1)

			if mode == randomWrite {
				add, err = w.Write(p[:i])
			} else {
				add, err = w.WriteString(string(p[:i]))
			}

			assert.NoError(t, err)

			size += add
			p = p[i:]
		}

		return size
	}
}

func testWriter(t *testing.T, tt testTable, genWriter func(io.Writer) writers.Writer) {
	t.Helper()

	for _, mode := range modes {
		mode := mode
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()

			for _, r := range tt {
				r := r
				t.Run(r.name, func(t *testing.T) {
					t.Parallel()

					b := &bytes.Buffer{}
					bqw := genWriter(b)

					size := testWrite(t, bqw, r.input, mode)
					assert.Equal(t, r.want, b.String(), "does not match wanted output")
					assert.Equal(t, len([]byte(r.input)), size, "mismatched size for %q", r.input)
				})
			}
		})
	}
}
