package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// type CdrType struct{
// 	rg float64
// 	pkg string
// 	product string
// 	rateProfile string
// 	rgFlag bool
// 	pkgFlag bool
// 	prdFlag bool
// 	rateProFlag bool
// }

func main() {

	cdr := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("packageSplit.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	//var finalChargeInforArray [2]interface{}
	var cdrTypeDeduced CdrType
	var cdrCategory string
	count := len((cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}))
	for i := 0; i < count; i++ {
		if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] != nil {
			if (CdrType{}) != cdrTypeDeduced{
				rg := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["ratingGroup"].(float64)
				pkg := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["packageId"].(string)
				prd := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["productId"].(string)
				var rp string
				if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"] != nil {
					rp = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["rateProfileId"].(string)
				}
				if cdrTypeDeduced.rg != rg && !cdrTypeDeduced.rgFlag{
					cdrTypeDeduced.rgFlag = true
				}
				if cdrTypeDeduced.pkg != pkg && !cdrTypeDeduced.pkgFlag{
					cdrTypeDeduced.pkgFlag = true
				}
				if cdrTypeDeduced.product != prd && !cdrTypeDeduced.prdFlag{
					cdrTypeDeduced.prdFlag = true
				}
				if cdrTypeDeduced.rateProfile != rp && !cdrTypeDeduced.rateProFlag{
					cdrTypeDeduced.rateProFlag = true
				}
			}else{
				if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"] != nil{
					cdrTypeDeduced.rg = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["ratingGroup"].(float64)
					cdrTypeDeduced.pkg = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["packageId"].(string)
					cdrTypeDeduced.product = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["productId"].(string)
				}
				if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"] != nil{
					cdrTypeDeduced.rateProfile = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["rateProfileId"].(string)
				}
			}
		}
	}
	// false_true_true_false
	
	for i := 0; i < 4; i++ {
		if cdrTypeDeduced.rgFlag && i == 0 {
			cdrCategory += "RGNS_"
			continue
		} 
		if !cdrTypeDeduced.rgFlag && i == 0  {
			cdrCategory += "RGS_"
			continue
		}
		if cdrTypeDeduced.pkgFlag && i == 1 {
			cdrCategory += "PKGNS_"
			continue
		}
		if !cdrTypeDeduced.pkgFlag && i == 1 {
			cdrCategory += "PKGS_"
			continue
		}
		if cdrTypeDeduced.prdFlag && i == 2 {
			cdrCategory += "PRDNS_"
			continue
		}
		if !cdrTypeDeduced.prdFlag  && i == 2{
			cdrCategory += "PRDS_"
			continue
		}
		if cdrTypeDeduced.rateProFlag && i == 3{
			cdrCategory += "RPNS_"
			continue
		}
		if !cdrTypeDeduced.rateProFlag && i == 3{
			cdrCategory += "RPS_"
			continue
		}
	}
	fmt.Println(cdrCategory)
	// cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = finalChargeInforArray 
	// delete(cdr, "subscriberIdentifier")
	// encodedData,_ := json.Marshal(cdr)
	// ioutil.WriteFile("btpOutputData",encodedData, 0644)
}	

	
/*
	if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"]; ok {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
					if finalChargeInforArray[0] == nil{ 
						finalChargeInforArray[0] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						continue
					}else{
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						continue
					}
				} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes"{
					finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
					finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
					(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] =  nil
					continue
				}
			}
		} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitAllowance"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
						if finalChargeInforArray[0] == nil{
							finalChargeInforArray[0] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
							finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
							(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
							continue
						}else {
							finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
							(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
							continue
						}	
					} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes"{
						if i == count -1 {
							finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]						
							finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
							(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] =  nil
							break
						}
						finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]						
						finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] =  nil
						break
					}
				}

		}
		*/