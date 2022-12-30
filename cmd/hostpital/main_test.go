//nolint:paralleltest // do not parallelize due to temporary changing global variables
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

// ============================================================================
//  Golden Cases
// ============================================================================

func Test_main_golden_show_version(t *testing.T) {
	// Backup and defer restore os.Args and osExit
	defer backupAndRestore(t)()

	// Mock os.Args
	os.Args = []string{
		t.Name(), // dummy app name
		"-v",     // show version
	}

	capturerdStatus := 1 // it shuld turn to 0 if "-v" is set

	// Mock osExit
	osExit = func(code int) {
		capturerdStatus = code

		panic("forced panic") // force panic instead of os.Exit
	}

	out := capturer.CaptureOutput(func() {
		require.Panics(t, func() { main() })
	})

	require.Equal(t, 0, capturerdStatus, "it shuld exit with 0 in case of version flag")
	require.Equal(t, "hostpital (devel)\n", out,
		"it should contain the app name without extension and with it's version")
}

func Test_main_golden_stdout(t *testing.T) {
	// Backup and defer restore os.Args and osExit
	defer backupAndRestore(t)()

	const pathDirFile = "testdata"

	listFiles := []string{
		filepath.Join(pathDirFile, "host1.txt"),
		filepath.Join(pathDirFile, "host2.txt"),
	}

	// Mock os.Args
	os.Args = []string{
		t.Name(),        // dummy app name
		"-l",            // sort by reversed label
		"-i", "0.0.0.0", // append IP address
		listFiles[0], // target file1
		listFiles[1], // target file2
	}

	// Mock osExit
	osExit = func(code int) {
		panic("os.Exit called") // force panic instead of os.Exit
	}

	out := capturer.CaptureOutput(func() {
		assert.NotPanics(t, func() { main() })
	})

	t.Log(out) // log in case of panic

	require.Contains(t, out, heredoc.Doc(`
		0.0.0.0 badboy1.example.com
		0.0.0.0 badboy2.example.com badboy3.example.com
		0.0.0.0 badboy1.example.jp
		0.0.0.0 badboy2.example.jp badboy3.example.jp
	`))
}

func Test_main_golden_file_out(t *testing.T) {
	// Backup and defer restore os.Args and osExit
	defer backupAndRestore(t)()

	pathFileOut := filepath.Join(t.TempDir(), "out.txt")

	const pathDirFile = "testdata"

	listFiles := []string{
		filepath.Join(pathDirFile, "host1.txt"),
		filepath.Join(pathDirFile, "host2.txt"),
	}

	// Mock os.Args
	os.Args = []string{
		t.Name(),          // dummy app name
		"-s",              // sort by host name
		"-o", pathFileOut, // output file
		listFiles[0], // target file1
		listFiles[1], // target file2
	}

	// Mock osExit
	osExit = func(code int) {
		// force to panic instead of os.Exit on error
		panic("unexpected os.Exit was called")
	}

	capturedOut := capturer.CaptureOutput(func() {
		assert.NotPanics(t, func() { main() })
	})

	t.Log(capturedOut)

	require.FileExists(t, pathFileOut)

	outFile, err := os.ReadFile(pathFileOut)
	require.NoError(t, err)

	require.Contains(t, capturedOut, "Output file:")
	require.Contains(t, string(outFile), heredoc.Doc(`
		badboy1.example.com
		badboy1.example.jp
		badboy2.example.com badboy3.example.com
		badboy2.example.jp badboy3.example.jp
	`))
}

// ============================================================================
//  Error Cases
// ============================================================================

// ----------------------------------------------------------------------------
//  appendFileTo()
// ----------------------------------------------------------------------------

func Test_appendFileTo_input_is_empty(t *testing.T) {
	err := appendFileTo("", nil)

	require.Error(t, err, "it should return error on empty input")
	assert.Contains(t, err.Error(), "failed to open the file")
}

func Test_appendFileTo_outfile_fails_to_write(t *testing.T) {
	pathFile := filepath.Join("testdata", "host1.txt")
	dummyFile := new(DummyFile) // dummy implementation to force error

	err := appendFileTo(pathFile, dummyFile)

	require.Error(t, err, "it should return if fails to write")
	assert.Contains(t, err.Error(), "failed to write to the file", "it should contain the error reason")
	assert.Contains(t, err.Error(), "dummy error to write", "it should contain the wrapped error")
}

