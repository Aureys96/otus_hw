package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("return code 0", func(t *testing.T) {
		file, err := os.Create("testdata/text_1.txt")
		if err != nil {
			assert.NoError(t, err, "Couldn't create temp file")
		}
		defer os.Remove(file.Name())

		_, err = file.Write([]byte("test data"))
		if err != nil {
			assert.NoError(t, err, "Couldn't write to temp file")
		}
		file.Close()

		cmd := []string{"cat", "testdata/text_1.txt"}
		env := Environment{}
		returnCode := RunCmd(cmd, env)
		require.Equal(t, 0, returnCode)
	})
	t.Run("return code not equal 0", func(t *testing.T) {
		cmd := []string{"cat", "text_2.txt"}
		env := Environment{}
		returnCode := RunCmd(cmd, env)
		require.NotEqual(t, 0, returnCode)
	})
}
