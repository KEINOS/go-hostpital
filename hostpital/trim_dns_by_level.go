package hostpital

import "strings"

// TrimDNSByLevel returns the host name with the specified number of DNS levels.
//
// E.g.
//
//	TrimDNSByLevel("www.example.com", 0) returns "com". Such as the top level domain.
//	TrimDNSByLevel("www.example.com", 1) returns "example.com". Such as the second level domain.
//	TrimDNSByLevel("www.example.com", 5) returns "www.example.com". As is.
func TrimDNSByLevel(host string, level int) string {
	parsed := strings.Split(host, ".")
	if len(parsed) <= level {
		return host
	}

	return strings.Join(parsed[len(parsed)-level-1:], ".")
}
