package main

import (
	"fmt"
	"os"
)

func fileReader() {
	fmt.Println("Please enter the file name to be read")
	var fileName string
	fmt.Scanf("%s",&fileName)

	file,err := os.Open(fileName)

	if err!=nil {
		fmt.Println("Error: ",err)
		return 
	}

	defer file.Close()

	buffer := make([]byte,10)
	for {
		n,err := file.Read(buffer)
		if err!=nil {
			break
		} else {
			fmt.Print(string(buffer[:n]))
		}
	}
}