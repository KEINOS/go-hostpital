package hostpital

import (
	"os"
)

// FileExists returns true if the file exists and is not a directory.
func FileExists(path string) bool {
	info, err := osLstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // unable to find file
		}

		return false // error when running os.Lstat
	}

	if info.IsDir() {
		return false // path is a directory
	}

	return true
}
