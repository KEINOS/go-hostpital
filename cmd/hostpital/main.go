package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// Flags holds the parsed flags of the command arguments and settings to parse
// the host file.
type Flags struct {
	Args       []string
	PathOutput string
	FlagSet    *pflag.FlagSet
	Parser     *hostpital.Parser
	IsHelp     bool
}

type IOFile interface {
	Write(b []byte) (n int, err error)
	Read(b []byte) (n int, err error)
}

// NameAppDefault is the name of the application for fallback. Usually the name
// is taken from the executable name.
const NameAppDefault = "hostpital"

// Function variables to ease testing.
var (
	// osExit is a copy of os.Exit to ease testing.
	osExit = os.Exit
	// osExecutable is a copy of os.Executable to ease testing.
	osExecutable = os.Executable
	// osCreateTemp is a copy of os.CreateTemp to ease testing.
	osCreateTemp = os.CreateTemp
)

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

func main() {
	flags := ParseFlags()

	flags.ShowHelpAndExitIfTrue(flags.IsHelp, "")
	flags.ShowHelpAndExitIfTrue(len(flags.Args) == 0, "Error: No file path(s) given")

	pathTmp, cleanup, err := MergeFiles(flags.Args)
	ExitOnError(err)

	defer cleanup()

	// fmt.Println(dd.Dump(flags))
	// fmt.Println(dd.Dump(pflag.Args()))
	// fmt.Println("Temp path:", pathTmp)

	var outFile *os.File = os.Stdout

	if flags.PathOutput != "" {
		var err error

		outFile, err = os.Create(flags.PathOutput)
		if err != nil {
			ExitOnError(errors.Wrap(err, "failed to create the output file"))
		}

		defer func() {
			if err := outFile.Close(); err == nil {
				fmt.Println("Output file:", flags.PathOutput)
			}
		}()
	}

	ExitOnError(flags.Parser.ParseFileTo(pathTmp, outFile))
}

// -----------------------------------------------------------------------------
//  Functions (methods follows below)
// -----------------------------------------------------------------------------

func appendFileTo(inputFile string, outFile IOFile) error {
	buf := make([]byte, bufio.MaxScanTokenSize)

	inFile, err := os.Open(inputFile)
	if err != nil {
		return errors.Wrap(err, "failed to open the file")
	}

	defer inFile.Close()

	// We do not early return on error here to ease capture buffer write errors.
	for {
		n, err := inFile.Read(buf)
		if err == nil {
			_, err = outFile.Write(buf[:n])
			if err == nil {
				continue
			}
		}

		if err == io.EOF {
			return nil
		}

		return errors.Wrap(err, "failed to write to the file")
	}
}

// ExitOnError prints the error message to the STDERR and exits the program.
//
// To mock the behavior of os.Exit() for testing, override the osExit function
// variable.
//
// Example:
//
//	oldOsExit := osExit
//	defer func() { osExit = oldOsExit }()
//	osExit = func(code int) { fmt.Println("Exit code:", code) }
func ExitOnError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		osExit(1)
	}
}

func MergeFiles(paths []string) (string, func() error, error) {
	outFile, err := osCreateTemp(os.TempDir(), "hostpital-*")
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to create a temporary file")
	}

	defer outFile.Close()

	pathFileTmp := outFile.Name()

	for _, pathFile := range paths {
		if err := appendFileTo(pathFile, outFile); err != nil {
			return "", nil, errors.Wrap(err, "failed to append the file")
		}
	}

	return pathFileTmp, func() error { return os.Remove(pathFileTmp) }, nil
}

// NameExec returns the name of the executable.
func NameExec() string {
	pathExec, err := osExecutable()
	if err != nil && len(os.Args) > 0 {
		pathExec = os.Args[0]
	}

	nameExec := filepath.Base(pathExec)
	if nameExec == "." {
		nameExec = NameAppDefault
	}

	return nameExec
}

// ParseFlags returns the Flags object with the parsed flags.
func ParseFlags() *Flags {
	flags := new(Flags)

	flags.FlagSet = pflag.NewFlagSet(NameExec(), pflag.ContinueOnError)
	flags.Parser = hostpital.NewParser()

	flags.FlagSet.BoolVarP(&flags.Parser.OmitEmptyLine, "emptyline", "e", flags.Parser.OmitEmptyLine,
		"remove empty line(s) from the output")
	flags.FlagSet.BoolVarP(&flags.IsHelp, "help", "h", flags.IsHelp, "show this message")
	flags.FlagSet.StringVarP(&flags.Parser.UseIPAddress, "ip", "i", flags.Parser.UseIPAddress,
		"set IP address to be replaced (suitable for sinkhole)")
	flags.FlagSet.StringVarP(&flags.PathOutput, "out", "o", flags.PathOutput,
		"set output file path (default: stdout)")
	flags.FlagSet.BoolVarP(&flags.Parser.IDNACompatible, "punycode", "p", flags.Parser.IDNACompatible,
		"convert unicode host names to ASCII/punycode")
	flags.FlagSet.BoolVarP(&flags.Parser.SortAfterParse, "sorthost", "s", flags.Parser.SortAfterParse,
		"sort the output by the host name")
	flags.FlagSet.BoolVarP(&flags.Parser.SortAsReverseDNS, "sortlabel", "l", flags.Parser.SortAsReverseDNS,
		"sort the output by the reversed labels of the DNS hosts. e.g. 'com.example.www'")
	flags.FlagSet.BoolVar(&flags.Parser.TrimComment, "no-comment", flags.Parser.TrimComment,
		"remove comment lines from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimIPAddress, "no-ip", flags.Parser.TrimIPAddress,
		"remove leading IP address in the line from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimLeadingSpace, "no-space-head", flags.Parser.TrimLeadingSpace,
		"remove leading space(s) from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimTrailingSpace, "no-space-tail", flags.Parser.TrimTrailingSpace,
		"remove trailing space(s) from the output")

	flags.FlagSet.Parse(os.Args[1:])

	flags.Args = flags.FlagSet.Args()

	return flags
}

// -----------------------------------------------------------------------------
//  Methods
// -----------------------------------------------------------------------------

// ShowHelpAndExit shows help and the msg if isTrue is true then exits with
// status 1.
func (f *Flags) ShowHelpAndExitIfTrue(isTrue bool, msg string) {
	if !isTrue {
		return
	}

	fmt.Fprintln(os.Stderr, NameExec()+" - Merge multiple hosts file(s) into one")

	fmt.Fprintln(os.Stderr, "Usage: "+NameExec()+" [options] <file path(s)>")

	fmt.Fprintln(os.Stderr, "Options:")
	f.FlagSet.PrintDefaults()

	if msg != "" {
		fmt.Fprintln(os.Stderr, "\n"+msg)
	}

	osExit(1)
}
