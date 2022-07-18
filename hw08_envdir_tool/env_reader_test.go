package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	BAR EnvValue = EnvValue{"bar", false}
	FOO string   = `   foo
with new line`
	HELLO string = `"hello"`
	UNSET string = ""
)

func TestReadDir(t *testing.T) {
	var expectedMap = Environment{
		"BAR":   EnvValue{Value: "BAR", NeedRemove: false},
		"FOO":   EnvValue{Value: "FOO", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
		"EMPTY": EnvValue{Value: " \n", NeedRemove: false},
	}
	const dir = "./testdata/env"
	env, err := ReadDir(dir)
	assert.NoError(t, err)

	assert.Len(t, env, 4)
	assert.Equal(t, expectedMap, env)
}
