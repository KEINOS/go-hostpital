package hostpital

import "strings"

// TrimIPAdd removes the leading IP address from the given string.
//
// Note that this function expects the given string to be a line from a hosts
// file. e.g. "0.0.0.0 example.com".
func TrimIPAdd(line string) string {
	line = TrimWordGaps(line)

	ipAdd := strings.Split(line, " ")[0]
	if IsIPAddress(ipAdd) {
		trimmed := strings.TrimLeft(strings.TrimLeft(line, ipAdd), Cutset)

		return TrimIPAdd(trimmed)
	}

	return line
}
