package hostpital

import (
	"slices"
	"strings"
)

// ReverseDNS converts the DNS labels in the reverse order.
//
// For example, "www.google.com" will be converted to "com.google.www".
func ReverseDNS(hostName string) string {
	chunks := strings.Split(hostName, string(DelimDNS))

	slices.SortFunc(chunks, func(_ string, _ string) int {
		return -1 // always a < b
	})

	return strings.Join(chunks, string(DelimDNS))
}

// Old version of ReverseDNS. This comment will be replaced to the current function
// if a faster function is found.
//
// // Benchmark results:
// //
// // goos: darwin
// // goarch: amd64
// // pkg: github.com/KEINOS/go-hostpital/hostpital
// // cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
// //
// // BenchmarkReverseDNS/ReverseDNS-4         	 5350540	       220.2 ns/op	      64 B/op	       2 allocs/op
// // BenchmarkReverseDNS/ReverseDNS2-4        	 3545688	       335.1 ns/op	     120 B/op	       4 allocs/op
//
// func ReverseDNS2(hostName string) string {
// 	chunks := strings.Split(hostName, string(DelimDNS))
//
// 	sort.Slice(chunks, func(i, j int) bool {
// 		return true
// 	})
//
// 	return strings.Join(chunks, string(DelimDNS))
// }
