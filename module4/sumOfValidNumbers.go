package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func sumOfValidNumbers() (total int,err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	numbersString := flag.String("Numbers","","a comma separated list of numbers")

	flag.Parse()

	if *numbersString == "" {
		fmt.Println("No numbers provided")
		return  0 , errors.New("no numbers are provided")
	}

	numbers := strings.Split(*numbersString,",")
	
	var sum int

	for _,numStr := range numbers {

		if strings.TrimSpace(numStr) == "" {
			continue
		}

		num , err := strconv.Atoi(numStr)

		if err!=nil {
			return 0, fmt.Errorf("invalid number '%s', unable to convert to integer",numStr)
		}

		sum += num
	}

	return sum,nil
}

func main() {
	total, err := sumOfValidNumbers()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Sum of valid numbers:", total)
	}
}