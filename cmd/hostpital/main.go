//nolint:forbidigo // fmt.Println() is allowed in the main() function.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/KEINOS/go-hostpital/hostpital"
	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// Flags holds the parsed flags of the command arguments and settings to parse
// the host file.
type Flags struct {
	Args       []string
	PathIntput string
	PathOutput string
	FlagSet    *pflag.FlagSet
	Parser     *hostpital.Parser
	ShowHelp   bool
	ShowVerion bool
}

// IOFile is an interface to write and read to a file. It is a dependency injection
// interface of os.File to ease testing.
type IOFile interface {
	Write(b []byte) (n int, err error)
	Read(b []byte) (n int, err error)
}

// NameAppDefault is the name of the application for fallback. Usually the name
// is taken from the executable name.
const NameAppDefault = "hostpital"

// version is the version of the application. It is set by the build script.
var version string

// Function variables to ease testing.
//
//nolint:gochecknoglobals // These are function variables to ease testing.
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
	flags, err := ParseFlags()
	ExitOnError(err)

	if flags.ShowVerion {
		ShowVerApp() // Print version and exit
	}

	if flags.ShowHelp {
		flags.showHelp(os.Stdout, "")
		osExit(0)
	}

	flags.ShowHelpAndExitIfTrue(
		flags.PathIntput == "" && len(flags.Args) == 0,
		"Error: No file path(s) given",
	)

	listFiles := flags.Args

	if flags.PathIntput != "" {
		pattern := "hosts*"
		if len(flags.Args) > 0 {
			pattern = flags.Args[0]
		}

		listFiles, err = hostpital.FindFile(pattern, flags.PathIntput)
		ExitOnError(err)
	}

	showStdOut := false
	pathFileOut := flags.PathOutput

	if flags.PathOutput != "" {
		defer func() {
			fmt.Println("Output file:", pathFileOut)
		}()
	}

	if flags.PathOutput == "" {
		showStdOut = true

		ptrOut, err := osCreateTemp(os.TempDir(), "hostpital-*")
		ExitOnError(err)

		pathFileOut = ptrOut.Name()
		ExitOnError(ptrOut.Close())

		defer func() {
			ExitOnError(os.Remove(pathFileOut))
		}()
	}

	err = hostpital.MergeSortFiles(listFiles, flags.Parser, pathFileOut)
	ExitOnError(err)

	if showStdOut {
		byteData, err := os.ReadFile(pathFileOut)
		ExitOnError(err)

		fmt.Println(string(byteData))
	}

	// pathTmp, cleanup, err := MergeFiles(listFiles)
	// ExitOnError(err)

	// defer func() {
	// 	ExitOnError(cleanup())
	// }()

	// outFile := os.Stdout

	// if flags.PathOutput != "" {
	// 	var err error

	// 	outFile, err = os.Create(flags.PathOutput)
	// 	if err != nil {
	// 		ExitOnError(errors.Wrap(err, "failed to create the output file"))
	// 	}

	// 	defer func() {
	// 		if err := outFile.Close(); err == nil {
	// 			fmt.Println("Output file:", flags.PathOutput)
	// 		}
	// 	}()
	// }

	// ExitOnError(flags.Parser.ParseFileTo(pathTmp, outFile))
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

		if errors.Is(err, io.EOF) {
			return nil
		}

		return errors.Wrap(err, "failed to write to the file")
	}
}

// ExitOnError prints the error message to the STDERR and exits the program if
// the error is not nil.
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

func getVersion() string {
	verBin := "(unknown)"

	if version != "" {
		verBin = version
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		verBin = buildInfo.Main.Version
	}

	return verBin
}

// MergeFiles merges the given files into a temporary file and returns the path
// to the temporary file and a function to remove the temporary file.
func MergeFiles(paths []string) (string, func() error, error) {
	outFile, err := osCreateTemp(os.TempDir(), "hostpital-*")
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to create a temporary file")
	}

	defer outFile.Close()

	pathFileTmp := outFile.Name()

	for _, pathFile := range paths {
		err := appendFileTo(pathFile, outFile)
		if err != nil {
			return "", nil, errors.Wrap(err, "failed to append the file")
		}
	}

	cleanup := func() error {
		return errors.Wrap(os.Remove(pathFileTmp), "failed to remove the temporary file")
	}

	return pathFileTmp, cleanup, nil
}

// NameExec returns the name of the executable.
func NameExec() string {
	const delimiter = '.'

	pathExec, err := osExecutable()
	if err != nil && len(os.Args) > 0 {
		pathExec = os.Args[0]
	}

	nameExec := filepath.Base(pathExec)
	if nameExec == string(delimiter) {
		nameExec = NameAppDefault
	}

	// trim "hostpital.test.exe" or "hostpital.test" to "hostpital"
	foundIndex := strings.IndexByte(nameExec, delimiter)
	if foundIndex != -1 {
		return nameExec[:foundIndex]
	}

	return nameExec
}

