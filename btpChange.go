package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	cdr := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("btpSampleCdr.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	var finalChargeInforArray [2]interface{}
	count := len((cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}))
	for i := 0; i < count; i++ {
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
	}
	for i := 0; i < count; i++ {
		if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] != nil {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
					}
				}
			}else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitAllowance"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
					}
				}
			}

		}
	}

	cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = finalChargeInforArray 
	delete(cdr, "subscriberIdentifier")
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("btpOutputData",encodedData, 0644)
}	

			/*		
					
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["debitCash"]; ok  {
						if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["ratingIndication"]; ok {
							if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes"{
								ratedUsage := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[0].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
								fmt.Println(ratedUsage)
								(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["modifiedDuration"] = ratedUsage 
								(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - ratedUsage
						}
					}
				}
			}	
		}
	}
}

}

}
*/
/*

Way - 1
for i := 0; i < count; i++ {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["debitCash"]; ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["debitCash"]; ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes"{
					ratedUsage := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[0].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
					fmt.Println(ratedUsage)
					(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["modifiedDuration"] = ratedUsage 
					(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - ratedUsage
				}
			}
		}

Way - 2

count := len((cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}))
	if count == 0 {
			return
	}else {
		for i := 0; i < count; i++ {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["debitCash"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No"{
						if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["debitCash"]; ok  {
							if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["ratingIndication"]; ok {
								if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok  && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes"{
									ratedUsage := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[0].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
									fmt.Println(ratedUsage)
									(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[0].(map[string]interface{})["modifiedDuration"] = ratedUsage 
									(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - ratedUsage
							}
						}
					}
				}	
			}
		}
	}
*/