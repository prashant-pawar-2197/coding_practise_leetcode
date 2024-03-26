package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	// "io/ioutil"
	// "reflect"
	"time"
	// "strconv"
	// "strings"
)

type RgPkgInfo struct {
	Package		string
	Product		string
}

//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\newMultipkg.json
func main(){
	// cdr := make(map[string]interface{})
	// //var action CustomFunction
	// readFileData, err := ioutil.ReadFile("sdr.json")
	// if err != nil {
	// 	fmt.Println("Error occured")
	// }
	// err = json.Unmarshal(readFileData, &cdr)
	// if err != nil {
	// 	fmt.Println("Error occured while unmarshalling", err)	
	// }
	 
	// 	fmt.Println(reflect.TypeOf( cdr["DELIVERY_TIME"]).Kind())
	// 	if fieldValue, ok := cdr["DELIVERY_TIME"].(float64); ok {
	// 		fmt.Println(fieldValue)
	// 		fmt.Println(int64(fieldValue))
	// 		UtcTime := time.Unix(int64(fieldValue), 0)
	// 		location, err := time.LoadLocation("Europe/Berlin")
	// 		if err != nil {
	// 			fmt.Println("Error occured while determining the location")
	// 		}
	// 		UtcTime = UtcTime.In(location)
	// 		fmt.Println(UtcTime)
	// 	} 
		
	UtcTime1 := time.UnixMilli(1684288217000)
	UtcTime2 := time.Unix(1684288217, 0)
	fmt.Println(UtcTime1)
	fmt.Println(UtcTime2)
	// encodedData,_ := json.Marshal(cdr)
	// ioutil.WriteFile("finalOutputPackageSplitCDR",encodedData, 0644)
}

