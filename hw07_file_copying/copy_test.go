package main

import (
	"crypto/md5"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCopy(t *testing.T) {
	var pwd, _ = os.Getwd()

	originalFile := path.Join(pwd, "testdata", "input.txt")

	tests := []struct {
		name         string
		originalFile string
		expectedFile string
		offset       int64
		limit        int64
	}{
		{name: "offset 0 limit 0",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset0_limit0.txt"),
			offset:       0,
			limit:        0},
		{name: "offset 0 limit 10",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset0_limit10.txt"),
			offset:       0,
			limit:        10},
		{name: "offset 0 limit 1000",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset0_limit1000.txt"),
			offset:       0,
			limit:        1000},
		{name: "offset 0 limit 10000",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset0_limit10000.txt"),
			offset:       0,
			limit:        10000},
		{name: "offset 100 limit 1000",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset100_limit1000.txt"),
			offset:       100,
			limit:        1000},
		{name: "offset 6000 limit 1000",
			originalFile: originalFile,
			expectedFile: path.Join(pwd, "testdata", "out_offset6000_limit1000.txt"),
			offset:       6000,
			limit:        1000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			file, err := os.CreateTemp(pwd+"/testdata", "goTest*.txt")
			if err != nil {
				assert.NoError(t, err, "Couldn't create temp file")
			}
			defer os.Remove(file.Name())

			err = Copy(test.originalFile, file.Name(), test.offset, test.limit)
			if err != nil {
				assert.NoError(t, err, "Error occurred while copying")
			}

			expectedFileHash := getFileHash(test.expectedFile)
			actualFileHash := getFileHash(file.Name())

			assert.Equal(t, expectedFileHash, actualFileHash)

		})
	}

}

func getFileHash(filename string) []byte {
	hash := md5.New()
	file, _ := ioutil.ReadFile(filename)
	return hash.Sum(file)
}
