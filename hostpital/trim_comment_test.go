package hostpital_test

import (
	"testing"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/stretchr/testify/require"
)

const hostTrimCommentExampleCom = "example.com"

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
			input: hostTrimCommentExampleCom,
			want:  hostTrimCommentExampleCom,
		},
		{
			name:  "comment",
			input: hostTrimCommentExampleCom + "#comment",
			want:  hostTrimCommentExampleCom,
		},
		{
			name:  "comment with space",
			input: hostTrimCommentExampleCom + " #comment",
			want:  hostTrimCommentExampleCom,
		},
		{
			name:  "comment with space and tab",
			input: hostTrimCommentExampleCom + " \t#comment",
			want:  hostTrimCommentExampleCom,
		},
		{
			name:  "indented line with comment",
			input: "    " + hostTrimCommentExampleCom + " \t#comment",
			want:  "    " + hostTrimCommentExampleCom,
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

	result, err := hostpital.TrimComment(hostTrimCommentExampleCom + "\n")

	require.Error(t, err, "if the given line contains a line break, it should error")
	require.Contains(t, err.Error(), "line break found", "error shuld contain the reason")
	require.Empty(t, result, "returned value should be empty on error")
}
