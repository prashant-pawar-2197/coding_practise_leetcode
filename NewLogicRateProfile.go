package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	//	"strconv"
)
type GenericData struct {
	rateProfileArr map[string]interface{}
	cost float64
	previousSlabHigh float64
	firstSlabFound bool
	tieredUnitSize float64
	pulsing []string
}

func rateProfileDerivation(cdr map[string]interface{}) error {
	var data GenericData
	data.rateProfileArr = make(map[string]interface{})
	if cdr["recordExtensions"] != nil {
		if cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] != nil {
			for _, chargeInformation := range cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}) {
				if chargeInformation.(map[string]interface{})["debitCash"] != nil{
					if chargeInformation.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"] != nil{
						for _, rateProfile := range chargeInformation.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}) {
							data.rateProfileArr[rateProfile.(map[string]interface{})["rateProfileId"].(string)] = rateProfile
						}
						lastAccumulatedRatedUsage :=  chargeInformation.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["lastAccumRatedUsage"].(float64)
						ratedUsage :=  chargeInformation.(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						rateProfileId := chargeInformation.(map[string]interface{})["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["rateProfileId"].(string)
						deriveRate(lastAccumulatedRatedUsage, ratedUsage, rateProfileId, &data)
						chargeInformation.(map[string]interface{})["cost"] = data.cost
						if len(data.pulsing) == 1{
							chargeInformation.(map[string]interface{})["pulsing"] = data.pulsing[0] +"/"+ data.pulsing[0]
							return nil
						}
						if len(data.pulsing) > 2{
							chargeInformation.(map[string]interface{})["pulsing"] = data.pulsing[len(data.pulsing)-2] + "/" + data.pulsing[len(data.pulsing)-1]
						}else{
							chargeInformation.(map[string]interface{})["pulsing"] = data.pulsing[0] + "/" + data.pulsing[1]w
						}
					}else{
						fmt.Println("RateProfile is absent in the current chargeInformation, moving to next chargeInformation")
						continue
					}
				}else{
					fmt.Println("DebitCash is absent in the current chargeInformation, moving to next chargeInformation")
					continue
				}
			}
		} else{
			return errors.New("ChargeInformation is absent in the cdr")
		}
	}else{
		return errors.New("RecordExtension is absent in the cdr")
	}
	return nil
}

