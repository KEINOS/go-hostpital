package hostpital_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-hostpital/hostpital"
)

func ExampleMergeSortFiles() {
	inFiles := []string{
		filepath.Join("testdata", "sort_me", "unsorted_file1.txt"),
		filepath.Join("testdata", "sort_me", "unsorted_file2.txt"),
		filepath.Join("testdata", "sort_me", "unsorted_file3.txt"),
	}

	outFile := filepath.Join(os.TempDir(), "hostpital-merge_sort_files-sorted_file.txt")
	defer func() {
		if err := os.Remove(outFile); err != nil {
			log.Fatal(err)
		}
	}()

	// Hosts file parser with default settings.
	// No-comment, no-empty lines, no-leading IP address, no-duplicate lines.
	parser := hostpital.NewParser()

	// Sort the combined file after parsing.
	parser.SortAfterParse = true

	err := hostpital.MergeSortFiles(inFiles, parser, outFile)
	if err != nil {
		panic(err) // panic to trigger the defer
	}

	byteOut, err := os.ReadFile(outFile)
	if err != nil {
		panic(err) // panic to trigger the defer
	}

	fmt.Println(string(byteOut))
	// Output:
	// alice.com
	// bob.jp
	// carol.jp
	// charlie.com
	// dave.com
	// ellen.com
	// eve.jp
	// frank.com
	// isaac.com
	// ivan.com
	// justin.jp
	// mallet.com
	// mallory.com
	// marvin.com
	// matilda.jp
	// oscar.com
	// pat.com
	// peggy.jp
	// steve.com
	// trent.jp
	// trudy.com
	// victor.com
	// walter.jp
	// zoe.com
}

func ExampleMergeSortFiles_sort_by_reversed_dns() {
	inFiles := []string{
		filepath.Join("testdata", "sort_me", "unsorted_file1.txt"),
		filepath.Join("testdata", "sort_me", "unsorted_file2.txt"),
		filepath.Join("testdata", "sort_me", "unsorted_file3.txt"),
	}

	outFile := filepath.Join(os.TempDir(), "hostpital-merge_sort_files-sorted_file.txt")
	defer func() {
		if err := os.Remove(outFile); err != nil {
			log.Fatal(err)
		}
	}()

	// Hosts file parser with default settings.
	// No-comment, no-empty lines, no-leading IP address, no-duplicate lines.
	parser := hostpital.NewParser()

	parser.SortAfterParse = true

	// If true, sort by reversed DNS. Which means "www.example.com" will be
	// treated as "com.example.www" while sorting.
	parser.SortAsReverseDNS = true

	err := hostpital.MergeSortFiles(inFiles, parser, outFile)
	if err != nil {
		panic(err) // panic to trigger the defer
	}

	byteOut, err := os.ReadFile(outFile)
	if err != nil {
		panic(err) // panic to trigger the defer
	}

	fmt.Println(string(byteOut))
	// Output:
	// alice.com
	// charlie.com
	// dave.com
	// ellen.com
	// frank.com
	// isaac.com
	// ivan.com
	// mallet.com
	// mallory.com
	// marvin.com
	// oscar.com
	// pat.com
	// steve.com
	// trudy.com
	// victor.com
	// zoe.com
	// bob.jp
	// carol.jp
	// eve.jp
	// justin.jp
	// matilda.jp
	// peggy.jp
	// trent.jp
	// walter.jp
}
