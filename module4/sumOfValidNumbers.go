package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func sumOfValidNumbers() (total int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	numbersString := flag.String("Numbers", "", "a comma separated list of numbers")
	flag.Parse()

	if *numbersString == "" {
		return 0, errors.New("no numbers provided")
	}

	numbers := strings.Split(*numbersString, ",")

	for _, numStr := range numbers {
		if strings.TrimSpace(numStr) == "" {
			continue
		}

		num, convErr := strconv.Atoi(numStr)
		if convErr != nil {
			return 0, fmt.Errorf("invalid number '%s', unable to convert to integer", numStr)
		}

		total += num
	}

	return total, nil
}

func main() {
	total, err := sumOfValidNumbers()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sum of valid numbers:", total)
	}
}