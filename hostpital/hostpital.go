/*
Package hostpital provides functions to assist in the maintenance of hosts files.
*/
package hostpital

import "os"

const (
	// DelimComnt is the delimiter for comments.
	DelimComnt = '#'
	// DelimDNS is the delimiter for DNS labels.
	DelimDNS = '.'
	// LF is the line feed, "\n".
	LF = int32(0x0a)
	// CR is the carriage return, "\r".
	CR = int32(0x0d)
	// Cutset is the set of characters for trimming white spaces.
	Cutset = "\t\n\v\f\r "
)

// Function variables for testing.
//
//nolint:gochecknoglobals // Allow global var for testing
var (
	// osOpen is a copy of os.Open to ease testing by mocking/monkey patching.
	osOpen = os.Open
	// osLstat is a copy of os.Lstat to ease testing.
	osLstat = os.Lstat
)
