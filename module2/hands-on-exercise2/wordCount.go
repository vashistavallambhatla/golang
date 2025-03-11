package main

import (
	"regexp"
	"strings"
)

type FrequencyMap map[string]int

func WordCount(s string) FrequencyMap {
	result := make(FrequencyMap)
	reg := regexp.MustCompile(`\w+('\w+)?`)

	for _, word := range reg.FindAllString(strings.ToLower(s), -1) {
		result[word]++
	}

	return result
}

