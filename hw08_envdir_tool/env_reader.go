package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrCannotOpenDir  = errors.New("cannot open directory")
	ErrWrongFileName  = errors.New("invalid '=' symbol in filename")
	ErrCannotReadFile = errors.New("cannot read file")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrCannotOpenDir
	}

	env := make(Environment)
	for _, file := range files {
		key, value, err := extractValue(dir, file)
		if err != nil {
			return nil, err
		}
		env[key] = *value
	}

	return env, nil
}

func extractValue(dir string, file os.DirEntry) (string, *EnvValue, error) {
	key := file.Name()
	if strings.Contains(key, "=") {
		return "", nil, ErrWrongFileName
	}

	filename := filepath.Join(dir, key)
	value, isEmpty, err := parseFile(filename)
	if err != nil {
		return "", nil, err
	}
	if isEmpty {
		return key, &EnvValue{value, true}, nil
	}

	return key, &EnvValue{value, false}, nil
}

func parseFile(filename string) (string, bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", false, ErrCannotReadFile
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", false, err
	}
	if stat.Size() == 0 {
		return "", true, nil
	}

	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", false, nil
		}
		return "", false, err
	}
	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))
	result := string(line)
	result = strings.TrimRight(result, "	")
	result = strings.TrimRight(result, " ")
	return result, false, nil
}
