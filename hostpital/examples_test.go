// ============================================================================
//
//	Examples for hostpital package
//
// ============================================================================
package hostpital_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Code-Hex/dd"
	"github.com/KEINOS/go-hostpital/hostpital"
)

// ----------------------------------------------------------------------------
//  FileExists()
// ----------------------------------------------------------------------------

func ExampleFileExists() {
	// FileExists returns true if the file exists.
	if ok := hostpital.FileExists("examples_test.go"); ok {
		fmt.Println("given path exists and is a file")
	}

	// FileExists returns false if the path is a directory even if the path exists.
	if ok := hostpital.FileExists(os.TempDir()); !ok {
		fmt.Println("given path is a directory")
	}

	// FileExists returns false if the path does not exist.
	if ok := hostpital.FileExists("unknown.txt"); !ok {
		fmt.Println("given path does not exist")
	}
	// Output:
	// given path exists and is a file
	// given path is a directory
	// given path does not exist
}

// ----------------------------------------------------------------------------
//  IsCommentLine()
// ----------------------------------------------------------------------------

func ExampleIsCommentLine() {
	for index, line := range []string{
		"# This is a comment line",
		"           # This is a comment line with spaces.",
		"\t\t\t# This is a comment line with tabs.", // white spaces are trimmed
		"example.com # This is an in-line comment.", // has a non-comment info
	} {
		// IsCommentLine returns true if the line is a single comment line.
		// It returns false if the line contains a non-comment info.
		fmt.Printf("#%d: %v\n", index+1, hostpital.IsCommentLine(line))
	}
	// Output:
	// #1: true
	// #2: true
	// #3: true
	// #4: false
}

// ----------------------------------------------------------------------------
//  IsCompatibleIDNA2008()
// ----------------------------------------------------------------------------

func ExampleIsCompatibleIDNA2008() {
	for index, sample := range []struct {
		input string
		want  bool
	}{
		// Golden cases
		{input: "example.com", want: true},
		{input: "0.0.0.0", want: true},
		{input: "xn--fa-hia.com", want: true},    // IDNA2008 compatible and is ASCII/punycoded
		{input: "xn--gpher-jua.com", want: true}, // same as above
		// Wrong cases
		{input: "27--m01police.55fifayellow.com"}, // Double hyphen with mal-formed punycode is not allowed
		{input: "my_host1.example.com"},           // Host contains under score
		{input: "faß.com"},                        // Must be in punycode/ASCII. Use TransformToASCII()
		{input: "www.аррӏе.com"},                  // Same as above
		{input: "*.faß.com"},                      // Wildcard is not allowed. Use IsCompatibleRFC6125Pattern()
		{input: ".example.com"},                   // Must not start with a dot
	} {
		// True if host name is ready for registration. False if it is a raw
		// punycode or not IDNA2008 compatible.
		got := hostpital.IsCompatibleIDNA2008(sample.input)

		if got != sample.want {
			log.Fatalf("failed test #%d. input: %s, want: %v, got: %v",
				index+1, sample.input, sample.want, got)
		}

		fmt.Printf("IsCompatibleIDNA2008(%#v) --> %v\n", sample.input, sample.want)
	}
	// Output:
	// IsCompatibleIDNA2008("example.com") --> true
	// IsCompatibleIDNA2008("0.0.0.0") --> true
	// IsCompatibleIDNA2008("xn--fa-hia.com") --> true
	// IsCompatibleIDNA2008("xn--gpher-jua.com") --> true
	// IsCompatibleIDNA2008("27--m01police.55fifayellow.com") --> false
	// IsCompatibleIDNA2008("my_host1.example.com") --> false
	// IsCompatibleIDNA2008("faß.com") --> false
	// IsCompatibleIDNA2008("www.аррӏе.com") --> false
	// IsCompatibleIDNA2008("*.faß.com") --> false
	// IsCompatibleIDNA2008(".example.com") --> false
}

// ----------------------------------------------------------------------------
//  IsCompatibleRFC6125
// ----------------------------------------------------------------------------

