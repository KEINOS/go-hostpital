package hostpital

import (
	"io/fs"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // do not parallelize due to temporary changing global variables
func TestFileExists_lstat_fail(t *testing.T) {
	oldOsLstat := osLstat
	defer func() { osLstat = oldOsLstat }()

	osLstat = func(name string) (fs.FileInfo, error) {
		return nil, errors.New("forced error")
	}

	result := FileExists("file_exists_test.go")

	require.False(t, result, "it should return false if os.Lstat fails")
}
