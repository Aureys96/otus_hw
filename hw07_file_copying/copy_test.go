package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var pwd, _ = os.Getwd()

func TestCopy(t *testing.T) {
	Copy(pwd+"/testdata/from/input.txt", pwd+"/testdata/to/dest.txt", 0, 0)
}

func TestCopy2(t *testing.T) {
	pwd, _ := os.Getwd()
	file, err := os.CreateTemp(pwd+"/testdata", "goTest*.txt")
	if err != nil {
		assert.NoError(t, err, "Couldn't create temp file")
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(file.Name())

	err = Copy(pwd+"/testdata/input.txt", file.Name(), 4, 9)
	if err != nil {
		assert.NoError(t, err, "Error occurred while copying")
	}

	err = file.Close()
	if err != nil {
		assert.NoError(t, err, "Error occurred while closing temp file")
	}

}
