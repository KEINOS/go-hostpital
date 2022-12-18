package hostpital

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Parser.CountLines()
// ----------------------------------------------------------------------------

func TestParser_CountLines(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	result, err := parser.CountLines("")

	require.Error(t, err)
	require.Empty(t, result, "it should be empty on error")
	assert.Contains(t, err.Error(), "failed to open the file", "it should contain the error reason")

	if runtime.GOOS == "windows" {
		assert.Contains(t, err.Error(), "The system cannot find the file specified", "it should contain the error reason")
	} else {
		assert.Contains(t, err.Error(), "no such file or directory", "it should contain the error reason")
	}
}

// ----------------------------------------------------------------------------
//  Parser.onlyIDNACompatible()
// ----------------------------------------------------------------------------

func TestParser_onlyIDNACompatible(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	{
		parser.IDNACompatible = false

		expect := "foo bar baz"
		actual := parser.onlyIDNACompatible(expect)

		require.Equal(t, expect, actual, "it should return as is if IDNACompatible is false")
	}
	{
		parser.IDNACompatible = true

		expect := "# this is a sample comment line"
		actual := parser.onlyIDNACompatible(expect)

		require.Equal(t, expect, actual, "it should return as is if the line is a comment")
	}
}

// ----------------------------------------------------------------------------
//  Parser.ParseFile()
// ----------------------------------------------------------------------------

func TestParser_ParseFile(t *testing.T) {
	t.Parallel()

	pathDirTemp := t.TempDir()
	parser := NewParser()

	result, err := parser.ParseFile(pathDirTemp)

	require.Error(t, err)
	require.Empty(t, result, "it should be empty on error")
	assert.Contains(t, err.Error(), "failed to read from reader",
		"it should contain the error reason")
}

// ----------------------------------------------------------------------------
//  Parser.ParseFileTo()
// ----------------------------------------------------------------------------

func TestParser_ParseFileTo_nil_arg(t *testing.T) {
	t.Parallel()

	pathDirTemp := t.TempDir()
	parser := NewParser()

	err := parser.ParseFileTo(pathDirTemp, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "the given io.Writer is nil",
		"it should contain the error reason")
}

func TestParser_ParseFileTo_outfile_is_dir(t *testing.T) {
	t.Parallel()

	pathDirTemp := t.TempDir()
	pathFileTemp := filepath.Join(pathDirTemp, "hosts_to_sort")

	err := os.WriteFile(
		pathFileTemp,
		[]byte(heredoc.Doc(`
			host3.example.com
			host2.example.com
			host1.example.com`,
		)),
		os.ModePerm,
	)
	require.NoError(t, err, "failed to create temp file for testing")

	parser := NewParser()

	// Set property to sort after parsing
	parser.SortAfterParse = true

	pt, err := os.OpenFile(pathDirTemp, os.O_RDONLY, os.ModePerm)
	require.NoError(t, err, "failed to open temp dir for testing")

	defer pt.Close()

	err = parser.ParseFileTo(pathFileTemp, pt)

	require.Error(t, err, "it should return an error if outfile is nil")
	assert.Contains(t, err.Error(), "failed to write to io.Writer",
		"it should contain the error reason")
}

// ----------------------------------------------------------------------------
//  Parser.prependIPAddress()
// ----------------------------------------------------------------------------

func TestParser_prependIPAddress(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	// Regular settings
	{
		parser.UseIPAddress = "0.0.0.0" // set an IP address to prepend

		input := "foo.example.com"
		expect := "0.0.0.0 foo.example.com"
		actual := parser.prependIPAddress(input)

		require.Equal(t, expect, actual, "it should return the input with prepended IP address set")
	}

	// IP address is empty
	{
		parser.UseIPAddress = "" // set to empty

		input := "bar.example.com"
		expect := "bar.example.com"
		actual := parser.prependIPAddress(input)

		require.Equal(t, expect, actual, "it should return as is if UseIPAddress is empty")
	}

	// Input is an IP address
	{
		parser.UseIPAddress = "0.0.0.0" // set an IP address to prepend

		input := "127.0.0.1"
		expect := "127.0.0.1"
		actual := parser.prependIPAddress(input)

		require.Equal(t, expect, actual, "it should return as is if the input is an IP address")
	}
}

// ----------------------------------------------------------------------------
//  Parser.scanFile()
// ----------------------------------------------------------------------------

type dummyReader struct {
	dummyFn func(p []byte) (int, error)
}

