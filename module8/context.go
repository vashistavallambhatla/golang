package main

import (
	"fmt"
	"context"
	"time"
)

func exampleTimeout() {
	ctx := context.Background()
	ctxWithTimeOut, cancel := context.WithTimeout(ctx,4 * time.Second)
	defer cancel()

	done := make(chan bool)

	go func(){
		time.Sleep(6 * time.Second)
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("API CALLED")
	case <-ctxWithTimeOut.Done():
		fmt.Println("Oh no timeout expired!",ctxWithTimeOut.Err())
	}
}