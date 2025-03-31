package main

import (
	"fmt"
	"net/http"
	"time"
	"context"
)

func ProcessHandler(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintf(w,"Starting the process")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w,"Process completed")
}  

func main() {
	ctx := context.Background()

	ctxWithTimeOut, cancel := context.WithTimeout(ctx,5 * time.Second)
	defer cancel()

	http.HandleFunc("/process",ProcessHandler)

	fmt.Println("Starting a server on port : 8080...") 

	if err := http.ListenAndServe(":8080",nil); err != nil {
		fmt.Println("Error occured while starting the server :",err)
	}

	select {
		case <- ctxWithTimeOut.Done() :
			fmt.Printf("The process ran for too long forcing main to quit")
	}
}