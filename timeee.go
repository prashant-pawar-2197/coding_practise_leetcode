package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

)

type RestartTimerAfterCoolingPeriod struct {
	CoolingPeriod		int			`json:"coolingPeriod"`
	Timeout				int			`json:"timeout"`
	Err 				error
}

func main() {
	var action RestartTimerAfterCoolingPeriod
	readFileData, err := ioutil.ReadFile("sampleStruct.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	fmt.Println(reflect.TypeOf(action.CoolingPeriod))
}