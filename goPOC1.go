package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main(){
	// actualCounter := 0
	// const COUNTER int = 5

	// for ; actualCounter < COUNTER ; actualCounter++ {
	// 	fmt.Println("Value of actualCounter -->",actualCounter)
	// 	if actualCounter == 3{
	// 		fmt.Println("Breaking from actualCounter at value -->",actualCounter)
	// 		break
	// 	}
	// }
	// for ; actualCounter < COUNTER ; actualCounter++{
	// 	fmt.Println("Resuming for loop with actualCounter at value -->",actualCounter)
	// 	fmt.Println(actualCounter)
	// }

	cdr := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("C:/Users/pawarpr/OneDrive - Mavenir Systems, Inc/Documents/GoPractise and notes/packageSplitSamples/packageSplit3(PackageChange).json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	val , ok := cdr["recordType"].(string)
	if !ok{
		fmt.Println(ok)
	}
	fmt.Println(val)
}