func deriveRate(lastAccumulatedRatedUsage float64, ratedUsage float64, rateProfileId string, data *GenericData) error {
	rateProfile, ok := data.rateProfileArr[rateProfileId] ; if !ok {
		fmt.Println("RateProfile ", rateProfileId, " is absent in the cdr")
		return errors.New("RateProfile "+ rateProfileId +" is absent in the cdr")	
	}
	if lastAccumulatedRatedUsage == 0 {
		val, ok := rateProfile.(map[string]interface{})["lastAccumRatedUsage"].(float64)
		if ok {
			lastAccumulatedRatedUsage = val
		}
		if !ok && rateProfile.(map[string]interface{})["ratingType"].(string) == "TELESCOPIC"{
			fmt.Println("lastAccumRatedUsage is absent for telescopic rateProfile with rateProfileID ", rateProfileId)
			return errors.New("lastAccumRatedUsage is absent for telescopic rateProfile with rateProfileID " + rateProfileId)
		}
	}
	lenOfRatesArr := len(rateProfile.(map[string]interface{})["rates"].([]interface{}))
	for i, rate := range rateProfile.(map[string]interface{})["rates"].([]interface{}) {
		currentSlabLow , ok := rate.(map[string]interface{})["low"].(float64)
		if !ok{
			fmt.Println("Field low is absent in the rateSlab with rateProfileId ", rateProfileId)
			return errors.New("Field low is absent in the rateSlab with rateProfileId " + rateProfileId)
		}
		currentSlabHigh, ok := rate.(map[string]interface{})["high"].(float64)
		if !ok{
			fmt.Println("Field high is absent in the rateSlab with rateProfileId ", rateProfileId)
			return errors.New("Field high is absent in the rateSlab with rateProfileId " + rateProfileId)
		}
		var currentSlabUnitSize float64
		if rateProfile.(map[string]interface{})["ratingType"] != nil && rateProfile.(map[string]interface{})["ratingType"].(string) == "TELESCOPIC"{
			currentSlabUnitSize = data.tieredUnitSize
		}else{
			currentSlabUnitSize, ok = rate.(map[string]interface{})["unitSize"].(float64)
			if !ok{
				fmt.Println("Field unitSize is absent in the rateSlab with rateProfileId ", rateProfileId)
				return errors.New("Field unitSize is absent in the rateSlab with rateProfileId " + rateProfileId)
			}
			currentSlabRatedUnit, ok := rate.(map[string]interface{})["ratedUnit"].(string)
			if !ok {
				fmt.Println("Field ratedUnit is absent in the rateSlab with rateProfileId ", rateProfileId)
				return errors.New("Field ratedUnit is absent in the rateSlab with rateProfileId " + rateProfileId)
			}
			currentSlabUnitSize = currentSlabUnitSize * lookupData[currentSlabRatedUnit].ValueInLeastUnit 

			rate.(map[string]interface{})["ratedUnit"] = lookupData[currentSlabRatedUnit].LeastUnit
			rate.(map[string]interface{})["unitSize"] = currentSlabUnitSize
		}
		currentSlabRate, ok := rate.(map[string]interface{})["rate"].(float64)
		if !ok && rateProfile.(map[string]interface{})["rateProfileId"] == nil {
			fmt.Println("Field rate and rateProfileID both are absent in the rateSlab with rateProfileId ", rateProfileId)
			return errors.New("Field rate and rateProfileID both are absent in the rateSlab with rateProfileId " + rateProfileId)		
		}
		currentSlabUsage := (currentSlabHigh - data.previousSlabHigh) * currentSlabUnitSize
		if lastAccumulatedRatedUsage == 0 && i > 0{
			lastAccumulatedRatedUsage = lastAccumulatedRatedUsage + ratedUsage
		}
		tierCount := lastAccumulatedRatedUsage/currentSlabUnitSize
		if (tierCount + data.previousSlabHigh ) >= currentSlabLow && (tierCount + data.previousSlabHigh ) <= currentSlabHigh{
			data.firstSlabFound = true
			if rate.(map[string]interface{})["rateProfileId"] != nil {
				var usage float64
				if i == lenOfRatesArr-1{
					usage = lastAccumulatedRatedUsage
				}else{
					if currentSlabUsage - lastAccumulatedRatedUsage > ratedUsage{
						usage = ratedUsage
					}else{
						usage = math.Abs(lastAccumulatedRatedUsage - currentSlabUsage)
					}
				}
				if usage == 0 {
					lastAccumulatedRatedUsage = lastAccumulatedRatedUsage -currentSlabUsage
					continue
				}
				data.tieredUnitSize = currentSlabUnitSize
				data.previousSlabHigh = currentSlabHigh
				data.pulsing = append(data.pulsing,strconv.FormatFloat(currentSlabUnitSize, 'f', -1, 64))
				newData := GenericData{tieredUnitSize: data.tieredUnitSize, rateProfileArr: data.rateProfileArr}
				deriveRate(0 , usage ,rate.(map[string]interface{})["rateProfileId"].(string) , &newData)
				data.cost = data.cost + newData.cost
			//	data.pulsing = append(data.pulsing,newData.pulsing[0])
				lastAccumulatedRatedUsage = lastAccumulatedRatedUsage - currentSlabUsage
				if lastAccumulatedRatedUsage < 0 && ratedUsage > 0{
					lastAccumulatedRatedUsage = lastAccumulatedRatedUsage + ratedUsage
					ratedUsage = 0
				}
				if lastAccumulatedRatedUsage < 0{
					break
				}
			}else{
				var unitsToCharge float64
				if i == lenOfRatesArr-1{
					unitsToCharge = lastAccumulatedRatedUsage/currentSlabUnitSize
				}else{
					if currentSlabUsage - lastAccumulatedRatedUsage > ratedUsage{
						unitsToCharge = ratedUsage/currentSlabUnitSize
					}else{
						unitsToCharge = math.Abs(lastAccumulatedRatedUsage - currentSlabUsage)/currentSlabUnitSize
					}
				}
				data.cost = data.cost + unitsToCharge * currentSlabRate
				if rateProfile.(map[string]interface{})["ratingType"].(string) != "TELESCOPIC"{
					data.pulsing = append(data.pulsing, strconv.FormatFloat(currentSlabUnitSize, 'f', -1, 64))
				}
				lastAccumulatedRatedUsage = lastAccumulatedRatedUsage - currentSlabUsage
				if lastAccumulatedRatedUsage < 0 && ratedUsage > 0{
					lastAccumulatedRatedUsage = lastAccumulatedRatedUsage + ratedUsage
					ratedUsage = 0
				}
				if lastAccumulatedRatedUsage < 0{
					break
				}
				data.previousSlabHigh = currentSlabHigh
				continue
			}
		}else{
			if data.firstSlabFound{
				var unitsToCharge float64
				if i == lenOfRatesArr-1 || lastAccumulatedRatedUsage < currentSlabUsage{
					unitsToCharge = lastAccumulatedRatedUsage/currentSlabUnitSize
				}else{
					unitsToCharge = currentSlabUsage/currentSlabUnitSize
				}
				data.cost = data.cost + unitsToCharge * currentSlabRate
				if rateProfile.(map[string]interface{})["ratingType"].(string) != "TELESCOPIC"{
					data.pulsing = append(data.pulsing, strconv.FormatFloat(currentSlabUnitSize, 'f', -1, 64))
				}
				lastAccumulatedRatedUsage = lastAccumulatedRatedUsage - currentSlabUsage
				if lastAccumulatedRatedUsage < 0 && ratedUsage > 0{
					lastAccumulatedRatedUsage = lastAccumulatedRatedUsage + ratedUsage
					ratedUsage = 0
				}
				if lastAccumulatedRatedUsage < 0{
					break
				}
				data.previousSlabHigh = currentSlabHigh
				continue
			}else{
				lastAccumulatedRatedUsage = lastAccumulatedRatedUsage - currentSlabUsage
				data.previousSlabHigh = currentSlabHigh
				continue
			}
		}
	}	
	return nil
}

type ConversionData struct{
	ValueInLeastUnit float64
	LeastUnit string
}
var lookupData map[string]ConversionData

func main() {
	cdr := make(map[string]interface{})
	//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\NewRateProfileSample\rateProfileWithTwoSlabTiered.json
	//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\NewRateProfileSample\rateProfileWithTwoRateSlabs.json
	//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\NewRateProfileSample\rateProfileTiered.json
	//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\NewRateProfileSample\rateProfileWithTwoRateSlabs.json
	//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\NewRateProfileSample\rateProfileSampleTiered(LastAccISZERO).json
	readFileData, err := ioutil.ReadFile("C:/Users/pawarpr/OneDrive - Mavenir Systems, Inc/Documents/GoPractise and notes/NewRateProfileSample/rateProfileTiered.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	lookupData = make(map[string]ConversionData)
	lookupData["KB"] = ConversionData{1000, "BYTES"}
	lookupData["MB"] = ConversionData{1000000, "BYTES"}
	lookupData["MIN"] = ConversionData{60, "SEC"}
	lookupData["SEC"] = ConversionData{1, "SEC"}
	
	rateProfileDerivation(cdr)
	
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("finalRateProfileCDR",encodedData, 0644)
	
}