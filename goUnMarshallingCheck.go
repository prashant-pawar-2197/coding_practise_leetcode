package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SubDocOperations struct {
	OperationType		string		`json:"operationType"`
	SubDocKey			string		`json:"subDocKey"`
	ListOfFields		[]string	`json:"listOfFields"`
	SubDocFieldLength 	string		`json:"subDocFieldLength"`
	Err 				error
	LenOfListOfFields	int
}
type actionList struct {
	ActionList		[]SubDocOperations	`json:"actionList"`
}

func main(){
	var action actionList
	readFileData, err := ioutil.ReadFile("unmarshallingCheck.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	fmt.Println(action)
}