// ----------------------------------------------------------------------------
//  Flags.ShowHelpAndExitIfTrue()
// ----------------------------------------------------------------------------

func TestFlags_ShowHelpAndExitIfTrue(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	capturedCode := 0 // captured exit code

	// Mock osExit
	osExit = func(code int) {
		capturedCode = code

		panic("os.Exit called") // force panic instead of os.Exit
	}

	out := capturer.CaptureOutput(func() {
		flags, err := ParseFlags()
		require.NoError(t, err)

		assert.PanicsWithValue(t, "os.Exit called", func() { flags.ShowHelpAndExitIfTrue(true, "foced error") })
	})

	require.Equal(t, 1, capturedCode)
	require.Contains(t, out, "foced error")
}

// ----------------------------------------------------------------------------
//  main()
// ----------------------------------------------------------------------------

func Test_main_out_file_is_dir(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	pathDirFile := "testdata"
	pathFileOut := t.TempDir()

	listFiles := []string{
		filepath.Join(pathDirFile, "host1.txt"),
		filepath.Join(pathDirFile, "host2.txt"),
	}

	// Mock os.Args
	os.Args = []string{
		t.Name(),          // dummy app name
		"-s",              // sort by host name
		"-o", pathFileOut, // output file
		listFiles[0], // target file1
		listFiles[1], // target file2
	}

	// Mock osExit
	osExit = func(code int) {
		panic("os.Exit called") // force panic instead of os.Exit
	}

	capturedOut := capturer.CaptureOutput(func() {
		assert.Panics(t, func() { main() })
	})

	require.Contains(t, capturedOut, "failed to create the output file")
}

func Test_main_show_help(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	// Mock os.Args
	os.Args = []string{
		t.Name(), // dummy app name
		"-h",     // show help
	}

	capturedCode := 0 // captured exit code

	// Mock osExit
	osExit = func(code int) {
		capturedCode = code // capture

		panic("forced panic. os.Exit called")
	}

	outStderr := capturer.CaptureStderr(func() {
		assert.Panics(t, func() { main() }, "it should panic on os.Exit call")
	})

	require.Equal(t, 1, capturedCode, "exit code should be 1 on error")
	assert.Contains(t, outStderr, "Merge multiple hosts file(s) into one but parse and sort them.",
		"help message should be shown on stderr")
	assert.Contains(t, outStderr, "Usage:",
		"help message should contain usage")
	assert.Contains(t, outStderr, "Options:",
		"help message should contain options")
}

// ----------------------------------------------------------------------------
//	MergeFiles()
// ----------------------------------------------------------------------------

func TestMergeFiles_golden(t *testing.T) {
	pathFile1 := filepath.Join("testdata", "host1.txt")
	require.FileExists(t, pathFile1, "test data file1 should exist")

	pathFile2 := filepath.Join("testdata", "host2.txt")
	require.FileExists(t, pathFile2, "test data file1 should exist")

	pathTmp, fnTest, err := MergeFiles([]string{pathFile1, pathFile2})

	require.NoError(t, err, "it should not return error on success")

	require.FileExists(t, pathTmp, "it should create a temporary file")

	// Read the file
	content, err := os.ReadFile(pathTmp)
	require.NoError(t, err, "it should read the file")

	assert.Contains(t, string(content), "This is a comment block for host1.txt", "it should contain the content")
	assert.Contains(t, string(content), "This is a comment block for host2.txt", "it should contain the content")

	// Cleanup check
	err = fnTest()
	require.NoError(t, err, "it should not return error on cleanup")

	require.NoFileExists(t, pathTmp, "cleanup function shuld remove the temporary file")
}

func TestMergeFiles_fail_create_temp_file(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	// Mock osCreateTemp
	osCreateTemp = func(dir string, pattern string) (*os.File, error) {
		return nil, errors.New("forced error")
	}

	pathTmp, fnTest, err := MergeFiles([]string{})

	require.Error(t, err, "it should return error on temp file creation failure")
	assert.Contains(t, err.Error(), "failed to create a temporary file")
	assert.Empty(t, pathTmp, "it should return empty path on error")
	assert.Nil(t, fnTest, "it should return nil function on error")
}

