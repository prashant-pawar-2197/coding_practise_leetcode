package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)
//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\newMultipkg.json
func main(){
	cdr := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("C:/Users/pawarpr/OneDrive - Mavenir Systems, Inc/Documents/GoPractise and notes/newMultipkg.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}

	actionMultiPackageAggregation(cdr, "")
	durationCalculator(cdr)
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("finalPAckageSplitCDR",encodedData, 0644)
}

func durationCalculatorr(cdr map[string]interface{})error{
	if cdr["recordExtensions"] == nil {
		fmt.Println("RecordExtension is nil")
		return nil
    }

    ChargeInformation, ok := cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{})
    if !ok {
    	fmt.Println("CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation parmeter not found in cdr")
    	return nil
    }
	var newChargeInformationArr []interface{}
	for _, val := range ChargeInformation {
		if val != nil {
			newChargeInformationArr = append(newChargeInformationArr, val)
		}
	}
	totalTime := cdr["listOfMultipleUnitUsage"].([]interface{})[0].(map[string]interface{})["usedUnitContainer"].([]interface{})[0].(map[string]interface{})["time"].(float64)
	for i := 0; i < len(newChargeInformationArr); i++ {
		if newChargeInformationArr[i] == nil {
			continue
		}
		if i == len(newChargeInformationArr)-1{
			if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitCash"]; ok{
				newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = totalTime
				break	
			}else if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitAllowance"]; ok{
				newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = totalTime
				break
			}
		}
		if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitCash"]; ok{
			newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = newChargeInformationArr[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"]
			totalTime -= newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"].(float64)
		}else if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitAllowance"]; ok{
			newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = newChargeInformationArr[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"]
			totalTime -= newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"].(float64)
		}
	}
	cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = newChargeInformationArr
	return nil
}

