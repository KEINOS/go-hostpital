package hostpital

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/net/idna"
)

// ----------------------------------------------------------------------------
//  Type: Validator
// ----------------------------------------------------------------------------

// Validator holds the settings and the rules for the validation. To clean up
// the hostfile, use the methods in the Parser type instead.
//
// It is recommended to use NewValidator() to create a new Validator due to the
// default values.
type Validator struct {
	mutx               sync.Mutex
	AllowComment       bool // If true, the line can be a comment (default: false).
	AllowEmptyLine     bool // If true, empty line returns true (default: true).
	AllowHyphen        bool // If true, the label can begin with hyphen (default: false).
	AllowHyphenDouble  bool // If true, unconvertable punycode with double hyphen is allowed (default: false).
	AllowIndent        bool // If true, the line can be indented (default: false).
	AllowIPAddressOnly bool // If true, the line can be only an IP address (default: false).
	AllowTrailingSpace bool // If true, the line can have trailing spaces (default: false).
	AllowUnderscore    bool // If true, the label can have underscore (default: false).
	IDNACompatible     bool // If true, the host must be compatible to IDNA2008 and false to RFC 6125 2.2 (default: true).
	isInitialized      bool
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewValidator returns a new Validator instance with the default values.
func NewValidator() *Validator {
	validator := new(Validator) // Set all to false

	validator.initialize()

	return validator
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// initialize sets the default values.
func (v *Validator) initialize() {
	// Set default values
	v.IDNACompatible = true
	v.AllowEmptyLine = true
	v.isInitialized = true
}

// validateChunk returns nil if the given chunk/part of line is valid according
// to the settings.
//
// This method validates RFC 6125 2.2 and IDNA2008 compatibility.
func (v *Validator) validateChunk(chunk string) error {
	if _, err := TransformToUnicode(chunk); err != nil && v.AllowHyphenDouble {
		chunk = strings.ReplaceAll(chunk, "--", "aa")
	}

	if strings.Contains(chunk, ".-") && v.AllowHyphen {
		chunk = strings.ReplaceAll(chunk, ".-", ".a")
	}

	// RFC 6125 2.2 compatible
	if !v.IDNACompatible && !IsCompatibleRFC6125(chunk) {
		return errors.Errorf("%#v is not RFC 6125 2.2 compatible", chunk)
	}

	// IDNA2008 compatible
	if v.IDNACompatible && !IsCompatibleIDNA2008(chunk) {
		if _, err := idna.Registration.ToASCII(chunk); err != nil {
			msgErr := fmt.Sprintf("%#v is not IDNA2008 compatible", chunk)

			return errors.Wrap(err, msgErr)
		}
	}

	return nil
}

// ValidateFile returns true if the file is valid according to the settings.
func (v *Validator) ValidateFile(pathFile string) bool {
	v.mutx.Lock()
	defer v.mutx.Unlock()

	osFile, err := os.Open(pathFile)
	if err != nil {
		return false
	}

	defer osFile.Close()

	if !IsExistingFile(pathFile) {
		return false
	}

	scanner := bufio.NewScanner(osFile)

	for scanner.Scan() {
		line := scanner.Text()

		if err := v.ValidateLine(line); err != nil {
			log.Println("invalid line:", line, err)

			return false
		}
	}

	return true
}

func (v *Validator) trimLine(line string) (string, error) {
	trimmed := strings.TrimLeft(line, " \t")

	if !v.AllowIndent && trimmed != line {
		return "", errors.New("indent is not allowed")
	}

	trimmed = strings.TrimRight(line, " \t")

	if !v.AllowTrailingSpace && trimmed != line {
		return "", errors.New("trailing space is not allowed")
	}

	trimmed = strings.TrimSpace(line)

	if trimmed == "" {
		if !v.AllowEmptyLine {
			return "", errors.New("empty line is not allowed")
		}
	}

	return trimmed, nil
}

// ValidateLine returns nil if the line is valid according to the settings.
func (v *Validator) ValidateLine(line string) error {
	trimmed, err := v.trimLine(line)
	if err != nil {
		return errors.Wrap(err, "failed to trim line")
	}

	if trimmed == "" {
		return nil
	}

	if v.AllowComment {
		if IsCommentLine(trimmed) {
			return nil
		}

		noComment, err := TrimComment(trimmed)
		if err != nil {
			return errors.Wrap(err, "failed to trim comment")
		}

		if trimmed != noComment {
			return v.ValidateLine(noComment)
		}
	}

	if !v.AllowIPAddressOnly && IsIPAddress(trimmed) {
		return errors.New("IP address only line is not allowed")
	}

	if v.AllowUnderscore {
		trimmed = strings.ReplaceAll(trimmed, "_", "-")
	}

	for _, chunk := range strings.Split(TrimWordGaps(trimmed), " ") {
		if err := v.validateChunk(chunk); err != nil {
			return errors.Wrap(err, "failed to validate chunk/part of line")
		}
	}

	return nil
}
