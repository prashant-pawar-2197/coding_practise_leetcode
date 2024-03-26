package main

import (
	"fmt"
	"sync"
//	"time"
)

func (rs *RandomStruct) printA(wg *sync.WaitGroup){
	rs.mu.Lock()
	defer wg.Done()
	for i := 0; i < 10; i++ {
		fmt.Println("A")	
	}
	rs.mu.Unlock()
}

func (rs *RandomStruct) printB(wg *sync.WaitGroup){
	//time.Sleep(time.Microsecond)
	rs.mu.Lock()
	defer wg.Done()
	
	for i := 0; i < 10; i++ {
		fmt.Println("B")	
	}
	rs.mu.Unlock()
}

type RandomStruct struct {
	mu sync.Mutex
}

func main(){
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	var rs *RandomStruct = &RandomStruct{}
	wg.Add(2)
	go rs.printA(wg)
	go rs.printB(wg)
	wg.Wait()
}