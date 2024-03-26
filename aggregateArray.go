package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type AggregateArray struct {
	ListOfArray []struct {
        Source string `json:"source"`
        Destination string `json:"destination"`
    }
}

func actionAggregateArray(cdr map[string]interface{}, action AggregateArray) error {
	var (
		counter int = 0
		sum int = 0
	)
	for _, arr := range action.ListOfArray {
		for {
			cdrFieldName := strings.Replace(arr.Source, "$i", strconv.Itoa(counter), -1)
			counter++
			if sourceFieldVal, ok := cdr[cdrFieldName].(float64); !ok {
				fmt.Println("Source field ", cdrFieldName, " is not present")
				if counter > 0 && sum != 0 {
					cdr[arr.Destination] = sum
				}
				sum = 0
				counter = 0
				break
			} else {
				sum = sum + int(sourceFieldVal)
			}
		}
	}
	return nil
}

func main() {
	var action AggregateArray
	var cdr map[string]interface{}
	readFileData, err := ioutil.ReadFile("sampleActionAggregateArray.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	//fmt.Println(action)
	readCDRData, err := ioutil.ReadFile("pgwCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	err = json.Unmarshal(readCDRData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	actionAggregateArray(cdr, action)
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("pgwOutputCdr",encodedData, 0644)
}
