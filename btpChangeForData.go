package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func panicHandler(funcName, transID string)  {
	if err := recover(); err != nil {
		fmt.Println(transID , "An Error occured in function --", funcName ," error is --", err)
	}
}

//Function to get rating group
func getRatingGroup(chargeInformation map[string]interface{}, transID string) (float64) {
    defer panicHandler("actionCustomFunction", transID)
 	ratingGroup, _ := chargeInformation["ratingIndication"].(map[string]interface{})["ratingGroup"].(float64)	
	return ratingGroup
}

func getLocalSeqNum(chargeInformation map[string]interface{}, transID string) (float64) {
    defer panicHandler("actionCustomFunction", transID)
 	ratingGroup, _ := chargeInformation["ratingIndication"].(map[string]interface{})["localSequenceNumber"].(float64)	
	return ratingGroup
}

func aggregateTwoChargeInfos(prevChargeInfo , currentChargeInfo map[string]interface{}, transID string, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg bool) {
	if currDebitCash, ok := currentChargeInfo["debitCash"].(map[string]interface{}); ok {
		if prevDebitCash, ok := prevChargeInfo["debitCash"].(map[string]interface{}); ok {
			// both are debitCash aggregate them
			aggregateTwoDebitCash(prevDebitCash, currDebitCash, transID, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg)
			return
		} else if _, ok := prevChargeInfo["debitAllowance"]; ok {
			prevChargeInfo["debitCash"] = currDebitCash
			return
		}
	} else if currDebitAllowance, ok := currentChargeInfo["debitAllowance"].(map[string]interface{}); ok {
		if prevDebitAllowance, ok := prevChargeInfo["debitAllowance"].(map[string]interface{}); ok {
			// both are debitCash aggregate them
			aggregateTwoDebitAllowances(prevDebitAllowance, currDebitAllowance)
			return
		} else if _, ok := prevChargeInfo["debitCash"]; ok {
			prevChargeInfo["debitAllowance"] = currDebitAllowance
			return
		}
	}
}

