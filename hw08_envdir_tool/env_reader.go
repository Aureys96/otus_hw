package main

import (
	"bytes"
	"errors"
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
		if key, value, err := extractValue(dir, file); err != nil {
			return nil, err
		} else {
			if value == "" {
				env[key] = EnvValue{value, true}
			} else {
				env[key] = EnvValue{value, false}
			}
		}
	}

	return env, nil
}

func extractValue(dir string, file os.DirEntry) (string, string, error) {
	key := file.Name()
	if strings.Contains(key, "=") {
		return "", "", ErrWrongFileName
	}

	filename := filepath.Join(dir, key)
	value, err := parseFile(filename)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func parseFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", ErrCannotReadFile
	}

	if len(content) == 0 {
		return "", nil
	}

	content = bytes.ReplaceAll(content, []byte("\x00"), []byte("\n"))
	return string(content), nil
}