func ExampleIsCompatibleRFC6125() {
	for index, test := range []struct {
		host string
		want bool
	}{
		// Golden cases
		{host: "example.com", want: true},
		{host: "eXample123-.com", want: true},
		{host: "example.com.", want: true},
		{host: "exa_mple.com", want: true},
		{host: "127.0.0.1", want: true},
		// Wrong cases
		{host: "0.0.0.0 example.com"}, // no space allowed
		{host: "-eXample123-.com"},
		{host: ""},
		{host: "."},
		{host: "example..com"},
		{host: ".example.com"},
		{host: "*.example.com."},
		{host: "*foo.example.com"},
		{host: "foo.*.example.com"},
		{host: "foo,bar"},
		{host: "project-dev:us-central1:main"},
	} {
		got := hostpital.IsCompatibleRFC6125(test.host)

		if got != test.want {
			log.Fatalf("failed test #%d: IsCompatibleRFC6125(%#v) --> %v, want: %v",
				index+1, test.host, got, test.want)
		}
	}

	fmt.Println("OK")
	// Output: OK
}

func ExampleIsCompatibleRFC6125Pattern() {
	for index, test := range []struct {
		host string
		want bool
	}{
		// Golden cases
		{host: "example.com", want: true},
		{host: "eXample123-.com", want: true},
		{host: "exa_mple.com", want: true},
		{host: "127.0.0.1", want: true},
		{host: "*.example.com", want: true}, // wildcard is allowed
		// Wrong cases
		{host: "0.0.0.0 example.com"}, // no space allowed
		{host: "example.com."},        // dot at the end is not allowed
		{host: "-eXample123-.com"},
		{host: ""},
		{host: "."},
		{host: "example..com"},
		{host: ".example.com"},
		{host: "*.example.com."},
		{host: "*foo.example.com"},
		{host: "foo.*.example.com"},
		{host: "foo,bar"},
		{host: "project-dev:us-central1:main"},
	} {
		got := hostpital.IsCompatibleRFC6125Pattern(test.host)

		if got != test.want {
			log.Fatalf("failed test #%d: IsCompatibleRFC6125Pattern(%#v) --> %v, want: %v",
				index+1, test.host, got, test.want)
		}
	}

	fmt.Println("OK")
	// Output: OK
}

// ----------------------------------------------------------------------------
//  IsExistingFile()
// ----------------------------------------------------------------------------

func ExampleIsExistingFile() {
	for index, sample := range []struct {
		nameFile string
		want     bool
	}{
		{nameFile: "hosts.txt", want: true},
		{nameFile: "hosts_not_exists.txt", want: false},
		{nameFile: "", want: false}, // is directory
	} {
		pathFile := filepath.Join("testdata", sample.nameFile)

		got := hostpital.IsExistingFile(pathFile)

		if got != sample.want {
			log.Fatalf("failed test #%d. input: %s, want: %v, got: %v",
				index+1, pathFile, sample.want, got)
		}
	}

	fmt.Println("OK")
	// Output: OK
}

// ----------------------------------------------------------------------------
//  IsIPAddress()
// ----------------------------------------------------------------------------

func ExampleIsIPAddress() {
	for index, test := range []struct {
		input string
		want  bool
	}{
		// Golden cases
		{input: "0.0.0.0", want: true},
		{input: "::", want: true},
		// Wrong cases
		{input: "0.0.0.0.0"},
	} {
		// Detects IP address (IPV4 and IPV6)
		got := hostpital.IsIPAddress(test.input)
		if got != test.want {
			log.Fatalf("test #%v failed. IsIPAddress(%#v) --> %v (want: %v)",
				index+1, test.input, got, test.want)
		}

		fmt.Printf("IsIPAddress(%#v) --> %v\n", test.input, got)
	}
	// Output:
	// IsIPAddress("0.0.0.0") --> true
	// IsIPAddress("::") --> true
	// IsIPAddress("0.0.0.0.0") --> false
}