func aggregateTwoDebitCash(prevDebitCash, currDebitCash map[string]interface{}, transID string, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg bool) {
	debugFlg := true // <------ TBR

	// aggregate TaxAmount 
	rateProfileArr, f2 := currDebitCash["rateProfile"].([]interface{})
    if f2 {
        rateProfileArr1, f1 := prevDebitCash["rateProfile"].([]interface{})
        if len(rateProfileArr) > 0 {
            if aggrTaxAmtArrFlg {
                rateProfileArr = aggrTaxAmountArray(rateProfileArr, transID, debugFlg)
				rateProfileArr1 = aggrTaxAmountArray(rateProfileArr1, transID, debugFlg)
            }

            if aggrTaxAmtFlg {
                if f1 && len(rateProfileArr1) > 0 {
                    if num1, ok := rateProfileArr[0].(map[string]interface{})["taxAmount"].(float64); ok {
                        if num2, ok := rateProfileArr1[0].(map[string]interface{})["taxAmount"].(float64); ok {
                            rateProfileArr1[0].(map[string]interface{})["taxAmount"] = num1 + num2
                        } else {
                            rateProfileArr1[0].(map[string]interface{})["taxAmount"] = num1
                        }
                    }
                } else {
                    rateProfileArr1 = rateProfileArr
                }
            }
            if aggrTaxAmtFlg || aggrTaxAmtArrFlg {
                prevDebitCash["rateProfile"] = rateProfileArr1
            }
        }
    }

	//Discount Aggregation
    if aggrDiscountFlg {
        discountProfile, f2 := currDebitCash["discountProfile"]
        if f2 {
            discountProfile1, f1 := prevDebitCash["discountProfile"]
             if f1 {
                    if num1, ok := discountProfile.(map[string]interface{})["discountAmount"].(float64); ok {
                        if num2, ok := discountProfile1.(map[string]interface{})["discountAmount"].(float64); ok {
                            discountProfile1.(map[string]interface{})["discountAmount"] = num1 + num2
                        } else {
                            discountProfile1.(map[string]interface{})["discountAmount"] = num1
                        }
                    }
                } else {
                    discountProfile1 = discountProfile
                }
				prevDebitCash["discountProfile"] = discountProfile1
        }
    }

    //Aggregate Rated usage :
    if aggrRateUsageFlg {
        if num1, ok := currDebitCash["ratedUsage"].(float64); ok {
            if num2, ok := prevDebitCash["ratedUsage"].(float64); ok {
                prevDebitCash["ratedUsage"] = num1 + num2
            } else {
            	prevDebitCash["ratedUsage"] = num1
            }
        }
    }

	// remainingBal
	var (
		prevBalChangeInfoArr []interface{}
		currBalChangeInfoArr []interface{}
		prevBalChangeInfoArrPresent	 bool
		currBalChangeInfoArrPresent	 bool
	)
	prevBalChangeInfoArr, prevBalChangeInfoArrPresent = prevDebitCash["balanceChangeInfo"].([]interface{})
	currBalChangeInfoArr, currBalChangeInfoArrPresent = currDebitCash["balanceChangeInfo"].([]interface{})

	if !prevBalChangeInfoArrPresent {
		if currBalChangeInfoArrPresent {
			prevDebitCash["balanceChangeInfo"] = currBalChangeInfoArr
		}
	} else {
		if currBalChangeInfoArrPresent {
			for i := 0; i < len(currBalChangeInfoArr); i++ {
				if balChangeInfo, ok := currBalChangeInfoArr[i].(map[string]interface{}); ok {
					if balChangeInfo["balanceType"] == "GC" {
						newBalSet := false
						if currRemainBal, ok := balChangeInfo["newBalance"].(float64); ok {
							for j := 0; j < len(prevBalChangeInfoArr); j++ {
								if balChangeInfoMap, ok := prevBalChangeInfoArr[j].(map[string]interface{}); ok {
									if balChangeInfoMap["balanceType"] == "GC" {
										//Reset newBalance
										if prevRemainBal, ok := balChangeInfoMap["newBalance"].(float64); ok && currRemainBal < prevRemainBal {
											balChangeInfoMap["newBalance"] = currRemainBal
											newBalSet = true
											break
										}
									}
								}
							}
						}
						if newBalSet {
							break
						}
					}
				}
			}
			// aggregating balChangeInfo's debitAmount within the balChangeInfo arr
			if aggrBalChangeDebitArrFlg {
				prevBalChangeInfoArr = aggrBalChangeDebitAmountArray(prevBalChangeInfoArr, transID, debugFlg)
				currBalChangeInfoArr = aggrBalChangeDebitAmountArray(currBalChangeInfoArr, transID, debugFlg)
			}
			// aggregating balChangeInfo's debitAmount
			if aggrBalChangeDebitFlg {
				if num1, ok := currBalChangeInfoArr[0].(map[string]interface{})["debitAmount"].(float64); ok {
					if num2, ok := prevBalChangeInfoArr[0].(map[string]interface{})["debitAmount"].(float64); ok {
						if debugFlg {
							fmt.Println("CustomFunName:", ". aggregating balanceChange.debitAmount:")
						}
						prevBalChangeInfoArr[0].(map[string]interface{})["debitAmount"] = num1 + num2
					} else {
						if debugFlg {
							fmt.Println("CustomFunName:", "balanceChange.debitAmount value not present in first instance of RG. Updating value from current instance")
						}
						prevBalChangeInfoArr[0].(map[string]interface{})["debitAmount"] = num1
					}
				}			
			}
			if aggrBalChangeDebitArrFlg || aggrBalChangeDebitFlg{
				prevDebitCash["balanceChangeInfo"] = prevBalChangeInfoArr
			}
		}
	}
}

func aggrBalChangeDebitAmountArray( balChangeList []interface{}, transID string, debugFlg bool) []interface{} {
    var sum float64
    counter := 0
    if debugFlg {
            fmt.Println("BalChange debit amount Array Aggregation enabled" , balChangeList)
    }
    if len(balChangeList) > 0 {
        for _, balChange := range balChangeList {
            if val, ok := balChange.(map[string]interface{})["debitAmount"].(float64); ok {
                sum += val
                counter++
            }
        }
        if counter > 0 {
            balChangeList[0].(map[string]interface{})["debitAmount"] = sum
            if debugFlg {
                fmt.Println("After balance Change debit amount arr aggregation. BalChangeLst:" , balChangeList)
            }
        }
    }
    return balChangeList
}