func actionMultiPackageAggregation(cdr map[string]interface{}, transID string) error {

    aggrTaxAmtArrFlg := true
    aggrTaxAmtFlg := true
    aggrBalChangeDebitArrFlg :=  true
    aggrBalChangeDebitFlg := true
    aggrDiscountFlg := true
    aggrDebitAllowanceAmtFlg := true
    aggrRateUsageFlg := true
    debugFlg := true
    addMultiPkgFieldFlg := true

    //There is CDR Format change for M1 & M2 release
    //M1 rate profile is not an array.
    //M2 rate profile is array.
    // Modified to handle release based on Configuration    
	rateProfileIsArray := true

    if cdr["recordExtensions"] == nil {
		fmt.Println("RecordExtension is nil")
		return nil
    }

    iChargeInformation, ok := cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]
    if !ok {
    	fmt.Println("CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation parmeter not found in cdr")
    	return nil
    }

    switch iChargeInformation.(type) {
    case []interface{}:
        rgMap := make(map[string]int)

        i := 0
        arrLen := len(iChargeInformation.([]interface{}))

        for j := 0; j < arrLen ; j++ {
            ChargeInfoMap := iChargeInformation.([]interface{})[i].(map[string]interface{})
            //getRatingGroup(ChargeInfoMap)
            if ChargeInfoMap["ratingIndication"] == nil {
                fmt.Println(". recordExtensions.chargeInformation.ratingIndication parmeter not found in cdr. continue with next chargeInformation")
                continue
            }

            ratingIndication, ok := ChargeInfoMap["ratingIndication"].(map[string]interface{})
            if !ok {
                fmt.Println(". recordExtensions.chargeInformation.ratingIndication parmeter not found in cdr. continue with next chargeInformation")
                continue
            }

            rg, ok := ratingIndication["ratingGroup"].(float64)
            if !ok {
                fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation.ratingIndication.ratingGroup parmeter not found in cdr. continue with next chargeInformation")
            }
			pkg, ok := ratingIndication["packageId"].(string)
			if !ok {
                fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation.ratingIndication.ratingGroup parmeter not found in cdr. continue with next chargeInformation")
        
            }
			prd, ok := ratingIndication["productId"].(string)
			if !ok {
                fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation.ratingIndication.ratingGroup parmeter not found in cdr. continue with next chargeInformation")
           
            }
			rg_pkg_prd := strconv.FormatFloat(rg, 'f', -1, 64)+ pkg+ prd

            if val, ok := rgMap[rg_pkg_prd]; ok {
                //do array level aggregation
                //debit amount.
                //tax amount
                //balance change info - debit amount

                fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". Duplicate ratingGroup found at arrIndex:", i, "Previous index:", val, "RatingGroup:" , rg)

                //Calculation of debit amount
                //firstInstanceOfChargeInfoWithSameRG -> ficiwsrg
                ficiwsrg := iChargeInformation.([]interface{})[val].(map[string]interface{})

                if aggrDebitAllowanceAmtFlg {
                //aggregate debit amount.
                    if ChargeInfoMap["debitAllowance"] != nil {
                        if ficiwsrg["debitAllowance"] == nil {
                            ficiwsrg["debitAllowance"] =  ChargeInfoMap["debitAllowance"]
                        } else {
                            debitAmount2, f2 := ChargeInfoMap["debitAllowance"].(map[string]interface{})["debitAmount"]
                            if f2 {
                                debitAmount1, f1 := ficiwsrg["debitAllowance"].(map[string]interface{})["debitAmount"]
                                if f1  {
                                        if debugFlg {
                                                fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating debitAllowance.debitAmount:")
                                        }
                                    ficiwsrg["debitAllowance"].(map[string]interface{})["debitAmount"] = (debitAmount1.(float64) + debitAmount2.(float64))
                                } else {
                                    if debugFlg {
                                            fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". debit amount not present in first instance of RG. updating debitAllowance.debitAmount")
                                    }
                                    ficiwsrg["debitAllowance"].(map[string]interface{})["debitAmount"] = debitAmount2.(float64)
                                }
                            }
                        }
                    }
                }

                //aggregate tax amount and balance change info
                if ChargeInfoMap["debitCash"] != nil {
                    if ficiwsrg["debitCash"] == nil {
                        ficiwsrg["debitCash"] =  ChargeInfoMap["debitCash"]
                    } else {

                        if rateProfileIsArray {
                            rateProfileArr, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})
                            if f2 {
                                rateProfileArr1, f1 := ficiwsrg["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})
                                if len(rateProfileArr) > 0 {
                                    if aggrTaxAmtArrFlg {
                                        rateProfileArr = aggrTaxAmountArray(rateProfileArr, transID, debugFlg)
                                    }

                                    if aggrTaxAmtFlg {
                                        if f1 && len(rateProfileArr1) > 0 {
                                            if num1, ok := rateProfileArr[0].(map[string]interface{})["taxAmount"].(float64); ok {
                                                if num2, ok := rateProfileArr1[0].(map[string]interface{})["taxAmount"].(float64); ok {
                                                    if debugFlg {
                                                        fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating tax amount:")
                                                    }
                                                    rateProfileArr1[0].(map[string]interface{})["taxAmount"] = num1 + num2
                                                } else {
                                                    if debugFlg {
                                                        fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "tax Amount value not present in first instance of RG. Updating value from current instance")
                                                    }
                                                    rateProfileArr1[0].(map[string]interface{})["taxAmount"] = num1
                                                }
                                            }
                                        } else {
                                            rateProfileArr1 = rateProfileArr
                                        }
                                    }
                                    if aggrTaxAmtFlg || aggrTaxAmtArrFlg {
                                        ficiwsrg["debitCash"].(map[string]interface{})["rateProfile"] = rateProfileArr1
                                    }
                                }
                            }
                        } else {
                            rateProfile, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"]
                            if f2 {
                                rateProfile1, f1 := ficiwsrg["debitCash"].(map[string]interface{})["rateProfile"]
                                if aggrTaxAmtFlg {
                                    if f1 {
                                        if num1, ok := rateProfile.(map[string]interface{})["taxAmount"].(float64); ok {
                                            if num2, ok := rateProfile1.(map[string]interface{})["taxAmount"].(float64); ok {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating tax amount:")
                                                }
                                                rateProfile1.(map[string]interface{})["taxAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "tax Amount value not present in first instance of RG. Updating value from current instance")
                                                }
                                                rateProfile1.(map[string]interface{})["taxAmount"] = num1
                                            }
                                        }
                                    } else {
                                        rateProfile1 = rateProfile
                                    }
                                    ficiwsrg["debitCash"].(map[string]interface{})["rateProfile"] = rateProfile1
                                }
                            }
                        }
                        //Discount Aggregation
                        //=======================================================================================
                        if aggrDiscountFlg {
                            discountProfile, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["discountProfile"]
                            if f2 {
                                discountProfile1, f1 := ficiwsrg["debitCash"].(map[string]interface{})["discountProfile"]
                                 if f1 {
                                        if num1, ok := discountProfile.(map[string]interface{})["discountAmount"].(float64); ok {
                                            if num2, ok := discountProfile1.(map[string]interface{})["discountAmount"].(float64); ok {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating discountProfilediscountAmount:")
                                                }
                                                discountProfile1.(map[string]interface{})["discountAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "discountProfilediscountAmount value not present in first instance of RG. Updating value from current instance")
                                                }
                                                discountProfile1.(map[string]interface{})["discountAmount"] = num1
                                            }
                                        }
                                    } else {
                                        discountProfile1 = discountProfile
                                    }
                                ficiwsrg["debitCash"].(map[string]interface{})["discountProfile"] = discountProfile1
                            }
                        }
                        //=======================================================================================
                        //Aggregate Rated usage :

                        if aggrRateUsageFlg {
                            if num1, ok := ChargeInfoMap["debitCash"].(map[string]interface{})["ratedUsage"].(float64); ok {
                                if num2, ok := ficiwsrg["debitCash"].(map[string]interface{})["ratedUsage"].(float64); ok {
                                    if debugFlg {
                                        fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating ratedUsage")
                                       }
                                    ficiwsrg["debitCash"].(map[string]interface{})["ratedUsage"] = num1 + num2
                                   } else {
                                    if debugFlg {
                                        fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "discountProfilediscountAmount value not present in first instance of RG. Updating value from current instance")
                                       }
                                    ficiwsrg["debitCash"].(map[string]interface{})["ratedUsage"] = num1
                                   }
                              }
                        }

                        //Balance Change info
                        balanceChangeInfoArr, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                        if f2 {
                            balanceChangeInfoArr1, f1 := ficiwsrg["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                            if len(balanceChangeInfoArr) > 0 {
                                if aggrBalChangeDebitArrFlg {
                                    balanceChangeInfoArr = aggrBalChangeDebitAmountArray(balanceChangeInfoArr, transID, debugFlg)
                                }

                                if aggrBalChangeDebitFlg {
                                    if f1 && len(balanceChangeInfoArr1) > 0 {
                                        if num1, ok := balanceChangeInfoArr[0].(map[string]interface{})["debitAmount"].(float64); ok {
                                            if num2, ok := balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"].(float64); ok {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". aggregating balanceChange.debitAmount:")
                                                }
                                                balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "balanceChange.debitAmount value not present in first instance of RG. Updating value from current instance")
                                                }
                                                balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"] = num1
                                            }
                                        }
                                    } else {
                                        balanceChangeInfoArr1 = balanceChangeInfoArr
                                    }
                                }
                                if aggrBalChangeDebitArrFlg || aggrBalChangeDebitFlg{
									//Adding remaining balance to the ficiwsrg
									balanceChangeInfoArr1[0].(map[string]interface{})["newBalance"] = balanceChangeInfoArr[0].(map[string]interface{})["newBalance"].(float64)
                                    ficiwsrg["debitCash"].(map[string]interface{})["balanceChangeInfo"] = balanceChangeInfoArr1
                                }
                            }
                        }
                    }
                }

                if addMultiPkgFieldFlg {
                    ficiwsrg["multiPackageAggregation"] = true
                }

                //Remove Duplicate element
                iChargeInformation.([]interface{})[val] = ficiwsrg
                chargeInfoArr := iChargeInformation.([]interface{})
                copy(chargeInfoArr[i:], chargeInfoArr[i+1:])
                chargeInfoArr[len(chargeInfoArr)-1] = nil
                chargeInfoArr = chargeInfoArr[:len(chargeInfoArr)-1]     // Truncate slice.
                iChargeInformation = chargeInfoArr

            } else {
                rgMap[rg_pkg_prd] = i
                var flag bool
                //aggregate tax amount
                if ChargeInfoMap["debitCash"] != nil {
                        if rateProfileIsArray {
                        rateProfileArr, f1 := ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})
                            if f1 {
                                flag = true
                                if aggrTaxAmtArrFlg {
                                        rateProfileArr = aggrTaxAmountArray(rateProfileArr, transID, debugFlg)
                                        ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"] = rateProfileArr
                                }
                            }
                        }
                        balanceChangeInfoArr, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                        if f2 {
                            if aggrBalChangeDebitArrFlg {
                                    balanceChangeInfoArr = aggrBalChangeDebitAmountArray(balanceChangeInfoArr, transID, debugFlg)
                                    ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"] = balanceChangeInfoArr
                            }
                        }
                        if ( flag || f2) && ( aggrTaxAmtArrFlg || aggrBalChangeDebitArrFlg) {
                            iChargeInformation.([]interface{})[i] = ChargeInfoMap
                        }
                }
                i++
            }
        }
    default:
       fmt.Println("CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ ". recordExtensions.chargeInformation invalid dataType")
        // if exitFromGppFlg {
        //     return errors.New("CustomFunName:" + action.CUSTOMFUNCTION.FunctionName + "recordExtensions.chargeInformation nvalid dataType")
        // } else {
        //     return nil
        // }
        return nil
    }

    if debugFlg {
            fmt.Println( "CustomFunName:", /*action.CUSTOMFUNCTION.FunctionName,*/ "CDR after multi package aggregation. CDR->:", cdr )
    }
    return nil
}

func aggrTaxAmountArrayy( rateProfileList []interface{}, transID string, debugFlg bool) []interface{} {
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

func aggrBalChangeDebitAmountArrayy( balChangeList []interface{}, transID string, debugFlg bool) []interface{} {
    var sum float64
    counter := 0
	var remainingBalance interface{} 
	var flag bool
    if debugFlg {
            fmt.Println( "BalChange debit amount Array Aggregation enabled" , balChangeList)
    }
    if len(balChangeList) > 0 {
        for _, balChange := range balChangeList {
            if val, ok := balChange.(map[string]interface{})["debitAmount"].(float64); ok {
                sum += val
                counter++
            }
			if balChange.(map[string]interface{})["balanceType"] == "GC"{
				remainingBalance, flag = balChange.(map[string]interface{})["newBalance"]
				if !flag {
					continue
				}	
			}
        }
        if counter > 0 {
            balChangeList[0].(map[string]interface{})["debitAmount"] = sum
			if remainingBalance != nil {
				balChangeList[0].(map[string]interface{})["newBalance"] = remainingBalance
			}
            if debugFlg {
                fmt.Println( "After balance Change debit amount arr aggregation. BalChangeLst:" , balChangeList)
            }
        }
    }
    return balChangeList
}
