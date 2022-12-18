package hostpital_test

import (
	"testing"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/stretchr/testify/require"
)

func TestIsIPAddress(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		input  string
		expect bool
	}{
		// Golden cases
		{"0.0.0.0", true}, // undefined address
		{"0:0:0:0:0:0:0:0", true},
		{"::", true},
		{"127.0.0.1", true}, // localhost, loopback
		{"0:0:0:0:0:0:0:1", true},
		{"::1", true},
		{"fe00::0", true},          // ip6-localnet
		{"ff00::0", true},          // ip6-mcastprefix
		{"ff02::1", true},          // ip6-allnodes
		{"ff02::2", true},          // ip6-allrouters
		{"::ffff:192.0.2.1", true}, // ipv4-mapped ipv6 address
		{"2001:0db8:bd05:01d2:288a:1fc0:0001:10ee", true}, // IPv6 examples from wikipedia
		{"2001:0db8:0020:0003:1000:0100:0020:0003", true},
		{"2001:db8:20:3:1000:100:20:3", true},
		{"2001:0db8:0000:0000:1234:0000:0000:9abc", true},
		{"2001:db8::1234:0:0:9abc", true},
		{"2001:db8::9abc", true},
		{"::ffff:192.0.2.1", true},
		// Wrong cases
		{"0.0.0.0 ", false}, // contains white space
		{"example.com", false},
		{"192.168.0.1.com", false},
		{"192.168.0.0.1", false},
		{"192.168.0.1/24", false}, // contains subnet
		{"fe80::0123:4567:89ab:cdef%4", false},
		{"fe80::0123:4567:89ab:cdef%fxp0", false},
	} {
		got := hostpital.IsIPAddress(test.input)

		if test.expect {
			require.True(t, got, "test #%d failed. input %s should be %v", index, test.input, test.expect)
		} else {
			require.False(t, got, "test #%d failed. input %s should be %v", index, test.input, test.expect)
		}
	}
}