func aggrTaxAmountArray( rateProfileList []interface{}, transID string, debugFlg bool) []interface{} {
    var sum float64
    counter := 0

    if debugFlg {
            fmt.Println( "TaxAmountArray Aggregation enabled" , rateProfileList)
    }

    if len(rateProfileList) > 0 {
        for _, rateProfile := range rateProfileList {
            if val, ok := rateProfile.(map[string]interface{})["taxAmount"].(float64); ok {
                sum += val
                counter++
            }
        }

        if counter > 0 {
            rateProfileList[0].(map[string]interface{})["taxAmount"] = sum
            if debugFlg {
                fmt.Println( "After tax amount arr aggregation. rateProfile:" , rateProfileList)
            }
        }
    }
    return rateProfileList
}

func aggregateTwoDebitAllowances(prevDebitAllowance, currDebitAllowance map[string]interface{}){
	//aggregate debit amount.
	debitAmount2, f2 := currDebitAllowance["debitAmount"]
	if f2 {
		debitAmount1, f1 := prevDebitAllowance["debitAmount"]
		if f1 {
			prevDebitAllowance["debitAmount"] = (debitAmount1.(float64) + debitAmount2.(float64))
		} else {
			prevDebitAllowance["debitAmount"] = debitAmount2.(float64)
		}
	}
}

type DataRgWiseInfo struct {
	TotalVol		float64
	UplinkVolume	float64
	DownlinkVolume	float64
}

func storeRgAndLocalSeqWiseVolumeInfo(cdr map[string]interface{}, transID string) (dataRgMap map[string]*DataRgWiseInfo) {
	dataRgMap = make(map[string]*DataRgWiseInfo)
	if listOfMultipleUnitUsage, ok := cdr["listOfMultipleUnitUsage"].([]interface{}); ok {
		for _, v := range listOfMultipleUnitUsage {
			var (
				ratingGroup 	float64
				ratingGroupStr 	string
				localSeqNum 	float64
				localSeqNumStr	string
			)
			if elem, ok := v.(map[string]interface{}); ok {
				if ratingGroup, ok = elem["ratingGroup"].(float64); !ok {
					continue
				} else {
					ratingGroupStr = strconv.FormatFloat(ratingGroup, 'f', -1, 64)
				}
				if usedUnitContainer, ok := elem["usedUnitContainer"].([]interface{}); ok {
					for _, val := range usedUnitContainer {
						if element, ok := val.(map[string]interface{}); ok {
							if localSeqNum, ok = element["localSequenceNumber"].(float64); !ok {
								continue
							} else {
								localSeqNumStr = strconv.FormatFloat(localSeqNum, 'f', -1, 64)
							}
							var dataRgWiseInfo *DataRgWiseInfo = &DataRgWiseInfo{}
							dataRgWiseInfo.TotalVol, _ = element["totalVolume"].(float64)
							dataRgWiseInfo.UplinkVolume, _ = element["uplinkVolume"].(float64)
							dataRgWiseInfo.DownlinkVolume, _ = element["downlinkVolume"].(float64)
							
							dataInfokey := ratingGroupStr + "_" + localSeqNumStr

							// store in dataRgMap
							dataRgMap[dataInfokey] = dataRgWiseInfo
						}
					}
				}
			}
		}	
	}
	return dataRgMap
}

func copyVolumeInfoIntoChargeInfo(chargeInfo map[string]interface{}, dataRgMap map[string]*DataRgWiseInfo, transID string){
	rg := getRatingGroup(chargeInfo, transID)
	if rg != 0.0 {
		rgStr := strconv.FormatFloat(rg, 'f', -1, 64)
		localSeqNum := getLocalSeqNum(chargeInfo, transID)
		if localSeqNum != 0.0 {
			localSeqNumStr := strconv.FormatFloat(localSeqNum, 'f', -1, 64)
			key := rgStr + "_" + localSeqNumStr
			if dataRgWiseInfo, ok := dataRgMap[key]; ok {
				chargeInfo["totalVolume"] = dataRgWiseInfo.TotalVol
				chargeInfo["uplinkVolume"] = dataRgWiseInfo.UplinkVolume
				chargeInfo["downlink"] = dataRgWiseInfo.DownlinkVolume
			} else {
				fmt.Println(key, "No corresponding data stats found for the combo of rg + localSeqNum")
			}
		}
	}
}

