package hostpital

// ============================================================================
//  Function isValidHostname is a copy of the internal function `validHostname`
//  from the crypto/x509 package with the following license.
// ============================================================================
//  Copyright (c) 2009 The Go Authors. All rights reserved.
//
//  Redistribution and use in source and binary forms, with or without
//  modification, are permitted provided that the following conditions are
//  met:
//
//     * Redistributions of source code must retain the above copyright
//  notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
//  copyright notice, this list of conditions and the following disclaimer
//  in the documentation and/or other materials provided with the
//  distribution.
//     * Neither the name of Google Inc. nor the names of its
//  contributors may be used to endorse or promote products derived from
//  this software without specific prior written permission.
//
//  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//  "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//  LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//  A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//  OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//  SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//  LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//  DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//  THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ----------------------------------------------------------------------------
//  Ref: https://github.com/golang/go/blob/master/src/crypto/x509/verify.go#L965-L1010
// ----------------------------------------------------------------------------

import "strings"

// IsCompatibleRFC6125 returns true if the given host name can be matched or matched
// against according to RFC 6125 2.2, with some leniency to accommodate legacy
// values.
//
// Note: This function is a copy of the internal function from the crypto/x509
// package under BSD-3-Clause license.
// Please see the source code in "is_valid_hostname.go" for more information.
func IsCompatibleRFC6125(host string) bool {
	return isValidHostname(host, false)
}

// IsCompatibleRFC6125Pattern is similar to IsValidHostname but it also allows
// wildcard patterns.
//
// Note: This function is a copy of the internal function from the crypto/x509
// package under BSD-3-Clause license.
// Please see the source code in "is_valid_hostname.go" for more information.
func IsCompatibleRFC6125Pattern(host string) bool {
	return isValidHostname(host, true)
}

func isValidHostname(host string, isPattern bool) bool {
	if !isPattern {
		host = strings.TrimSuffix(host, ".")
	}

	if len(host) == 0 {
		return false
	}

	for index, label := range strings.Split(host, ".") {
		if label == "" {
			return false
		}

		if isPattern && index == 0 && label == "*" {
			// Only allow full left-most wildcards, as those are the only ones
			// we match, and matching literal '*' characters is probably never
			// the expected behavior.
			continue
		}

		if !isValidLabel(label) {
			return false
		}
	}

	return true
}

func isValidLabel(label string) bool {
	for pos, char := range label {
		if 'a' <= char && char <= 'z' {
			continue
		}

		if '0' <= char && char <= '9' {
			continue
		}

		if 'A' <= char && char <= 'Z' {
			continue
		}

		if char == '-' && pos != 0 {
			continue
		}

		if char == '_' {
			// Not a valid character in hostnames, but commonly
			// found in deployments outside the WebPKI.
			continue
		}

		return false
	}

	return true
}
