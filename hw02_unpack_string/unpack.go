package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	err := validate(s)
	if err != nil {
		return "", err
	}

	rs := []rune(s)
	var builder strings.Builder
	for i, r := range rs {
		if isLastRuneInString(i, rs) && unicode.IsDigit(r) {
			continue
		} else if isLastRuneInString(i, rs) || !unicode.IsDigit(r) && !unicode.IsDigit(rs[i+1]) {
			builder.WriteRune(r)
		} else {
			digit, _ := strconv.Atoi(string(rs[i+1]))
			for j := 0; j < digit; j++ {
				builder.WriteRune(r)
			}
		}
	}

	return builder.String(), nil
}

func isLastRuneInString(i int, rs []rune) bool {
	return i == len(rs)-1
}

func validate(s string) error {
	for i, r := range s {
		if i == 0 && unicode.IsDigit(r) {
			return ErrInvalidString
		}
		if unicode.IsDigit(r) && unicode.IsDigit(rune(s[i-1])) {
			return ErrInvalidString
		}
	}
	return nil
}
