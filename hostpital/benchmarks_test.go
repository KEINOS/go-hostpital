// ============================================================================
//
//	Benchmarks for hostpital package
//
// ============================================================================
package hostpital_test

import (
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-hostpital/hostpital"
)

func BenchmarkParser(b *testing.B) {
	const wantIP = "0.0.0.0"

	for range b.N {
		pathFile := filepath.Join("testdata", "hosts.txt")

		// For the default settings, see the NewValidator() example.
		parser := hostpital.NewParser()

		// Set the IP address to use for all the hosts.
		parser.UseIPAddress = wantIP

		_, err := parser.ParseFile(pathFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReverseDNS(b *testing.B) {
	for range b.N {
		_ = hostpital.ReverseDNS("www.example.com")
	}
}
