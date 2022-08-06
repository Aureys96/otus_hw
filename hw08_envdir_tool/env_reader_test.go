package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	expectedMap := Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: false},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}
	const dir = "./testdata/env"
	env, err := ReadDir(dir)
	assert.NoError(t, err)
	assert.Equal(t, expectedMap, env)
}