func (d *dummyReader) Read(p []byte) (int, error) {
	return d.dummyFn(p)
}

func TestParser_scanFile(t *testing.T) {
	t.Parallel()

	parser := NewParser()
	dummy := new(dummyReader)

	dummy.dummyFn = func(p []byte) (int, error) {
		return 0, errors.New("forced error")
	}

	err := parser.scanFile(dummy, []string{"foo.example.com"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read/scan the file",
		"it should error if the reader returns an error")
	assert.Contains(t, err.Error(), "forced error",
		"it should contain the error reason")
}

// ----------------------------------------------------------------------------
//  Parser.sortAsReverseDNS()
// ----------------------------------------------------------------------------

func TestParser_sortAsReverseDNS(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	lines := []string{
		"one.example.co.jp",
		"three.example.com",
		"two.example.jp",
	}

	// Regulat settings (no sorting)
	{
		parser.SortAsReverseDNS = false

		expect := []string{
			"one.example.co.jp",
			"three.example.com",
			"two.example.jp",
		}

		actual := parser.sortAsReverseDNS(lines)

		require.Equal(t, expect, actual,
			"it should return as is if the SortAsReverseDNS field is false")
	}

	// Sort by reversed DNS. For regular sort use sortSlices() instead.
	{
		parser.SortAsReverseDNS = true

		expect := []string{
			"three.example.com",
			"one.example.co.jp",
			"two.example.jp",
		}

		actual := parser.sortAsReverseDNS(lines)

		require.Equal(t, expect, actual,
			"it should sort the lines in reverse DNS order if the SortAsReverseDNS field is true")
	}
}

// ----------------------------------------------------------------------------
//  Parser.sortSlices()
// ----------------------------------------------------------------------------

func TestParser_sortSlices(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	lines := []string{
		"3.example.co.jp",
		"2.example.jp",
		"1.example.com",
	}

	// Regulat settings (no sorting)
	{
		parser.SortAsReverseDNS = false

		expect := []string{
			"1.example.com",
			"2.example.jp",
			"3.example.co.jp",
		}

		actual := parser.sortSlices(lines)

		require.Equal(t, expect, actual,
			"it should sort the lines in regular way if the SortAsReverseDNS field is false")
	}

	// Sort by reversed DNS. For regular sort set SortAsReverseDNS to true instead.
	{
		parser.SortAsReverseDNS = true

		expect := []string{
			"1.example.com",
			"3.example.co.jp",
			"2.example.jp",
		}

		actual := parser.sortSlices(lines)

		require.Equal(t, expect, actual,
			"it should sort the lines in reverse DNS order if the SortAsReverseDNS field is true")
	}
}

// ----------------------------------------------------------------------------
//  Parser.trimComment()
// ----------------------------------------------------------------------------

func TestParser_trimComment(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	parser.TrimComment = false

	expect := "# this is a comment line"
	actual := parser.trimComment(expect)
	require.Equal(t, expect, actual, "it should return as is if TrimComment is false")
}

// ----------------------------------------------------------------------------
//  Parser.trimIPAddress()
// ----------------------------------------------------------------------------

func TestParser_trimIPAddress(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	expect := "# this is a comment line"
	actual := parser.trimIPAddress(expect)

	require.Equal(t, expect, actual, "it should return as is if the line is a comment")
}

// ----------------------------------------------------------------------------
//  Parser.trimSpace()
// ----------------------------------------------------------------------------

func TestParser_trimSpace(t *testing.T) {
	t.Parallel()

	parser := NewParser()

	const line = "    foo bar baz    "

	{
		parser.TrimLeadingSpace = false
		parser.TrimTrailingSpace = false

		expect := line
		actual := parser.trimSpace(line)

		require.Equal(t, expect, actual,
			"it should return as is if both TrimLeadingSpace and TrimTrailingSpace are false")
	}
	{
		parser.TrimLeadingSpace = true
		parser.TrimTrailingSpace = false

		expect := "foo bar baz    "
		actual := parser.trimSpace(line)

		require.Equal(t, expect, actual,
			"it should trim leading space if TrimLeadingSpace is true")
	}
	{
		parser.TrimLeadingSpace = false
		parser.TrimTrailingSpace = true

		expect := "    foo bar baz"
		actual := parser.trimSpace(line)

		require.Equal(t, expect, actual,
			"it should trim trailing space if TrimTrailingSpace is true")
	}
}
