package hostpital

import (
	"net"
)

// IsIPAddress returns true if the given string is a valid IPv4 or IPv6 address.
// Note that white spaces are not allowed.
func IsIPAddress(hostName string) bool {
	return net.ParseIP(hostName) != nil
}
