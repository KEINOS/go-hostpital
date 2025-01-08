package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Benchmark to find the fastest way to trim the extension from a file name.
//
//  It must trim "hostpital.test.exe" to "hostpital" on Windows in the CI.
// ----------------------------------------------------------------------------

func trimExt1(fileName string) string {
	trimmed := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	if fileName != trimmed {
		return trimExt1(trimmed)
	}

	return trimmed
}

func trimExt2(fileName string) string {
	trimmed := fileName[:len(fileName)-len(filepath.Ext(fileName))]

	if fileName != trimmed {
		return trimExt2(trimmed)
	}

	return trimmed
}

func trimExt3(fileName string) string {
	const delimiter = '.'

	foundIndex := strings.IndexByte(fileName, delimiter)
	if foundIndex != -1 {
		return fileName[:foundIndex]
	}

	return fileName
}

func Benchmark_trimExt1(b *testing.B) {
	const fileName = "hostpital.test.exe"

	require.Equal(b, "hostpital", trimExt1(fileName))

	b.ResetTimer()

	for range b.N {
		_ = trimExt1(fileName)
	}
}

func Benchmark_trimExt2(b *testing.B) {
	const fileName = "hostpital.test.exe"

	require.Equal(b, "hostpital", trimExt2(fileName))

	b.ResetTimer()

	for range b.N {
		_ = trimExt2(fileName)
	}
}

func Benchmark_trimExt3(b *testing.B) {
	const fileName = "hostpital.test.exe"

	require.Equal(b, "hostpital", trimExt3(fileName))

	b.ResetTimer()

	for range b.N {
		_ = trimExt3(fileName)
	}
}
