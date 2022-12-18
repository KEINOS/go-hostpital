package hostpital_test

import (
	"testing"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/stretchr/testify/require"
)

func TestTrimComment(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "no comment",
			input: "example.com",
			want:  "example.com",
		},
		{
			name:  "comment",
			input: "example.com#comment",
			want:  "example.com",
		},
		{
			name:  "comment with space",
			input: "example.com #comment",
			want:  "example.com",
		},
		{
			name:  "comment with space and tab",
			input: "example.com \t#comment",
			want:  "example.com",
		},
		{
			name:  "indented line with comment",
			input: "    example.com \t#comment",
			want:  "    example.com",
		},
	} {
		expect := test.want
		actual, err := hostpital.TrimComment(test.input)

		require.NoError(t, err, "test #%d: %s failed", index, test.name)
		require.Equal(t, expect, actual, "test #%d: %s failed", index, test.name)
	}
}

func TestTrimComment_contains_line_break(t *testing.T) {
	t.Parallel()

	result, err := hostpital.TrimComment("example.com\n")

	require.Error(t, err, "if the given line contains a line break, it should error")
	require.Contains(t, err.Error(), "line break found", "error shuld contain the reason")
	require.Equal(t, "", result, "returned value should be empty on error")
}
