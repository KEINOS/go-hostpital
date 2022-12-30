package hostpital

import (
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
)

// FindFile returns a list of file paths found under the given directory.
// It will not include directories even if it matches. If no files are found,
// it will return an fs.ErrNotExist error.
//
// e.g. FindFile("hosts*", "/home/user") will return:
//
//	["/home/user/.ssh/hosts", "/home/user/.ssh/hosts.deny", "/home/user/.ssh/hosts.allow"]
func FindFile(patternFile, pathDirSearch string) ([]string, error) {
	patternSearch := patternFile

	findList := []string{}
	err := filepath.WalkDir(pathDirSearch, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}

		matched, errMatch := filepath.Match(patternSearch, filepath.Base(path))
		if errMatch != nil {
			return errors.Wrap(errMatch, "failed filepath.Match")
		}

		if info.IsDir() || !matched {
			return nil
		}

		findList = append(findList, path)

		return nil
	})

	if err == nil && len(findList) == 0 {
		err = fs.ErrNotExist
	}

	return findList, errors.Wrap(err, "failed to search directory")
}
