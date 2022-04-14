package yaml_test

import (
	"fmt"
	"testing"

	"github.com/faetools/format/yaml"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	for i, tt := range []struct {
		in, out string
	}{
		{},
		{`foo: "bar"`, "foo: bar\n"},
		{`foo: 'bar'`, "foo: bar\n"},
		{`foo: "3"`, "foo: '3'\n"},
		{`foo: 3`, "foo: 3\n"},
		{`foo: "foo\nblub"`, "foo: |-\n  foo\n  blub\n"},
		{`foo: "[]"`, "foo: '[]'\n"},
		{`# a comment
foo: "foo\nblub"`, `# a comment
foo: |-
  foo
  blub
`},
		{`# a sequence
foo:
- a
- b`, `# a sequence
foo:
  - a
  - b
`},
		{`# an anchor
defaults: &resources
  assets_expire_days: '[]'

service_name_prefix: production
<<: *resources`, `# an anchor
defaults: &resources
  assets_expire_days: '[]'
service_name_prefix: production
<<: *resources
`},
	} {
		i, tt := i, tt
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			t.Parallel()

			res, err := yaml.Format([]byte(tt.in))
			require.NoError(t, err)

			require.Equal(t, tt.out, string(res))
		})
	}
}

func TestFormat_Error(t *testing.T) {
	t.Parallel()

	_, err := yaml.Format([]byte(`[`))
	require.EqualError(t, err, "unmarshalling: yaml: line 1: did not find expected node content")
}
