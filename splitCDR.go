package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)


type SplitOnCdrFieldChange struct{
	ParamList []ParamDetail	`json:"paramList"`
	SctName	string 			`json:"sctName"`
	YesNoFields []string 	`json:"yesNoFields"`
} 

type ParamDetail struct {
    FieldName string	`json:"fieldName"`
    IsMandatary bool	`json:"isMandatory"`
}

func main(){
	var action SplitOnCdrFieldChange
	fmt.Println(action.ParamList)
	var cdr map[string]interface{}
	readFileData, err := ioutil.ReadFile("splitCDRAction.json")
	fmt.Println(string(readFileData))
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	fmt.Println(action)
	readFileData2, err := ioutil.ReadFile("sampleCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData2, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}

		var splitKeyToInsert = make(map[string]interface{})

		action.ParamList = nil  
		var NameArr []string
		for _, v := range NameArr{
			fmt.Println("Entered",v)
		}
		if splitKeyToInsert == nil {
			fmt.Println("Empty")
		}
}