package hostpital

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  ValidateFile
// ----------------------------------------------------------------------------

func TestValidator_ValidateFile_path_not_exist(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	path := filepath.Join(t.TempDir(), "not_exist")
	ok := validator.ValidateFile(path)

	require.False(t, ok, "it should return false if the path is a directory")
}

func TestValidator_ValidateFile_path_is_dir(t *testing.T) {
	t.Parallel()

	pathDir := t.TempDir()
	validator := NewValidator()

	validator.AllowEmptyLine = false

	ok := validator.ValidateFile(pathDir)

	require.False(t, ok, "it should return false if the path is a directory")
}

func TestValidator_ValidateFile_malformed_file(t *testing.T) {
	t.Parallel()

	pathFile := filepath.Join("testdata", "default.txt")
	validator := NewValidator()

	validator.AllowComment = false

	ok := validator.ValidateFile(pathFile)

	require.False(t, ok,
		"it should return false if the file is malformed according to the configuration")
}

// ----------------------------------------------------------------------------
//  ValidateLine
// ----------------------------------------------------------------------------

func TestValidator_ValidateLine_allow_underscore(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	// Hostname that the label contains underscore
	line := "0.0.0.0 foo_bar.example.com"

	// Disallow underscore
	{
		validator.AllowUnderscore = false // default setting

		err := validator.ValidateLine(line)
		require.Error(t, err, "it should return an error if the line contains underscore by default")
	}

	// Allow underscore
	{
		validator.AllowUnderscore = true

		err := validator.ValidateLine(line)
		require.NoError(t, err, "it should not error by the configuration")
	}
}

func TestValidator_ValidateLine_allow_double_hyphen(t *testing.T) {
	t.Parallel()

	// Hostname that contains malformed punycode with double hyphen
	const line = "127.0.0.0 123.--foo.bar.example.com"

	validator := NewValidator()

	// Disallow host that begins with number
	{
		validator.AllowHyphenDouble = false // default setting

		err := validator.ValidateLine(line)

		require.Error(t, err, "it should return an error if host begins with number")
		require.Contains(t, err.Error(), "is not IDNA2008 compatible", "error should contain the reason")
		require.Contains(t, err.Error(), "idna: invalid label", "error should contain the reason")
	}

	// Allow host that begins with number
	{
		validator.AllowHyphenDouble = true

		err := validator.ValidateLine(line)
		require.NoError(t, err, "it should not error by the configuration")
	}
}

func TestValidator_ValidateLine_trailing_space(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	err := validator.ValidateLine("0.0.0.0 example.com         ")

	require.Error(t, err, "it should return an error if the line contains trailing space by default")
	assert.Contains(t, err.Error(), "trailing space is not allowed")
}

func TestValidator_ValidateLine_disallow_empty_line(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	validator.AllowEmptyLine = false // Disallow empty line

	err := validator.ValidateLine("")

	require.Error(t, err, "it should return an error if the line is empty by the configuration")
	assert.Contains(t, err.Error(), "empty line is not allowed")
}

func TestValidator_ValidateLine_fail_trim_comment(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	validator.AllowComment = true // Allow comment

	err := validator.ValidateLine("127.0.0.0 example.com \n # comment")

	require.Error(t, err, "it should return an error if the input is multiple lines")
	assert.Contains(t, err.Error(), "failed to trim comment")
}

func TestValidator_ValidateLine_ip_address_only_line(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	validator.AllowComment = true        // Allow comment
	validator.AllowIPAddressOnly = false // Disallow IP address only line

	err := validator.ValidateLine("127.0.0.0 # comment")

	require.Error(t, err, "it should return an error if the input is multiple lines")
	assert.Contains(t, err.Error(), "IP address only line is not allowed")
}

func TestValidator_ValidateLine_incompatible_to_IDNA(t *testing.T) {
	t.Parallel()

	validator := NewValidator()

	validator.IDNACompatible = false

	err := validator.ValidateLine("-eXample123-.com")

	require.Error(t, err, "it should return an error if the input is multiple lines")
	assert.Contains(t, err.Error(), "not RFC 6125 2.2 compatible")
}
