package main

import (
	"fmt"
	"net/http"
	"time"
	"context"
)

func LongRunningTask(done chan bool){
	fmt.Println("Process started running..")
	time.Sleep(5 * time.Second)
	fmt.Println("Process completed")
	done <- true
}

func ProcessHandler(w http.ResponseWriter, r * http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(),7 * time.Second)
	defer cancel()
	
	done := make(chan bool) 
	defer close(done)
	
	go LongRunningTask(done)

	select { 
		case <- done:
			fmt.Fprintf(w,"Process completed successfully")
		case <- ctx.Done():
			if ctx.Err() == context.Canceled {
				http.Error(w, "Request was cancelled by the client", http.StatusRequestTimeout)
				fmt.Println("Request was cancelled by the client.")
			} else if ctx.Err() == context.DeadlineExceeded{
				http.Error(w,"Requqest timed out",http.StatusRequestTimeout)
				fmt.Println("Request timed out")
			}
	}
}  

func main() {
	http.HandleFunc("/process",ProcessHandler)

	fmt.Println("Starting a server on port : 8080...") 

	if err := http.ListenAndServe(":8080",nil); err != nil {
		fmt.Println("Error occured while starting the server :",err)
	}
}