package hostpital

import (
	"golang.org/x/net/idna"
)

// IsCompatibleIDNA2008 returns true if the given hostName is properly formatted for
// registration and is in ASCII/Punycode.
// False if it contains hosts or labels in Unicode or not IDNA2008 compatible.
//
// Use TransformASCII() to convert raw punycode to IDNA2008 compatible ASCII.
//
// Note that this function returns false if the label (host name or part of the
// subdomain) contains "_" (underscore).
func IsCompatibleIDNA2008(hostName string) bool {
	hostNameASCII, err := idna.Registration.ToASCII(hostName)
	if err != nil {
		return false
	}

	return hostNameASCII == hostName
}
