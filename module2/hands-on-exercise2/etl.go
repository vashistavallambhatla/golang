package main

import (
	"fmt"
	"strings"
)

func transform(m map[int][]string) map[string]int {
	res := make(map[string]int)
	for key,value := range m {
		for _,str := range value {
			res[strings.ToLower(str)] = key
		}
	}
	return res
}

func etl() {
	example := map[int][]string{1 : {"A","E","I","O","U"}}
	fmt.Println(transform(example))
}