// ----------------------------------------------------------------------------
//  Type: Parser
// ----------------------------------------------------------------------------

// This example parses the hosts file as a DNS sinkhole. Which all the hosts will
// point to 0.0.0.0 and will not be able to connect to the Internet.
func ExampleParser() {
	pathFile := filepath.Join("testdata", "default.txt")

	// For the default settings, see the NewValidator() example.
	parser := hostpital.NewParser()

	// Set the IP address to use for all the hosts.
	parser.UseIPAddress = "0.0.0.0"

	parsed, err := parser.ParseFile(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parsed)
	// Output:
	// 0.0.0.0 dummy1.example.com
	// 0.0.0.0 dummy2.example.com
	// 0.0.0.0 dummy3.example.com
	// 0.0.0.0 dummy4.example.com
	// 0.0.0.0 dummy5.example.com dummy6.example.com
}

func ExampleParser_ParseFileTo() {
	pathFile := filepath.Join("testdata", "default.txt")

	parser := hostpital.NewParser()

	// Set property to sort after parsing
	parser.SortAfterParse = true

	// Create a buffer to store the parsed hosts file
	var buf bytes.Buffer

	// Execute ParseFileTo
	if err := parser.ParseFileTo(pathFile, &buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	// Output:
	// dummy1.example.com
	// dummy2.example.com
	// dummy3.example.com
	// dummy4.example.com
	// dummy5.example.com dummy6.example.com
}

func ExampleParser_ParseString() {
	hosts := `# this is a comment
badboy5.example.com      badboy6.example.com

# this is another comment
123.123.123.120 badboy4.example.com
123.123.123.121 badboy3.example.com
123.123.123.122 badboy2.example.com
123.123.123.123 badboy1.example.com
`

	parser := hostpital.NewParser()

	// Set property for user custom settings
	parser.SortAfterParse = true
	parser.UseIPAddress = "0.0.0.0"

	parsed := parser.ParseString(hosts)

	fmt.Println(parsed)
	// Output:
	// 0.0.0.0 badboy1.example.com
	// 0.0.0.0 badboy2.example.com
	// 0.0.0.0 badboy3.example.com
	// 0.0.0.0 badboy4.example.com
	// 0.0.0.0 badboy5.example.com badboy6.example.com
}

// ----------------------------------------------------------------------------
//  PickRandom()
// ----------------------------------------------------------------------------

func ExamplePickRandom() {
	items := []string{
		"one.example.com",
		"two.example.com",
		"three.example.com",
		"four.example.com",
		"five.example.com",
	}

	changed := false
	picked1st := hostpital.PickRandom(items)

	// Pick 10 times since the number of items is too few.
	for i := 0; i < 10; i++ {
		// Sleep for a random time from 0 to 999 milliseconds to avoid the same
		// seed for the random number generator. The CI server is too fast.
		hostpital.SleepRandom(1)

		pickedCurr := hostpital.PickRandom(items)

		if pickedCurr != picked1st {
			changed = true
		}
	}

	fmt.Println("changed:", changed)
	// Output:
	// changed: true
}

// ----------------------------------------------------------------------------
//  ReverseDNS()
// ----------------------------------------------------------------------------

func ExampleReverseDNS() {
	// ReverseDNS reverses the order of the labels in a domain name.
	// Useful for grouping hosts by domain name.
	fmt.Println(hostpital.ReverseDNS("www.example.com"))
	fmt.Println(hostpital.ReverseDNS("com.example.www"))
	// Output:
	// com.example.www
	// www.example.com
}

// ----------------------------------------------------------------------------
//  TransformToASCII()
// ----------------------------------------------------------------------------

func ExampleTransformToASCII() {
	// Unicode --> ASCII/Punycode conversion.
	// For the opposite, see TransformToUnicode().
	hostASCII, err := hostpital.TransformToASCII("göpher.com")
	fmt.Println("TransformToASCII(\"göpher.com\") -->", hostASCII, err)
	// Output:
	// TransformToASCII("göpher.com") --> xn--gpher-jua.com <nil>
}

// ----------------------------------------------------------------------------
//  TransformToUnicode()
// ----------------------------------------------------------------------------

func ExampleTransformToUnicode() {
	// ASCII/Punycode --> Unicode conversion. It will error if the input is not
	// convertable to Unicode. For the opposite, see TransformToASCII().
	hostPunycode, err := hostpital.TransformToUnicode("xn--gpher-jua.com")
	fmt.Println("TransformToUnicode(\"xn--gpher-jua.com\") -->", hostPunycode, err)
	// Output:
	// TransformToUnicode("xn--gpher-jua.com") --> göpher.com <nil>
}

// ----------------------------------------------------------------------------
//  TrimComment()
// ----------------------------------------------------------------------------

func ExampleTrimComment() {
	line := "127.0.0.0 localhost # this is a line comment"

	// Trim a comment from a line
	hostTrimmed, err := hostpital.TrimComment(line)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hostTrimmed)
	// Output: 127.0.0.0 localhost
}

