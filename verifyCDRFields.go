package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	sampleData := make(map[string]interface{})
	actualData := make(map[string]interface{})

	readFileData, err := ioutil.ReadFile("BSSRequiredFields.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &sampleData)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}

	
	readFileData, err = ioutil.ReadFile("actualCDRoutput.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &actualData)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}

	for _, v := range sampleData["Body"].([]interface{}) {
		for _, y := range v.(map[string]interface{})["Parameters"].([]interface{}) {
			if _, ok := actualData[y.(map[string]interface{})["InternalParameter"].(string)]; !ok{
				fmt.Println(y.(map[string]interface{})["InternalParameter"].(string), " is absent in final cdr")
			}
		}
		for _, x := range v.(map[string]interface{})["Parameters"].([]interface{}) {
			if _, ok := actualData[x.(map[string]interface{})["InternalParameter"].(string)]; ok{
				if actualData[x.(map[string]interface{})["InternalParameter"].(string)] == ""{
					fmt.Println(x.(map[string]interface{})["InternalParameter"].(string), " is present in final cdr but is empty")
				}
			}
		}
	}
	
}