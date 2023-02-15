package hostpital

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/KEINOS/go-sortfile/sortfile"
	"github.com/pkg/errors"
)

// MergeSortFiles merges the files and sorts the lines.
func MergeSortFiles(inFiles []string, parser *Parser, outFile string) error {
	pathFilesParsed, cleanUp, err := parseFiles(inFiles, parser)
	if err != nil {
		return errors.Wrap(err, "failed to combine the files")
	}

	defer cleanUp()

	pathFileTmp, err := combineFiles(pathFilesParsed)
	if err != nil {
		return errors.Wrap(err, "failed to combine the parsed files")
	}

	if !parser.SortAfterParse {
		err := os.Rename(pathFileTmp, outFile)

		return errors.Wrap(err, "failed to move the generated file")
	}

	forceExternalSort := false // auto detect in-memory sort or external sort
	isLess := func(a, b string) bool {
		return a < b
	}

	if parser.SortAsReverseDNS {
		isLess = func(a, b string) bool {
			aRev := ReverseDNS(a)
			bRev := ReverseDNS(b)

			return aRev < bRev
		}
	}

	err = sortfile.FromPathFunc(pathFileTmp, outFile, forceExternalSort, isLess)

	return errors.Wrap(err, "failed to sort the file")
}

func parseFiles(inFiles []string, parser *Parser) ([]string, func() error, error) {
	numFiles := len(inFiles)
	wg := sync.WaitGroup{}
	pathFilesParsed := make([]string, numFiles)

	// Create parsed files to temp dir
	wg.Add(numFiles)
	for index, pathFile := range inFiles {
		ptrOut, err := os.CreateTemp(os.TempDir(), "hostpital-parsed*")
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create a temporary file")
		}

		pathFilesParsed[index] = ptrOut.Name()

		go func(pathFileIn string, parser *Parser, pathFileOut io.Writer) {
			defer func() {
				if rec := recover(); rec != nil {
					fmt.Fprintf(os.Stderr, "Panic: %v (recovered)", rec)
				}

				wg.Done()
			}()

			if err := parser.ParseFileTo(pathFileIn, pathFileOut); err != nil {
				panic(err)
			}
			if err := ptrOut.Close(); err != nil {
				panic(err)
			}
		}(pathFile, parser, ptrOut)
	}

	wg.Wait()

	return pathFilesParsed, func() error {
		// Clean up the temporary file
		for _, pathFile := range pathFilesParsed {
			if err := os.Remove(pathFile); err != nil {
				return errors.Wrap(err, "failed to remove the temp parsed file")
			}
		}
		return nil
	}, nil
}

// combineFiles combines the files into one file.
// It returns the temporary path to the combined file. It is the caller's
// responsibility to remove the temp file.
func combineFiles(inFiles []string) (string, error) {
	ptrOut, err := os.CreateTemp(os.TempDir(), "hostpital-*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create a temporary file")
	}

	defer ptrOut.Close()

	for _, pathFile := range inFiles {
		ptrIn, err := os.Open(pathFile)
		if err != nil {
			return "", errors.Wrap(err, "failed to open the temp parsed file")
		}

		if _, err := io.Copy(ptrOut, ptrIn); err != nil {
			return "", errors.Wrap(err, "failed to copy the temp parsed file to output")
		}

		if err := ptrIn.Close(); err != nil {
			return "", errors.Wrap(err, "failed to close the temp parsed file")
		}
	}

	return ptrOut.Name(), nil
}