// ParseFlags returns the Flags object with the parsed flags.
func ParseFlags() (*Flags, error) {
	flags := new(Flags)

	flags.FlagSet = pflag.NewFlagSet(NameExec(), pflag.ContinueOnError)
	flags.Parser = hostpital.NewParser()

	flags.FlagSet.StringVarP(&flags.PathIntput, "dir", "d", flags.PathOutput,
		"set directory path to search for hosts files")
	flags.FlagSet.BoolVarP(&flags.ShowHelp, "help", "h", flags.ShowHelp, "show this message")
	flags.FlagSet.StringVarP(&flags.PathOutput, "out", "o", flags.PathOutput,
		"set output file path (default: stdout)")
	flags.FlagSet.BoolVarP(&flags.Parser.IDNACompatible, "punycode", "p", flags.Parser.IDNACompatible,
		"convert unicode host names to ASCII/punycode")
	flags.FlagSet.BoolVarP(&flags.Parser.TrimComment, "remove-comment", "c", flags.Parser.TrimComment,
		"remove comment lines from the output")
	flags.FlagSet.BoolVarP(&flags.Parser.OmitEmptyLine, "remove-emptyline", "e", flags.Parser.OmitEmptyLine,
		"remove empty line(s) from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimIPAddress, "remove-ip-head", flags.Parser.TrimIPAddress,
		"remove leading IP address in the line from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimLeadingSpace, "remove-space-head", flags.Parser.TrimLeadingSpace,
		"remove leading space(s) from the output")
	flags.FlagSet.BoolVar(&flags.Parser.TrimTrailingSpace, "remove-space-tail", flags.Parser.TrimTrailingSpace,
		"remove trailing space(s) from the output")
	flags.FlagSet.BoolVarP(&flags.Parser.SortAfterParse, "sorthost", "s", flags.Parser.SortAfterParse,
		"sort the output by the host name")
	flags.FlagSet.BoolVarP(&flags.Parser.SortAsReverseDNS, "sortlabel", "l", flags.Parser.SortAsReverseDNS,
		"sort the output by the reversed labels of the DNS hosts. e.g. 'com.example.www'")
	flags.FlagSet.StringVarP(&flags.Parser.UseIPAddress, "use-ip", "i", flags.Parser.UseIPAddress,
		"set IP address to be replaced (suitable for sinkhole)")
	flags.FlagSet.BoolVarP(&flags.ShowVerion, "version", "v", flags.ShowVerion,
		"prints the version of the application")

	var err error

	const minLenArg = 1

	if len(os.Args) < minLenArg+1 {
		err = errors.New("no arguments")
	} else {
		err = flags.FlagSet.Parse(os.Args[1:])
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the flags")
	}

	flags.Args = flags.FlagSet.Args()

	return flags, nil
}

// ShowVerApp prints the version of the application. This will exit the
// application with status 0.
func ShowVerApp() {
	nameApp := NameExec()

	fmt.Printf("%s %s\n", nameApp, getVersion())

	osExit(0)
}

// -----------------------------------------------------------------------------
//  Methods
// -----------------------------------------------------------------------------

func (f *Flags) getExamples() string {
	examples := heredoc.Doc(`
		Examples:
		  $ # Merge multiple hosts files into one and print to stdout.
		  $ %%NAME_EXEC%% ./path/to/hosts ./path/to/hosts.txt ./path/to/another/file.txt

		  $ # Merge multiple hosts files into one and sort them by hostname. Then
		  $ # print to stdout.
		  $ %%NAME_EXEC%% -s ./path/to/hosts ./path/to/hosts.txt ./path/to/another/file.txt
		  $ %%NAME_EXEC%% --sorthost ./path/to/hosts ./path/to/hosts.txt ./path/to/another/file.txt

		  $ # Merge multiple hosts files into one and output to a file.
		  $ %%NAME_EXEC%% ./path/to/hosts ./path/to/hosts.txt -o ./path/to/output/merged_hosts.txt

		  $ # Search for hosts files in the directory and merge them into one and
		  $ # print to stdout ('hosts*' by default).
		  $ %%NAME_EXEC%% -d ./path/to/dir/to/search

		  $ # Search for hosts files with the given pattern ('hostfile*') in the
		  $ # directory and merge them into one and print to stdout. This will
		  $ # search for 'hostfile', 'hostfiles', 'hostfile.txt', etc.
		  $ %%NAME_EXEC%% -d ./path/to/dir/to/search hostfile
	`)

	examples = strings.ReplaceAll(examples, "%%NAME_EXEC%%", NameExec())
	examples = f.grayOutComments(examples)

	return examples
}

func (f *Flags) grayOutComments(inText string) string {
	lines := strings.Split(inText, "\n")
	grayOut := color.New(color.FgHiBlack).SprintFunc()
	outText := ""

	for _, line := range lines {
		isComment := false
		buf := ""

		for _, char := range line {
			if char == hostpital.DelimComnt {
				isComment = true
			}

			if isComment {
				buf += string(char)
			} else {
				outText += string(char)
			}
		}

		outText += grayOut(buf) + "\n"
	}

	return outText
}

func (f *Flags) showHelp(output *os.File, msg string) {
	defer color.Unset()

	f.FlagSet.SetOutput(output)

	fmt.Fprintln(output, NameExec()+" - Merge multiple hosts file(s) into one but parse and sort them.")
	fmt.Fprintln(output, "Usage:")
	fmt.Fprintf(output, "  %s [options] <file path> [<file path(s)> ...]\n", NameExec())
	fmt.Fprintf(output, "  %s [options] -d <directory path> [<search pattern>]\n", NameExec())

	fmt.Fprintln(output, "Options:")
	f.FlagSet.PrintDefaults()

	fmt.Fprintln(output, f.getExamples())

	if msg != "" {
		fmt.Fprintln(output, "\n"+msg)
	}
}

// ShowHelpAndExitIfTrue shows help and the msg to STDERR if isTrue is true.
// Then exits with status 1.
func (f *Flags) ShowHelpAndExitIfTrue(isTrue bool, msg string) {
	if !isTrue {
		return
	}

	f.showHelp(os.Stderr, msg)

	osExit(1)
}
