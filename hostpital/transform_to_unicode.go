package hostpital

import (
	"github.com/pkg/errors"
	"golang.org/x/net/idna"
)

// TransformToUnicode converts the given punycoded host name in ASCII to UNICODE
// format.
//
// It does the opposite of TransformToASCII().
func TransformToUnicode(hostASCII string) (string, error) {
	hostPunycode, err := idna.Lookup.ToUnicode(hostASCII)

	return hostPunycode, errors.Wrap(err, "failed to convert host name to punycode")
}