func aggregateMultipleChargeInformationForBtpChangeData(cdr map[string]interface{}, transID string, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg bool) {
	defer panicHandler("aggregateMultipleChargeInformationForBtpChange", transID)
	// get RG+LocalSeqNum wise voulme details to be populated in chargeInformation further
	dataRgMap := storeRgAndLocalSeqWiseVolumeInfo(cdr,transID)

	// Creating an array which will store chargeInformation with btpChangeFlag No and  chargeInformation with btpChangeFlag Yes
	// Once both chargeInformations are determined then aggregate all chargeinformations with btpChangeFlag as No
	finalChargeInforArray := make([]interface{}, 0)
	chargeInfoRgWise := make(map[float64]map[string]interface{})
	count := len((cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}))
	for i := 0; i < count; i++ {
		chargeInfo, _ := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})
		if _, ok := chargeInfo["ratingIndication"]; ok {
			if val, ok := chargeInfo["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && val == "No" {
				// get RG from chargeInfo
				rg := getRatingGroup(chargeInfo, transID)
				if rg != 0.0 {
					if btpChargeInfoMap, ok := chargeInfoRgWise[rg]; ok {
						// this means for the given RG there already exists some btpChargeInformations
						// so we will aggregate the prev and current chargeInfo
						if _, ok := btpChargeInfoMap["No"]; !ok {
							copyVolumeInfoIntoChargeInfo(chargeInfo, dataRgMap, transID)
							btpChargeInfoMap["No"] = chargeInfo
							chargeInfoRgWise[rg] = btpChargeInfoMap
						} else {
							prevChargeInfo, _ := btpChargeInfoMap["No"].(map[string]interface{})
							aggregateTwoChargeInfos(prevChargeInfo, chargeInfo, transID, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg)
							btpChargeInfoMap["No"] = prevChargeInfo
						}
					} else {
						copyVolumeInfoIntoChargeInfo(chargeInfo, dataRgMap, transID)
						btpChargeInfoMap := make(map[string]interface{})
						btpChargeInfoMap["No"] = chargeInfo
						chargeInfoRgWise[rg] = btpChargeInfoMap
					}
				} else {
					continue
				}
			} else if val, ok := chargeInfo["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && val == "Yes" {
				rg := getRatingGroup(chargeInfo, transID)
				if rg != 0.0 {
					if btpChargeInfoMap, ok := chargeInfoRgWise[rg]; ok {
						// this means for the given RG there already exists some btpChargeInformations
						// so we will aggregate the prev and current chargeInfo
						if _, ok := btpChargeInfoMap["Yes"]; !ok {
							copyVolumeInfoIntoChargeInfo(chargeInfo, dataRgMap, transID)
							btpChargeInfoMap["Yes"] = chargeInfo
							chargeInfoRgWise[rg] = btpChargeInfoMap
						} else {
							prevChargeInfo, _ := btpChargeInfoMap["Yes"].(map[string]interface{})
							aggregateTwoChargeInfos(prevChargeInfo, chargeInfo, transID, aggrTaxAmtArrFlg, aggrTaxAmtFlg, aggrDiscountFlg, aggrRateUsageFlg, aggrBalChangeDebitArrFlg, aggrBalChangeDebitFlg)
							btpChargeInfoMap["Yes"] = prevChargeInfo
						}
					} else {
						copyVolumeInfoIntoChargeInfo(chargeInfo, dataRgMap, transID)
						btpChargeInfoMap := make(map[string]interface{})
						btpChargeInfoMap["Yes"] = chargeInfo
						chargeInfoRgWise[rg] = btpChargeInfoMap
					}
				} else {

					break
				}
			}
		}
	}
	for _, v := range chargeInfoRgWise {
		if chargeInfo, ok := v["No"]; ok {
			finalChargeInforArray = append(finalChargeInforArray, chargeInfo)
		}
		if chargeInfo, ok := v["Yes"]; ok {
			finalChargeInforArray = append(finalChargeInforArray, chargeInfo)
		}
	}
	if len(finalChargeInforArray) > 0 {
		cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = finalChargeInforArray
	}
}

func main(){
	cdr := make(map[string]interface{})
	//var action SplitCdr
	readFileData, err := ioutil.ReadFile("btpSampleData5FUP.json")
	if err != nil {
		fmt.Println("Error occured")
	}

	// readActionData, err := ioutil.ReadFile("multipackageAction.json")
	// if err != nil {
	// 	fmt.Println("Error occured")
	// }
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	// err = json.Unmarshal(readActionData, &action)
	// if err != nil {
	// 	fmt.Println("Error occured while unmarshalling", err)
	// }
	aggregateMultipleChargeInformationForBtpChangeData(cdr, "transid", false, true, true, true, true, true)
	encodedData, _ := json.Marshal(cdr)
	ioutil.WriteFile("btpOutput", encodedData, 0644)
}