package hostpital

import (
	"github.com/pkg/errors"
	"golang.org/x/net/idna"
)

// TransformToASCII converts the given hostName in UNICODE to ASCII/punycode.
//
// It does the opposite of TransformToUnicode().
func TransformToASCII(hostName string) (string, error) {
	hostASCII, err := idna.Lookup.ToASCII(hostName)

	return hostASCII, errors.Wrap(err, "failed to convert host name to ASCII")
}
