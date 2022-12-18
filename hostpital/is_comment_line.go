package hostpital

import "strings"

// IsCommentLine returns true if the given line is a comment line.
func IsCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), string(DelimComnt))
}
