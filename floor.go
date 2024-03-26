package main

import (
	"fmt"
	//"io/ioutil"
)
func mainn()  {
	// sampleData := make(map[string]interface{})
	// readFileData, err := ioutil.ReadFile("sampleCDR.json")
	// if err != nil {
	// 	fmt.Println("Error occured")
	// }
	// err = json.Unmarshal(readFileData, &sampleData)
	// if err != nil {
	// 	fmt.Println("Error occured while unmarshalling", err)	
	// }
	// if sampleData["recordType"] == nil {
	// 	fmt.Println("recordType field is absent")
	// }else if sampleData["recordType"] == ""{
	// 	fmt.Println("recordType field is present but empty string")
	// }else{
	// 	fmt.Println("recordType field is present")
	// }
	
	var finalChargeInforArray [2]interface{}

	if finalChargeInforArray[0] == nil || finalChargeInforArray[1] == nil{
		fmt.Println("array is empty")
	}
	
}
