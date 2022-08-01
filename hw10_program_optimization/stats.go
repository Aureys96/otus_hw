package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsonIterator "github.com/json-iterator/go"
)

var (
	json             = jsonIterator.ConfigCompatibleWithStandardLibrary
	ErrUnmarshalling = errors.New("error while unmarshalling data")
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	user := User{}
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, ErrUnmarshalling
		}
		if strings.Contains(user.Email, "."+domain) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}

	return result, nil
}
