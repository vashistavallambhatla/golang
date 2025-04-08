package main

import (
	"fmt"
	"sync"
)

func mockDBCall(i int,wg *sync.WaitGroup,res *int,mut *sync.Mutex) {
	defer wg.Done()
	mut.Lock()
	*res += i
	fmt.Printf("%v\n",i)
	mut.Unlock()
}

func mainn() {
	wg := sync.WaitGroup{}
	mut := sync.Mutex{}
	var res int

	wg.Add(5)
	for i := 0; i < 5 ; i++ {
		go mockDBCall(i,&wg,&res,&mut)
	}

	wg.Wait()
	fmt.Printf("result %v\n",res)
	fmt.Println("Successfully completed")
}