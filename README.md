<!-- markdownlint-disable MD001 MD041 MD050 MD033 -->
[![go1.18+](https://img.shields.io/badge/Go-1.18--latest-blue?logo=go)](https://github.com/KEINOS/go-hostpital/blob/main/.github/workflows/unit-tests.yml "Supported versions")
[![Go Reference](https://pkg.go.dev/badge/github.com/KEINOS/go-hostpital.svg)](https://pkg.go.dev/github.com/KEINOS/go-hostpital#section-documentation "Read generated documentation of the app")

# go-hostpital

- `hostpital` is a simple library written in go to maintain and manage `hosts` files.
- `hostpital` es una sencilla librería escrita en go para mantener y gestionar archivos `hosts`.
- `hostpital` は、`hosts` ファイルを維持・管理するための go で書かれたシンプルなライブラリです。

## Usage

```go
go get "github.com/KEINOS/go-hostpital"
```

```go
import "github.com/KEINOS/go-hostpital/hostpital"

func ExampleValidator_ValidateFile() {
    // Validator with default settings
    validator := hostpital.NewValidator()

    validator.AllowComment = true    // Allow comment lines in the hostfile.
    validator.IDNACompatible = false // Want RFC 6125 2.2 compatibility. If true, IDNA2008 compatible.

    // Validate a file
    pathFile := filepath.Join("testdata", "hosts.txt")

    if validator.ValidateFile(pathFile) {
        fmt.Println("The hostfile is valid.")
    }

    // Output: The hostfile is valid.
}
```

```go
import "github.com/KEINOS/go-hostpital/hostpital"

func ExampleParser() {
    // For the default settings, see the NewValidator() example.
    parser := hostpital.NewParser()

    // Set the IP address to use for all the hosts. Suitable for DNS sinkhole.
    parser.UseIPAddress = "0.0.0.0"

    // Parse a file to clean up
    pathFile := filepath.Join("testdata", "hosts.txt")

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
```

```go
import "github.com/KEINOS/go-hostpital/hostpital"

// Various functions
func Example() {
    // Detects IPV4 and IPV6
    fmt.Println("IsIPAddress(\"0.0.0.0\") -->", hostpital.IsIPAddress("0.0.0.0"))
    fmt.Println("IsIPAddress(\"::\") -->", hostpital.IsIPAddress("::"))
    fmt.Println("IsIPAddress(\"0.0.0.0.0\") -->", hostpital.IsIPAddress("0.0.0.0.0"))

    // True if host name is ready for registration. False if it is a raw punycode or not IDNA2008 compatible.
    fmt.Println("IsIDNAComatible(\"xn--gpher-jua.com\") -->", hostpital.IsIDNAComatible("xn--gpher-jua.com"))
    fmt.Println("IsIDNAComatible(\"göpher.com\") -->", hostpital.IsIDNAComatible("göpher.com"))

    // ASCII/Punycode <---> Unicode conversion
    hostASCII, err := hostpital.TransformToASCII("göpher.com")
    fmt.Println("TransformToASCII(\"göpher.com\") -->", hostASCII, err)

    hostUnicode, err := hostpital.TransformToUnicode("xn--gpher-jua.com")
    fmt.Println("TransformToUnicode(\"xn--gpher-jua.com\") -->", hostUnicode, err)

    // Trim a comment from a line
    hostTrimmed, err := hostpital.TrimComment("127.0.0.0 localhost # this is a line comment")
    fmt.Println("TrimComments(\"127.0.0.0 localhost # this is a line comment\") --->", hostTrimmed, err)

    /* And more ... */

    // Output:
    // IsIPAddress("0.0.0.0") --> true
    // IsIPAddress("::") --> true
    // IsIPAddress("0.0.0.0.0") --> false
    // IsIDNAComatible("xn--gpher-jua.com") --> true
    // IsIDNAComatible("göpher.com") --> false
    // TransformToASCII("göpher.com") --> xn--gpher-jua.com <nil>
    // TransformToPunycode("xn--gpher-jua.com") --> göpher.com <nil>
    // TrimComments("127.0.0.0 localhost # this is a line comment") ---> 127.0.0.0 localhost  <nil>
}
```

- [View more examples](https://pkg.go.dev/github.com/KEINOS/go-hostpital/hostpital#pkg-examples) @ pkg.go.dev

## Statuses

[![UnitTests](https://github.com/KEINOS/go-hostpital/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/unit-tests.yml)
[![golangci-lint](https://github.com/KEINOS/go-hostpital/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/golangci-lint.yml)
[![CodeQL-Analysis](https://github.com/KEINOS/go-hostpital/actions/workflows/codeQL-analysis.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/codeQL-analysis.yml)
[![PlatformTests](https://github.com/KEINOS/go-hostpital/actions/workflows/platform-tests.yml/badge.svg)](https://github.com/KEINOS/go-hostpital/actions/workflows/platform-tests.yml)

[![codecov](https://codecov.io/gh/KEINOS/go-hostpital/branch/main/graph/badge.svg?token=IQKfPZPiU1)](https://codecov.io/gh/KEINOS/go-hostpital)
[![Go Report Card](https://goreportcard.com/badge/github.com/KEINOS/go-hostpital)](https://goreportcard.com/report/github.com/KEINOS/go-hostpital)

## Contributing

[![go1.18+](https://img.shields.io/badge/Go-1.18--latest-blue?logo=go)](https://github.com/KEINOS/go-hostpital/blob/main/.github/workflows/unit-tests.yml "Supported versions")
[![Go Reference](https://pkg.go.dev/badge/github.com/KEINOS/go-hostpital.svg)](https://pkg.go.dev/github.com/KEINOS/go-hostpital#section-documentation "Read generated documentation of the app")

- Branch to PR: `main`
- [CONTRIBUTING.md](https://github.com/KEINOS/go-hostpital/blob/main/.github/CONTRIBUTING.md)
- [CIs](https://github.com/KEINOS/go-hostpital/actions) on PR/Push: `unit-tests` `golangci-lint` `codeQL-analysis` `platform-tests`
- [Security policy](https://github.com/KEINOS/go-hostpital/blob/main/.github/SECURITY.md)

## License/Copyright

- [MIT License](https://github.com/KEINOS/go-hostpital/blob/main/LICENSE)
  - Copyright [KEINOS and the Hostpital contributors](https://github.com/KEINOS/go-hostpital/graphs/contributors)
- [BSD-3-Clause license](https://github.com/golang/go/blob/master/LICENSE)
  - Copyright of [`Is_compatible_rfc6125.go` by The Go Authors](https://github.com/KEINOS/go-hostpital/blob/main/hostpital/Is_compatible_rfc6125.go)
