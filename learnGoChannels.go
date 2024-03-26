package main

import (
	"fmt"
	"sync"
	"time"
)

func sendNum(ch chan int64, wg *sync.WaitGroup){
	defer wg.Done()
	counter := 0
	time.Sleep(time.Millisecond)
	for i := 1; i < 11; i++ {
		counter++
		ch <- int64(counter)
		fmt.Println(counter)
	}
}

func main(){
	nameChan := make(chan int64)
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(1)
	go sendNum(nameChan, wg)
	fmt.Println("Came out of go routine")
	for i := 0; i < 10; i++ {
		fmt.Println(<- nameChan)
	}
	wg.Wait()
}