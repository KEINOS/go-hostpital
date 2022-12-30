# Example Application

This is a simple example application that uses the `hostpital` package.

## Install

```go
// Install the latest release version
go install "github.com/KEINOS/go-hostpital/cmd/hostpital@latest"
```

```go
// Install the latest dev version
go install "github.com/KEINOS/go-hostpital/cmd/hostpital@main"
```

## Usage

```shellsession
$ hostpital -h
hostpital - Merge multiple hosts file(s) into one but parse and sort them.
Usage: hostpital [options] <file path(s)>
Options:
  -e, --emptyline           remove empty line(s) from the output (default true)
  -h, --help                show this message
  -o, --out string          set output file path (default: stdout)
  -p, --punycode            convert unicode host names to ASCII/punycode (default true)
      --remove-comment      remove comment lines from the output (default true)
      --remove-ip-head      remove leading IP address in the line from the output (default true)
      --remove-space-head   remove leading space(s) from the output (default true)
      --remove-space-tail   remove trailing space(s) from the output (default true)
  -s, --sorthost            sort the output by the host name
  -l, --sortlabel           sort the output by the reversed labels of the DNS hosts. e.g. 'com.example.www'
  -i, --use-ip string       set IP address to be replaced (suitable for sinkhole)
  -v, --version             prints the version of the application
```

```shellsession
$ hostpital ./testdata/host1.txt ./testdata/host2.txt
badboy1.example.com
badboy2.example.com badboy3.example.com
badboy2.example.jp badboy3.example.jp
badboy1.example.jp

$ hostpital --sorthost ./testdata/host1.txt ./testdata/host2.txt
badboy1.example.com
badboy1.example.jp
badboy2.example.com badboy3.example.com
badboy2.example.jp badboy3.example.jp

$ hostpital ./testdata/host1.txt ./testdata/host2.txt --sorthost --use-ip "0.0.0.0"
0.0.0.0 badboy1.example.com
0.0.0.0 badboy1.example.jp
0.0.0.0 badboy2.example.com badboy3.example.com
0.0.0.0 badboy2.example.jp badboy3.example.jp
```
