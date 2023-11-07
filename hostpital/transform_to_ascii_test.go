package hostpital_test

import (
	"testing"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/stretchr/testify/require"
)

func TestTransformToASCII(t *testing.T) {
	t.Parallel()

	for index, test := range []struct {
		input   string
		wantOut string
		wantErr bool
	}{
		{"example.com", "example.com", false},
		{"0.0.0.0", "0.0.0.0", false},
		{"www.xn--gpher-jua.com", "www.xn--gpher-jua.com", false},
		{"www.GÖPHER.com", "www.xn--gpher-jua.com", false},
		// Example of confusing characters.
		// U+0430 "а" --> U+0061 "a"
		// U+0440 "р" --> U+0070 "p"
		// U+04cf "ӏ" --> U+0069 "i"
		// U+0435 "е" --> U+0065 "e"
		{"www.аррӏе.com", "www.xn--80ak6aa92e.com", false},
	} {
		out, err := hostpital.TransformToASCII(test.input)

		if test.wantErr {
			require.Error(t, err, "test #%d failed. input '%s' should be an error", index+1, test.input)
			require.Contains(t, err.Error(), "failed to convert host name to ASCII",
				"test #%d failed. input '%s' should be an error", index+1, test.input)
			require.Empty(t, out, "test #%d failed. output should be empty on error. got %s", index+1, out)
		} else {
			require.NoError(t, err, "test #%d failed. input '%s' should not error", index+1, test.input)
			require.Contains(t, out, test.wantOut, "test #%d failed. it did not contain the expected output", index+1)
		}
	}
}
