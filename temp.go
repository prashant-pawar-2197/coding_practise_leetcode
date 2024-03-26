package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main(){
sampleData := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("sampleCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &sampleData)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}

}