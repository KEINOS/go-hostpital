package hostpital

import "strings"

// TrimWordGaps reduces redundant and repetitive whitespace in the input string.
// It removes the line breaks, tabs, leading and trailing spaces as well.
func TrimWordGaps(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