// ----------------------------------------------------------------------------
//  TrimDNSByLevel()
// ----------------------------------------------------------------------------

func ExampleTrimDNSByLevel() {
	for index, test := range []struct {
		host  string
		want  string
		level int
	}{
		{host: "www.example.com", level: 0, want: "com"},
		{host: "www.example.com", level: 1, want: "example.com"},
		{host: "www.example.com", level: 2, want: "www.example.com"},
		{host: "www.example.com", level: 5, want: "www.example.com"},
	} {
		got := hostpital.TrimDNSByLevel(test.host, test.level)
		if got != test.want {
			log.Fatalf("failed test #%d: host: %s, level: %d, want: %s, got: %s",
				index, test.host, test.level, test.want, got)
		}

		fmt.Printf("Level %d: %s --> %s\n", test.level, test.host, got)
	}
	// Output:
	// Level 0: www.example.com --> com
	// Level 1: www.example.com --> example.com
	// Level 2: www.example.com --> www.example.com
	// Level 5: www.example.com --> www.example.com
}

// ----------------------------------------------------------------------------
//  TrimIPAdd()
// ----------------------------------------------------------------------------

func ExampleTrimIPAdd() {
	for index, line := range []struct {
		input string
		want  string
	}{
		{"", ""},
		{" ", ""},
		{"123.123.123.123", ""},
		{"example.com", "example.com"},
		{"\texample.com", "example.com"},
		{"      example.com", "example.com"},
		{"123.123.123.123 example.com ", "example.com"},
		{"123.123.123.123       example.com ", "example.com"},
		{"123.123.123.123\texample.com ", "example.com"},
		{
			"123.123.123.123    0.0.0.0    sub1.example.com    sub2.example.com",
			"sub1.example.com sub2.example.com",
		},
	} {
		expect := line.want
		actual := hostpital.TrimIPAdd(line.input)

		fmt.Printf("#%d: %q --> %q ... ", index+1, line.input, actual)

		if actual != expect {
			fmt.Printf("FAIL (want: %q, got: %q)\n", expect, actual)
		} else {
			fmt.Println("PASS")
		}
	}
	//nolint:lll // The last line is too long, but leave it as is for the sake of readability.
	// Output:
	// #1: "" --> "" ... PASS
	// #2: " " --> "" ... PASS
	// #3: "123.123.123.123" --> "" ... PASS
	// #4: "example.com" --> "example.com" ... PASS
	// #5: "\texample.com" --> "example.com" ... PASS
	// #6: "      example.com" --> "example.com" ... PASS
	// #7: "123.123.123.123 example.com " --> "example.com" ... PASS
	// #8: "123.123.123.123       example.com " --> "example.com" ... PASS
	// #9: "123.123.123.123\texample.com " --> "example.com" ... PASS
	// #10: "123.123.123.123    0.0.0.0    sub1.example.com    sub2.example.com" --> "sub1.example.com sub2.example.com" ... PASS
}

