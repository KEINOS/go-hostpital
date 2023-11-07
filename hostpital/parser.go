package hostpital

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/KEINOS/go-countline/cl"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

// ----------------------------------------------------------------------------
//  Type: Parser
// ----------------------------------------------------------------------------

// Parser holds the settings and the rules for the parsing. To simply validate
// the hostfile, use the methods in the Validator type instead.
type Parser struct {
	UseIPAddress      string // If not empty and 'TrimIPAddress' is true, use this IP address instead (default: "").
	mutx              sync.Mutex
	IDNACompatible    bool // If true, punycode is converted to IDNA2008 compatible (default: true).
	OmitEmptyLine     bool // If true, empty lines are omitted (default: true).
	SortAfterParse    bool // If true, sort the lines after parsing (default: false).
	SortAsReverseDNS  bool // If true, sort the lines as reversed DNS hosts (default: false).
	TrimComment       bool // If true, comment is trimmed (default: true).
	TrimIPAddress     bool // If true, leading IP address is trimmed (default: true).
	TrimLeadingSpace  bool // If true, leading spaces are trimmed (default: true).
	TrimTrailingSpace bool // If true, trailing spaces are trimmed (default: true).
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewParser returns a new Parser instance with the default values.
func NewParser() *Parser {
	parser := new(Parser) // Set all to false

	// Set default values. Non mentioned values are set to false.
	parser.UseIPAddress = ""
	parser.IDNACompatible = true
	parser.OmitEmptyLine = true
	parser.TrimComment = true
	parser.TrimIPAddress = true
	parser.TrimLeadingSpace = true
	parser.TrimTrailingSpace = true

	return parser
}

// ----------------------------------------------------------------------------
//  Methods (Public)
// ----------------------------------------------------------------------------

// CountLines counts the number of lines in the file.
func (p *Parser) CountLines(pathFile string) (int, error) {
	p.mutx.Lock()
	defer p.mutx.Unlock()

	osFile, err := os.Open(pathFile)
	if err != nil {
		return 0, errors.Wrap(err, "failed to open the file")
	}

	defer osFile.Close()

	numLines, err := cl.CountLines(osFile)

	return numLines, errors.Wrap(err, "failed to count lines")
}

// ParseFile reads the file from pathFile and returns the parsed lines as a string
// according to the settings in the Parser.
func (p *Parser) ParseFile(pathFile string) (string, error) {
	outBuf := new(bytes.Buffer)

	if err := p.ParseFileTo(pathFile, outBuf); err != nil {
		return "", errors.Wrap(err, "failed to parse the file")
	}

	return outBuf.String(), nil
}

// ParseFileTo reads the file from pathFileIn and writes the parsed lines to fileOut.
func (p *Parser) ParseFileTo(pathFileIn string, fileOut io.Writer) error {
	if fileOut == nil {
		return errors.New("the given io.Writer is nil")
	}

	// Prepare a slice to store the parsed lines.
	numLines, err := p.CountLines(pathFileIn)
	if err != nil {
		return errors.Wrap(err, "failed to read and count lines from the input file")
	}

	lines := make([]string, numLines)

	// Open the file.
	// Error check is omitted because it is done in the above p.CountLines() so
	// it will never reach here if the file does not exist or is not readable.
	osFile, _ := osOpen(pathFileIn)
	defer osFile.Close()

	// Returned error not checked as it is done in the above p.CountLines().
	_ = p.scanFile(osFile, lines)

	if p.SortAfterParse || p.SortAsReverseDNS {
		lines = p.sortSlices(lines)
	}

	for _, line := range lines {
		if _, err = fileOut.Write([]byte(line)); err != nil {
			return errors.Wrap(err, "failed to write to io.Writer")
		}
	}

	return nil
}

// ParseLine parses the given line and returns the parsed line as a string
// according to the settings in the Parser.
func (p *Parser) ParseLine(line string) (string, bool) {
	trimmed := line

	trimmed = p.trimSpace(trimmed)
	trimmed = p.trimComment(trimmed)
	trimmed = p.trimIPAddress(trimmed)
	trimmed = p.onlyIDNACompatible(trimmed)

	if p.OmitEmptyLine && strings.TrimSpace(trimmed) == "" {
		return "", false
	}

	trimmed = p.prependIPAddress(trimmed)

	return trimmed, true
}

// ParseString parses the given string and returns the parsed lines as a string
// according to the settings in the Parser.
func (p *Parser) ParseString(input string) string {
	lines := strings.Split(input, string(LF))
	parsed := make([]string, len(lines))

	for index, line := range lines {
		if trimmed, ok := p.ParseLine(line); ok {
			parsed[index] = trimmed
		}
	}

	if p.SortAfterParse || p.SortAsReverseDNS {
		parsed = p.sortSlices(parsed)
	}

	return strings.Join(parsed, string(LF))
}

// ----------------------------------------------------------------------------
//  Methods (Private)
// ----------------------------------------------------------------------------

func (p *Parser) onlyIDNACompatible(line string) string {
	if !p.IDNACompatible {
		return line
	}

	if IsCommentLine(line) {
		return line
	}

	trimmed := strings.Split(TrimWordGaps(line), " ")

	for index, chunk := range trimmed {
		hostASCII, err := TransformToASCII(chunk)

		if err != nil || !IsCompatibleIDNA2008(hostASCII) {
			hostASCII = ""
		}

		trimmed[index] = hostASCII
	}

	return strings.Join(trimmed, " ")
}

func (p *Parser) prependIPAddress(line string) string {
	if IsCommentLine(line) || p.UseIPAddress == "" || !p.TrimIPAddress {
		return line
	}

	if p.UseIPAddress != "" && IsIPAddress(line) {
		return line
	}

	return p.UseIPAddress + " " + line
}

func (p *Parser) scanFile(inFile io.Reader, lines []string) error {
	// Prepare reading the file.
	countLines := 0
	wgrp := new(sync.WaitGroup)
	scanBuf := bufio.NewScanner(inFile)

	for scanBuf.Scan() {
		line := scanBuf.Text()

		wgrp.Add(1)

		go func(line string, index int) {
			defer wgrp.Done()

			parsed, ok := p.ParseLine(line)
			if ok {
				lines[index] = parsed + string(LF)
			}
		}(line, countLines)

		countLines++
	}

	if scanBuf.Err() != nil {
		return errors.Wrap(scanBuf.Err(), "failed to read/scan the file")
	}

	wgrp.Wait()

	return nil
}

// sortAsReverseDNS sorts the given slice as reversed DNS hosts.
// For example, "www.example.com" will be sorted as "com.example.www".
// Use `sortSlices` to sort the slice as normal.
func (p *Parser) sortAsReverseDNS(lines []string) []string {
	if !p.SortAsReverseDNS {
		return lines
	}

	// Sort as reversed DNS hosts.
	slices.SortFunc(lines, func(a string, b string) int {
		return strings.Compare(ReverseDNS(a), ReverseDNS(b))
	})

	return lines
}

func (p *Parser) sortSlices(lines []string) []string {
	if p.SortAsReverseDNS {
		return p.sortAsReverseDNS(lines)
	}

	slices.Sort(lines)

	return lines
}

func (p *Parser) trimComment(line string) string {
	if !p.TrimComment {
		return line
	}

	trimmed := strings.TrimLeft(line, Cutset)

	trimmed, _ = TrimComment(trimmed) // Error is ignored as it is already right trimmed.

	return trimmed
}

func (p *Parser) trimSpace(line string) string {
	trimmed := line

	if p.TrimLeadingSpace && p.TrimTrailingSpace {
		return strings.TrimSpace(trimmed)
	}

	if p.TrimLeadingSpace {
		trimmed = strings.TrimLeft(trimmed, Cutset)
	}

	if p.TrimTrailingSpace {
		trimmed = strings.TrimRight(trimmed, Cutset)
	}

	return trimmed
}

func (p *Parser) trimIPAddress(line string) string {
	if IsCommentLine(line) || !p.TrimIPAddress {
		return line
	}

	return TrimIPAdd(line)
}
