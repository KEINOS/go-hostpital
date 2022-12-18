package hostpital

import (
	"errors"
	"io/fs"
	"os"
)

// IsExistingFile returns true if the path is an existing file.
// It returns false if the path is a directory or does not exist.
func IsExistingFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return !errors.Is(err, fs.ErrNotExist)
	}

	return !info.IsDir()
}