// ----------------------------------------------------------------------------
//  TrimWordGaps()
// ----------------------------------------------------------------------------

func ExampleTrimWordGaps() {
	for index, input := range []string{
		"0.0.0.0                  example.com",
		"0.0.0.0\texample.com",
		" \t\n\t 127.0.0.1\texample.com\n\n\t#    inline     comment",
		"127.0.0.1\t              example.com      # inline\tcomment",
	} {
		result := hostpital.TrimWordGaps(input)

		fmt.Printf("#%d: %q\n", index+1, result)
	}
	// Output:
	// #1: "0.0.0.0 example.com"
	// #2: "0.0.0.0 example.com"
	// #3: "127.0.0.1 example.com # inline comment"
	// #4: "127.0.0.1 example.com # inline comment"
}

// ----------------------------------------------------------------------------
//  Type: Validator
// ----------------------------------------------------------------------------

func ExampleValidator() {
	validator := hostpital.NewValidator()

	// Print default settings
	fmt.Println(dd.Dump(validator))

	// Output:
	// &hostpital.Validator{
	//   mutx: sync.Mutex{
	//     state: 0,
	//     sema: 0,
	//   },
	//   AllowComment: false,
	//   AllowEmptyLine: true,
	//   AllowHyphen: false,
	//   AllowHyphenDouble: false,
	//   AllowIndent: false,
	//   AllowIPAddressOnly: false,
	//   AllowTrailingSpace: false,
	//   AllowUnderscore: false,
	//   IDNACompatible: true,
	//   isInitialized: true,
	// }
}

func ExampleValidator_ValidateFile() {
	// Validator with default settings
	validator := hostpital.NewValidator()

	// Want RFC 6125 2.2 compatibility. If true, IDNA2008 compatible.
	validator.IDNACompatible = false
	// Allow comment lines in the hostfile.
	validator.AllowComment = true
	// Allow labels to begin with hyphen.
	// This setting is useful to manage hosts file for DNS sinkhole. "traditional
	// domain name" does not allow hyphen in the first character position of
	// their labels. Such as "m.-www99a.abc.example.com" for example. Usually
	// it is blocked by the client side's router, browser, etc. However, some
	// malicious domains owns their name server configured to resolve it to find
	// out web clients who doesn't care about it.
	validator.AllowHyphen = true

	// Validate a file
	pathFile := filepath.Join("testdata", "hosts.txt")

	if validator.ValidateFile(pathFile) {
		fmt.Println("The hostfile is valid.")
	}

	// Output: The hostfile is valid.
}

func ExampleValidator_ValidateLine() {
	// Validator with default settings
	//   AllowComment: false
	//   AllowEmptyLine: true
	//   AllowHyphen: false
	//   AllowHyphenDouble: false
	//   AllowIndent: false
	//   AllowIPAddressOnly: false
	//   AllowTrailingSpace: false
	//   AllowUnderscore: false
	//   IDNACompatible: true
	validator := hostpital.NewValidator()

	// User custom settings
	validator.AllowTrailingSpace = true

	for index, line := range []string{
		// Valid cases according to the settings.
		"example.com",
		"example.com           ",
		"",
		// Invalid cases according to the settings.
		"           example.com",
		"123.123.123.123",
		"# This is a comment line",
		"example.com # This is an in-line comment",
		"         ", // false. Empty line is allowed but indent is not allowed.
	} {
		err := validator.ValidateLine(line)
		fmt.Printf("#%d: %v\n", index+1, err)
	}
	// Output:
	// #1: <nil>
	// #2: <nil>
	// #3: <nil>
	// #4: failed to trim line: indent is not allowed
	// #5: IP address only line is not allowed
	// #6: failed to validate chunk/part of line: "#" is not IDNA2008 compatible: idna: disallowed rune U+0023
	// #7: failed to validate chunk/part of line: "#" is not IDNA2008 compatible: idna: disallowed rune U+0023
	// #8: failed to trim line: indent is not allowed
}
