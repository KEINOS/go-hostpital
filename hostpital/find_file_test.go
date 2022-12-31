package hostpital

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindFile_bad_search_pattern(t *testing.T) {
	t.Parallel()

	paths, err := FindFile("[]a]", t.TempDir())

	require.Error(t, err, "it should error if the search pattern is mal-formed")
	require.Contains(t, err.Error(), "failed filepath.Match", "it should contain the wrapped error")
	require.Contains(t, err.Error(), "failed to search directory", "it should contain the error reason")
	require.Empty(t, paths, "it should return empty list on error")
}

func TestFindFile_no_files_matched(t *testing.T) {
	t.Parallel()

	paths, err := FindFile(`unknownfile`, t.TempDir())

	require.Error(t, err, "it should error if no files are found")
	require.ErrorIs(t, err, fs.ErrNotExist, "the error should be fs.ErrNotExist if no files are found")
	require.Empty(t, paths, "it should return empty list on error")
}

func TestFindFile_serch_dir_is_empty(t *testing.T) {
	t.Parallel()

	paths, err := FindFile("hosts*", "")

	require.Error(t, err, "it should error if the search path is empty")
	require.Contains(t, err.Error(), "failed filepath.WalkDir", "it should contain the wrapped error")
	require.Contains(t, err.Error(), "failed to search directory", "it should contain the error reason")
	require.Empty(t, paths, "it should return empty list on error")
}