func TestMergeFiles_fail_append_to_file(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	tmpDirAsDummyFile, err := os.Open(t.TempDir())
	require.NoError(t, err, "it should open a dummy file")

	defer tmpDirAsDummyFile.Close()

	// Mock osCreateTemp
	osCreateTemp = func(dir string, pattern string) (*os.File, error) {
		return tmpDirAsDummyFile, nil
	}

	pathTmp, fnTest, err := MergeFiles([]string{
		filepath.Join("testdata", "host1.txt"),
	})

	require.Error(t, err, "it should return error on temp file creation failure")
	assert.Contains(t, err.Error(), "failed to append the file")
	assert.Empty(t, pathTmp, "it should return empty path on error")
	assert.Nil(t, fnTest, "it should return nil function on error")
}

// ----------------------------------------------------------------------------
//  NameExec()
// ----------------------------------------------------------------------------

func TestNameExec(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	beforeMock := NameExec()
	t.Log("before mock:", beforeMock)

	// Mock osExecutable
	osExecutable = func() (string, error) {
		return "", errors.New("forced error")
	}

	// Use os.Args[0] as app name
	{
		expectAfter := "dummy-app-name"

		// Mock os.Args
		os.Args = []string{
			expectAfter,
		}

		actualAfter := NameExec()

		require.NotEqual(t, beforeMock, actualAfter)
		require.Equal(t, expectAfter, actualAfter, "if os.Executable() fails, it should return os.Args[0]")
	}

	// Use pre-defined name as the app name
	{
		expectAfter := NameAppDefault

		// Mock os.Args
		os.Args = []string{}

		actualAfter := NameExec()

		require.Equal(t, expectAfter, actualAfter,
			"if os.Executable() fails and os.Args is empty, it should return the default app name")
	}
}

// ----------------------------------------------------------------------------
//  ParseFlags()
// ----------------------------------------------------------------------------

func TestParseFlags(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	// Mock os.Args
	os.Args = []string{}

	flags, err := ParseFlags()

	require.Error(t, err, "it should return error on empty os.Args")
	assert.Contains(t, err.Error(), "failed to parse the flags", "it should return error on empty os.Args")
	assert.Contains(t, err.Error(), "no arguments", "the error should contain the reason")
	assert.Nil(t, flags, "it should return nil flags on error")
}

// ----------------------------------------------------------------------------
//  ShowVerApp
// ----------------------------------------------------------------------------

func TestShowVerApp(t *testing.T) {
	// Backup and defer restore os.Args and function variables
	defer backupAndRestore(t)()

	capturedStatus := 1 // it shuld turn to 0 on success

	// Mock osExit
	osExit = func(code int) {
		capturedStatus = code
	}

	{
		out := capturer.CaptureStdout(func() {
			ShowVerApp()
		})

		require.Equal(t, 0, capturedStatus, "it should exit with status 0")
		assert.Contains(t, out, "hostpital (devel)",
			"it should contain the app name with no extension with the version")
	}

	// version is set (via build args)
	{
		version = "v1.2.3456789" // pretend that version is set via build args

		out := capturer.CaptureStdout(func() {
			ShowVerApp()
		})

		expect := "hostpital v1.2.3456789"

		assert.Contains(t, out, expect, "if version variable is set, it should use that version")
	}
}

// ============================================================================
//  Helper Functions and Types for Testing
// ============================================================================

// ----------------------------------------------------------------------------
//  Type DummyFile
// ----------------------------------------------------------------------------

// DummyFile implements main.IOFile interface to return dummy error.
type DummyFile struct{}

// Read implements main.IOFile.Read interface to return dummy error.
func (d *DummyFile) Read(p []byte) (int, error) {
	return 0, errors.New("dummy error to read")
}

// Write implements main.IOFile.Write interface to return dummy error.
func (d *DummyFile) Write(p []byte) (int, error) {
	return 0, errors.New("dummy error to write")
}

// ----------------------------------------------------------------------------
//  backupAndRestore()
// ----------------------------------------------------------------------------

// backupAndRestore backups the values of os.Args and global function variables
// and returns a function to restore the original values.
//
// To avoid side effects, the caller test functions must follow the below two
// rules:
//  1. MUST NOT PARALLELIZE due to the temporary change of the global variables.
//  2. MUST defer execute the returned function to restore the original values.
func backupAndRestore(t *testing.T) func() {
	t.Helper()

	oldArgs := os.Args
	oldOsExit := osExit
	oldOsExecutable := osExecutable
	oldOsCreateTemp := osCreateTemp
	oldVersion := version

	return func() {
		os.Args = oldArgs
		osExit = oldOsExit
		osExecutable = oldOsExecutable
		osCreateTemp = oldOsCreateTemp
		version = oldVersion
	}
}
