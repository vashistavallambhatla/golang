package main

import (
	"fmt"
	"time"
)

func doWork(done chan bool){
	fmt.Println("doing some work")
	time.Sleep(2 * time.Second)
	fmt.Println("Work done")
	done <- true
}

func main() {
	done := make(chan bool)

	go doWork(done)

	<-done

	fmt.Println("Main completed!")
}