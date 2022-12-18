package hostpital

import (
	"strings"

	"github.com/pkg/errors"
)

// TrimComment removes hash ("#") comments from the given string and trailing
// spaces.
//
// This function expects the given string to be a line from a hosts file and
// will error if the given line contains a line break.
func TrimComment(line string) (string, error) {
	result := []rune{}

	for _, char := range line {
		if char == LF {
			return "", errors.New("line break found")
		}

		if char == DelimComnt {
			break // comment delimiter found
		}

		result = append(result, char)
	}

	return strings.TrimRight(string(result), Cutset), nil
}
