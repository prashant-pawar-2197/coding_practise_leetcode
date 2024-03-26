package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// func panicHandler(funName string)  {
// 	if err := recover(); err != nil {
// 		fmt.Println("Some panic occurred in function --", funName, " error is ", err)
// 	}
// }

func DeriveRateProfile(rateProfileId string, rateProfileArr []interface{}) []interface{} {
	var rateProfile map[string]interface{}
	if rateProfileId == "" {
		if len(rateProfileArr) == 0 {
			return nil
		}else{
			rateProfile = rateProfileArr[0].(map[string]interface{})
		}
	}else{
		found := false
		for _, v := range rateProfileArr {
			if rateProfileId == v.(map[string]interface{})["rateProfileId"]{
				rateProfile = v.(map[string]interface{})
				found = true
                break
			} 
		}
		if found == false{
			return nil
		}
	}

	switch rateProfile["ratingType"] {
	case "TIERED": for _, v := range rateProfile["rates"].([]interface{}) {
						if rateProfile["lastAccumRatedUsage"].(float64) >= (v.(map[string]interface{})["low"]).(float64) && rateProfile["lastAccumRatedUsage"].(float64) <= (v.(map[string]interface{})["high"]).(float64){
							if _, ok := v.(map[string]interface{})["rateProfileId"] ; !ok || v.(map[string]interface{})["rateProfileId"] == ""{
							var tempRateArr [1]interface{} 
							tempRateArr[0] = v
							rateProfile["rates"] = tempRateArr
							break
							}else{
								return DeriveRateProfile(v.(map[string]interface{})["rateProfileId"].(string), rateProfileArr)
							}
						}else{
							continue	
						} 
				}
	case "TELESCOPIC":	var tempRateArr [1]interface{}
						tempRateArr[0] = rateProfile["rates"].([]interface{})[0]
						rateProfile["rates"] = tempRateArr
	}
	var tempRateProfileArr []interface{}
	tempRateProfileArr = append(tempRateProfileArr, rateProfile)
	return tempRateProfileArr
}

func main() {
	sampleData := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("sampleCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &sampleData)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	defer panicHandler("deriveRateProfile")

	if _, ok := sampleData["recordExtensions"].(map[string]interface{})["chargeInformation"]; !ok {
		fmt.Println("ChargeInformation array is absent in the cdr")
		return	
	}
	
	for _, v := range (sampleData["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}) {
		var taxAmt float64
		if _, ok := v.(map[string]interface{})["debitCash"]; ok {
			if _, ok := v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"]; ok {
				if _, ok := v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["taxAmount"]; ok {
					taxAmt  = (v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}))[0].(map[string]interface{})["taxAmount"].(float64) 
				}else{
					fmt.Println("taxAmount is not present in the first RateProfile")
				}
				v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"] = DeriveRateProfile("",v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}))
				(v.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}))[0].(map[string]interface{})["taxAmount"] = taxAmt
			}else{
				continue
			}
		}else{
			continue
		}
	}
	encodedData,_ := json.Marshal(sampleData)

	ioutil.WriteFile("RATEPROFILEPROCESSEDCDR",encodedData, 0644)
}

































//fmt.Println((sampleData["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"])
//fmt.Println(reflect.TypeOf(rateArr.([]interface{})[0].(map[string]interface{})["rates"]))
	//fmt.Println((sampleData["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"])

/*
var finalRateObj Rate
	var rateProfileArr []RateProfile
	for _, v := range (sampleData["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}) {
		rateProfileData, err  := json.Marshal(v)
		if err != nil{
			fmt.Println("error occured while unmarshalling")
		}
		var rateProfileInstances RateProfile
		err = json.Unmarshal([]byte(rateProfileData), &rateProfileInstances)
		rateProfileArr = append(rateProfileArr, rateProfileInstances)
	}
	if rateProfileArr[0].RatingType == "TELESCOPIC" {
		finalRateObj = rateProfileArr[0].Rates[0]
	}else if rateProfileArr[0].RatingType == "TIERED"  {
		
	}
	fmt.Println(finalRateObj)
*/