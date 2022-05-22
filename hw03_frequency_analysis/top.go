package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile("([.,/#!$%^&*;:{}=_~()])|(^-)+")

func Top10(input string) []string {
	splitStrings := strings.Fields(input)
	freqMap := make(map[string]int)
	for _, word := range splitStrings {
		word = re.ReplaceAllString(strings.ToLower(word), "")
		if word == "" {
			continue
		}
		if _, ok := freqMap[word]; ok {
			freqMap[word]++
		} else {
			freqMap[word] = 1
		}
	}

	tuples := make([]tuple, 0, len(freqMap))
	for word, count := range freqMap {
		tuples = append(tuples, tuple{
			word:  word,
			count: count,
		})
	}

	sort.Slice(tuples, func(i, j int) bool {
		if tuples[i].count > tuples[j].count {
			return true
		}
		if tuples[i].count < tuples[j].count {
			return false
		}
		return tuples[i].word < tuples[j].word
	})

	var answer []string
	for i, t := range tuples {
		if i == 10 {
			break
		}
		answer = append(answer, t.word)
	}

	return answer
}

type tuple struct {
	word  string
	count int
}
