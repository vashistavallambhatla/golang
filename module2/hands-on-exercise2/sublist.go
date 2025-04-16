package main

import (
	"fmt"
)

func CustomDeepEqual(a []int,b []int) bool { // Wrote a custom DeepEqual that basically checks if the arrays are identical
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i]!=b[i] {
			return false
		}
	}
	return true
}

func findRelation(a []int,b []int) string {
	if CustomDeepEqual(a,b) {
		return "equal"
	} else if len(a) > len(b) {
		if(isSubList(a,b)){
			return "Superlist"
		}
	} else if len(b) > len(a) {
		if(isSubList(b,a)){
			return "Sublist"
		}
	} 
	return "Unequal"
}

func isSubList(big,small []int) bool {
	small_len := len(small)
	if small_len == 0 {
		return true
	}
	for i := 0 ; i <= len(big)-small_len ; i++ {
		if CustomDeepEqual(big[i:i+small_len],small){
			return true
		}
	}
	return false
}


func sublist() {
	testCases := []struct {
		A, B      []int
		expected  string
	}{
		{[]int{}, []int{}, "equal"},
		{[]int{1, 2, 3}, []int{}, "superlist"},
		{[]int{}, []int{1, 2, 3}, "sublist"},
		{[]int{1, 2, 3}, []int{1, 2, 3, 4, 5}, "sublist"},
		{[]int{3, 4, 5}, []int{1, 2, 3, 4, 5}, "sublist"},
		{[]int{3, 4}, []int{1, 2, 3, 4, 5}, "sublist"},
		{[]int{1, 2, 3}, []int{1, 2, 3}, "equal"},
		{[]int{1, 2, 3, 4, 5}, []int{2, 3, 4}, "superlist"},
		{[]int{1, 2, 4}, []int{1, 2, 3, 4, 5}, "unequal"},
		{[]int{1, 2, 3}, []int{1, 3, 2}, "unequal"},
	}

	for _, tc := range testCases {
		result := findRelation(tc.A, tc.B)
		fmt.Printf("A: %v, B: %v → %s (Expected: %s)\n", tc.A, tc.B, result, tc.expected)
	}
}
