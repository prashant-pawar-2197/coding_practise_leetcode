package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"
	//	"time"
)

type RgPkgInfo struct {
	Package string
	Product string
}

//C:\Users\pawarpr\OneDrive - Mavenir Systems, Inc\Documents\GoPractise and notes\newMultipkg.json
func main() {
	cdr := make(map[string]interface{})
	var action CustomFunction
	readFileData, err := ioutil.ReadFile("multipkgsample2.json")
	if err != nil {
		fmt.Println("Error occured")
	}

	readActionData, err := ioutil.ReadFile("multipkgsample2SampleAction.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	err = json.Unmarshal(readActionData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}

	actionMultiPackageAggregationn(action, cdr, "transID")
	//durationCalculator(cdr, "")
	encodedData, _ := json.Marshal(cdr)
	ioutil.WriteFile("pulseSplitOutput", encodedData, 0644)
}

func panicHandler(funcName, transID string) {
	if err := recover(); err != nil {
		fmt.Println("An Error occured in function --", funcName, " error is --", err)
	}
}

type CustomFunction struct {
	FunctionName string          `json:"functionName"`
	Lookup       map[string]bool `json:"lookup,omitempty"`
	Data         []struct {
		Key   string `json:"key"`
		Value bool   `json:"value"`
	} `json:"data,omitempty"`
	BtpFlagEnabled     bool            `json:"btpFlagEnabled"`
	MultiPackageConfig MultiPackageCfg `json:"multiPackageConfig"`
}

type UnitDet struct {
	Value      float64 `json:"value"`
	MinuteUnit string  `json:"minuteUnit"`
}

type MultiPackageCfg struct {
	SplitOnPremiumTypes	[]struct {
        Key   string    `json:"key"`
        Value string 	`json:"value"`
    } `json:"splitOnPremiumTypes,omitempty"`
	FUPCopyPkgFromFirstOccranceOfNormalRate		bool	 `json:"fUPCopyPkgFromFirstOccranceOfNormalRate"`
	SplitOnPremiumTypesMap 						map[string]string
	RemoveChargeInfoFields struct {
		DeleteListType 				EnumListType
		ListType					string `json:"listType"`
		DeleteFields				[]string `json:"deleteFields"`
		ApplcablePremiumRateTypes	[]string `json:"applcablePremiumRateTypes"`	
		PremiumRateTypesMap			map[string]bool
	} `json:"removeChargeInfoFields,omitempty"`

	SplitOnPulseChange  bool `json:"splitOnPulseChange"`
    UnitLookup          map[string]UnitDet
    Units []struct {
        Key   string    `json:"key"`
        Value UnitDet   `json:"value"`
    } `json:"units,omitempty"`
    ProfilePriority     	string `json:"profilePriority"`
    ProfilePriorityEnum 	EnumProfilePriority
	TimeFormat				string `json:"timeFormat"`
	PremiumRateTypesLookup	map[string]string
	PremiumRateTypes 		[]string	`json:"premiumRateTypes"`
	RateProfilePrefixForReplacingToNuccc string `json:"rateProfilePrefixForReplacingToNuccc"`
}

type CopyableMap map[string]interface{}
type CopyableSlice []interface{}

// DeepCopy will create a deep copy of this map. The depth of this
// copy is all inclusive. Both maps and slices will be considered when
// making the copy.
func (m CopyableMap) DeepCopy() map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range m {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result[k] = CopyableMap(mapvalue).DeepCopy()
			continue
		}

		// Handle slices
		slicevalue, isSlice := v.([]interface{})
		if isSlice {
			result[k] = CopyableSlice(slicevalue).DeepCopy()
			continue
		}

		result[k] = v
	}

	return result
}

// DeepCopy will create a deep copy of this slice. The depth of this
// copy is all inclusive. Both maps and slices will be considered when
// making the copy.
func (s CopyableSlice) DeepCopy() []interface{} {
	result := []interface{}{}

	for _, v := range s {
		// Handle maps
		mapvalue, isMap := v.(map[string]interface{})
		if isMap {
			result = append(result, CopyableMap(mapvalue).DeepCopy())
			continue
		}

		// Handle slices
		slicevalue, isSlice := v.([]interface{})
		if isSlice {
			result = append(result, CopyableSlice(slicevalue).DeepCopy())
			continue
		}
		result = append(result, v)
	}

	return result
}

func (m *MultiPackageCfg) UnmarshalJSON(data []byte) error {
    type MultiPackageCfg2 MultiPackageCfg
    if err := json.Unmarshal(data, (*MultiPackageCfg2)(m)); err != nil {
        return err
    }

    m.ProfilePriorityEnum = ConvertProfilePriorityFromStringToEnum(m.ProfilePriority)
    m.UnitLookup = make(map[string]UnitDet)

    for _, fvalue := range m.Units {
        m.UnitLookup[fvalue.Key] = fvalue.Value
    }

	m.SplitOnPremiumTypesMap = make(map[string]string)
	for _, fvalue := range m.SplitOnPremiumTypes {
		m.SplitOnPremiumTypesMap[fvalue.Key] = fvalue.Value
    }
	m.RemoveChargeInfoFields.PremiumRateTypesMap = make(map[string]bool)
	for _, fvalue := range m.RemoveChargeInfoFields.ApplcablePremiumRateTypes {
		m.RemoveChargeInfoFields.PremiumRateTypesMap[fvalue] = true
	}
	m.PremiumRateTypesLookup = make(map[string]string)
	for _, v := range m.PremiumRateTypes {
		m.PremiumRateTypesLookup[v] = v
	}
	m.RemoveChargeInfoFields.DeleteListType = ConvertListTypeFromStringToEnum(m.RemoveChargeInfoFields.ListType)
    return nil
}

//Convert profile priority from string to enum
func ConvertProfilePriorityFromStringToEnum(typestr string) EnumProfilePriority {
	switch typestr {
	case "QUOTA_PROFILE":
		return QUOTA_PROFILE
	case "RATE_PROFILE":
		return RATE_PROFILE
	}
	return RATE_PROFILE
}

type EnumProfilePriority int

const (
	QUOTA_PROFILE = iota + 1
	RATE_PROFILE
)

type EnumListType int

func ConvertListTypeFromStringToEnum(typestr string) EnumListType {
	switch typestr {
	case "CGFM_DELETE_FROM_LIST":
		return CGFM_DELETE_FROM_LIST
	case "CGFM_ALLOWED_IN_LIST":
		return CGFM_ALLOWED_IN_LIST
	}
	return CGFM_DELETE_FROM_LIST
}

const (
	CGFM_DELETE_FROM_LIST = iota
	CGFM_ALLOWED_IN_LIST
)

func DeleteNestedMapField(doc map[string]interface{}, keyPath string) (error, interface{}) {
	pathVarList := strings.Split(keyPath, ".")
	err, val := DeleteNestedFieldFromMap(doc, pathVarList...)
	if err != nil {
		return err, nil
	}
	return err, val
}

//Function to delete subdocument from sub levels
func DeleteNestedFieldFromMap(doc map[string]interface{}, keyset ...string) (err error, rval interface{}) {
	var ok bool
	if len(keyset) == 0 { // degenerate input
		return errors.New("Delete-NestedValue failed. Path/Key list is empty"), nil
	}
	if rval, ok = doc[keyset[0]]; !ok {
		return errors.New("Key/Path not found"), nil
	} else if len(keyset) == 1 { // we've reached the final key
		delete(doc, keyset[0])
		return nil, rval
	} else if doc, ok = rval.(map[string]interface{}); !ok {
		return errors.New("Delete-NestedValue failed. Malformed structure"), nil
	} else { // 1+ more keys
		return DeleteNestedFieldFromMap(doc, keyset[1:]...)
	}
}

//Function to delete nested document
func DeleteNestedField(doc map[string]interface{}, keyPath string) (error, interface{}) {
	pathVarList := strings.Split(keyPath, ".")
	err, val := DeleteNestedValue(doc, pathVarList...)
	if err != nil {
		return err, nil
	}
	return err, val
}

//Function to insert subdocuments in sub levels.
func InsertNestedValue(doc interface{}, value interface{}, keyset ...string) (error, interface{}) {
	if len(keyset) == 0 { // degenerate input
		return errors.New("Nested-Insert failed. Path/Key list is empty"), nil
	}

	switch doc.(type) {
	case []interface{}:
		elements := doc.([]interface{})
		index, err := strconv.Atoi(keyset[0])
		if err != nil {
			return errors.New("Get-NestedValue failed. invalid number."), nil
		}
		if index > len(elements) {
			return errors.New("Get-NestedValue failed. Invalid array index"), nil
		}

		if len(keyset) == 1 {
			if index < len(elements) {
				elements[index] = value
			} else {
				elements = append(elements, value)
			}
			return nil, elements
		} else {
			isArray := false
			if _, err := strconv.Atoi(keyset[1]); err == nil {
				isArray = true
			}
			if index < len(elements) {
				return InsertNestedValue(elements[index], value, keyset[1:]...)
			} else {
				if isArray {
					newArray := make([]interface{}, 0, 0)
					elements = append(elements, newArray)
					err, opArr := InsertNestedValue(newArray, value, keyset[1:]...)
					if err == nil {
						elements[len(elements)-1] = opArr
					}
					return err, elements
				} else {
					newMap := make(map[string]interface{})
					elements = append(elements, newMap)
					err, _ := InsertNestedValue(elements[index], value, keyset[1:]...)
					return err, elements
				}
			}
		}
	case map[string]interface{}:
		docMap := doc.(map[string]interface{})
		if rval, ok := docMap[keyset[0]]; !ok {
			if len(keyset) == 1 {
				docMap[keyset[0]] = value
				return nil, nil
			} else {
				if _, err := strconv.Atoi(keyset[1]); err != nil {
					newMap := make(map[string]interface{})
					docMap[keyset[0]] = newMap
					return InsertNestedValue(newMap, value, keyset[1:]...)
				} else {
					newArray := make([]interface{}, 0, 0)
					docMap[keyset[0]] = newArray
					err, opArr := InsertNestedValue(newArray, value, keyset[1:]...)
					if opArr != nil {
						docMap[keyset[0]] = opArr
					}
					return err, nil
				}
			}
		} else if len(keyset) == 1 { // we've reached the final key
			docMap[keyset[0]] = value
			return nil, rval
		} else {
			switch rval.(type) {
			case map[string]interface{}:
				return InsertNestedValue(rval.(map[string]interface{}), value, keyset[1:]...)
			case []interface{}:
				err, opArr := InsertNestedValue(rval.([]interface{}), value, keyset[1:]...)
				if opArr != nil {
					docMap[keyset[0]] = opArr
				}
				return err, nil
			default:
				return errors.New("Nested-Insert failed. Malformed structure"), nil
			}
		}
	default:
		return errors.New("Nested-Insert failed. Malformed structure"), nil
	}
	return nil, nil
}

//Function to copy nested parameter from one map to another
func CopyNestedFields(src map[string]interface{}, dest map[string]interface{}, keyPath string) error {
	pathVarList := strings.Split(keyPath, ".")
	err, val := GetNestedValue(src, pathVarList...)
	if err != nil {
		return err
	}

	err, _ = InsertNestedValue(dest, val, pathVarList...)
	if err != nil {
		return err
	}
	return nil
}

//Function to get subdocument from sub levels
func GetNestedValue(doc interface{}, keyset ...string) (err error, rval interface{}) {
	var ok bool
	if len(keyset) == 0 { // degenerate input
		return errors.New("Get-NestedValue failed. Path/Key list is empty"), nil
	}

	switch doc.(type) {
	case []interface{}:
		elements := doc.([]interface{})
		index, err := strconv.Atoi(keyset[0])
		if err != nil {
			return errors.New("Get-NestedValue failed. invalid number."), nil
		}
		if index >= len(elements) {
			return errors.New("Get-NestedValue failed. Invalid array index"), nil
		}
		if len(keyset) == 1 {
			return err, elements[index]
		} else {
			return GetNestedValue(elements[index], keyset[1:]...)
		}
	case map[string]interface{}:
		subdoc := doc.(map[string]interface{})
		if rval, ok = subdoc[keyset[0]]; !ok {
			return errors.New("Get-NestedValue failed. key not found."), nil
		} else if len(keyset) == 1 { // we've reached the final key
			return nil, rval
		} else {
			switch rval.(type) {
			case map[string]interface{}:
				return GetNestedValue(rval.(map[string]interface{}), keyset[1:]...)
			case []interface{}:
				return GetNestedValue(rval.([]interface{}), keyset[1:]...)
			default:
				return errors.New("Get-NestedValue failed. Malformed structure"), nil
			}
		}
	default:
		return errors.New("DefaultCase: Get-NestedValue failed. Malformed structure"), nil
	}
}

//Function to get subdocument from sub levels
func DeleteNestedValue(doc interface{}, keyset ...string) (err error, rval interface{}) {
	var ok bool
	if len(keyset) == 0 { // degenerate input
		return errors.New("Delete-NestedValue failed. Path/Key list is empty"), nil
	}

	switch doc.(type) {
	case []interface{}:
		elements := doc.([]interface{})
		index, err := strconv.Atoi(keyset[0])
		if err != nil {
			return errors.New("Delete-NestedValue failed. invalid number."), nil
		}
		if index >= len(elements) {
			return errors.New("Delete-NestedValue failed. Invalid array index"), nil
		}
		if len(keyset) == 1 {
			val := elements[index]
			elements[index] = nil
			return err, val
		} else {
			return DeleteNestedValue(elements[index], keyset[1:]...)
		}
	case map[string]interface{}:
		subdoc := doc.(map[string]interface{})
		if rval, ok = subdoc[keyset[0]]; !ok {
			return errors.New("Delete-NestedValue failed. key not found."), nil
		} else if len(keyset) == 1 { // we've reached the final key
			delete(subdoc, keyset[0])
			return nil, rval
		} else {
			switch rval.(type) {
			case map[string]interface{}:
				return DeleteNestedValue(rval.(map[string]interface{}), keyset[1:]...)
			case []interface{}:
				return DeleteNestedValue(rval.([]interface{}), keyset[1:]...)
			default:
				return errors.New("Delete-NestedValue failed. Malformed structure"), nil
			}
		}
	default:
		return errors.New("DefaultCase: Delete-NestedValue failed. Malformed structure"), nil
	}
}

//get config value from custom config map
func getConfig(cfgmap map[string]bool, key string, defVal bool) bool {

	cfg, ok := cfgmap[key]
	if !ok {
		return defVal
	}
	return cfg
}

//GetPremiumType
func getPremiumType(chargeInfo map[string]interface{}, transID string) (premiumType string) {
	defer panicHandler("getPremiumType", transID)
	if debitCash, ok := chargeInfo["debitCash"].(map[string]interface{}); ok {
		if rateProfileArr, ok := debitCash["rateProfile"].([]interface{}); ok && len(rateProfileArr) > 0 {
			if rateProfile, ok := rateProfileArr[0].(map[string]interface{}); ok {
				if premiumType, ok := rateProfile["premiumRateType"].(string); ok {
					return premiumType
				}
			}
		}
	}
	return ""
}

//GetPremiumType
func getPackageProduct(chargeInfo map[string]interface{}, transID string) (pkg string, product string) {
	defer panicHandler("getPackageProduct", transID)
	if ratingIndication, ok := chargeInfo["ratingIndication"].(map[string]interface{}); ok {
		pkg, _ = ratingIndication["packageId"].(string)
		product, _ = ratingIndication["productId"].(string)
	}
	return pkg, product
}

//Copy new balance
func copyNewBalance(currBalChangeInfo, firstBalChangeInfo map[string]interface{}, transID string) {
	defer panicHandler("copyNewBalance", transID)
	if currBalChangeInfo["balanceType"] == "GC" {
		if newBalance, ok := currBalChangeInfo["newBalance"]; ok {
			firstBalChangeInfo["newBalance"] = newBalance
		}
	}
}

type MsccInfoData struct {
	totalVolume                       float64
	totalUplink                       float64
	totalDownlink                     float64
	finalTotalVol                     float64
	finalDownlink                     float64
	finalUplink                       float64
	arrIndexAtListOfMultipleUsageUnit int
	prevDuration                      float64
	remainingVol                      float64
	chargedFirstChargeInfo            bool
	recordOpeningTime                 string
}

type MsccInfoVoice struct {
	ccTime   float64
	prevTime float64
}

type MsccInfo struct {
	data  MsccInfoData
	voice MsccInfoVoice
}

type MultiPkgAdditionalInfo struct {
	msccRGMap         map[float64]*MsccInfo
	recordOpeningTime string
	newBalance        float64
	serviceType       string
	premiumRateType   string
}

//MultiPackageAggregation
func actionMultiPackageAggregationn(action CustomFunction, cdr map[string]interface{}, transID string) error {
    aggrTaxAmtArrFlg := getConfig(action.Lookup, "AggregateRateProfileArray", true)
    aggrTaxAmtFlg := getConfig(action.Lookup, "AggregateTaxAmt", true)
    aggrBalChangeDebitArrFlg := getConfig(action.Lookup, "AggregateBalanceChangeInfoArray", true)
    aggrBalChangeDebitFlg := getConfig(action.Lookup, "AggregateBalChangeDebitBal", true)
    aggrDiscountFlg := getConfig(action.Lookup, "AggregateDiscount", true)
    aggrDebitAllowanceAmtFlg := getConfig(action.Lookup, "AggregateDebitAllowanceAmount", true)
    aggrRateUsageFlg := getConfig(action.Lookup, "AggregateRateUsage", true)
    debugFlg := getConfig(action.Lookup, "debug", false)
    addMultiPkgFieldFlg := getConfig(action.Lookup, "AddMultiPkgIndicator", true)
    exitFromGppFlg := getConfig(action.Lookup, "ExitFromGpp", false)
	aggrIgnorePkgProductChangeFlag := getConfig(action.Lookup, "AggregateIgnorePkgProductChange", true)

    //There is CDR Format change for M1 & M2 release
    //M1 rate profile is not an array.
    //M2 rate profile is array.
    // Modified to handle release based on Configuration
    rateProfileIsArray := getConfig(action.Lookup, "RateProfileIsArray", true)
	DataRgMap := make(map[float64]*MsccInfoData, 0)

    if cdr["recordExtensions"] == nil {
        fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions parmeter not found in cdr")
        if exitFromGppFlg {
            return errors.New("CustomFunName:" + action.FunctionName + "recordExtensions parmeter not found in cdr")
        } else {
            return nil
        }
    }

	chargeInfoExtendedArr := make([]interface{}, 0)
    iChargeInformation, ok := cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]
    if !ok {
        fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation parmeter not found in cdr")
        if exitFromGppFlg {
            return errors.New("CustomFunName:" + action.FunctionName + "recordExtensions.chargeInformation  parmeter not found in cdr")
        } else {
            return nil
        }
    }
	var btpChangeFlag bool
		if action.BtpFlagEnabled {
		aggregateMultipleChargeInformationForBtpChange(cdr, &btpChangeFlag, transID)
		if btpChangeFlag{
			fmt.Println("NOT CONTINUING MULTIPACKAGE AGGR")
			return nil
		}
	}
	
    switch iChargeInformation.(type) {
    case []interface{}:
        rgMap := make(map[string]int)
		rgPkgMap := make(map[float64]RgPkgInfo)

        i := 0
        arrLen := len(iChargeInformation.([]interface{}))
		var multiPkgInfo MultiPkgAdditionalInfo
		listOfMultipleUnitUsage := make([]interface{},0)
		if action.MultiPackageConfig.SplitOnPulseChange {
			if multiPkgInfo.msccRGMap == nil {
				multiPkgInfo.msccRGMap = make(map[float64]*MsccInfo)
			}
			multiPkgInfo.serviceType, _ = cdr["basicService"].(string)
			multiPkgInfo.recordOpeningTime, _ = cdr["recordOpeningTime"].(string)

			if arrLen >= 1 {
				if lastChargeInfo, ok := iChargeInformation.([]interface{})[(arrLen-1)].(map[string]interface{}); ok {
					if debitCash, ok := lastChargeInfo["debitCash"].(map[string]interface{}); ok {
						if balChangeInfoArr, ok := debitCash["balanceChangeInfo"].([]interface{}); ok {
							for i := 0; i < len(balChangeInfoArr); i++ {
								if balChangeInfo, ok := balChangeInfoArr[i].(map[string]interface{}); ok {
									if balChangeInfo["balanceType"] == "GC" {
										multiPkgInfo.newBalance, _ = balChangeInfo["newBalance"].(float64)
										break
									}
								}
							}
						}
					}
				}
			}
			
			//Fetch MSCC info
			if listOfMultipleUnitUsage, ok = cdr["listOfMultipleUnitUsage"].([]interface{}); ok {
				for i, element := range listOfMultipleUnitUsage {
					if mscc, ok := element.(map[string]interface{}); ok {
						if ratingGroup, ok := mscc["ratingGroup"].(float64); ok {
							switch (multiPkgInfo.serviceType) {
								case "Voice": 
									if usedUnitContainerList, ok := mscc["usedUnitContainer"].([]interface{}); ok {
										for _, usedUnitContainerInf := range usedUnitContainerList {
											if usedUnitContainer, ok := usedUnitContainerInf.(map[string]interface{}); ok {
												if ccTime, ok := usedUnitContainer["time"].(float64); ok {
													multiPkgInfo.msccRGMap[ratingGroup] = &MsccInfo {
														voice : MsccInfoVoice {
															ccTime: ccTime, 
														},
													}
												}
											}
										}
									}
								case "Data", "MMS":
									if usedUnitContainerList, ok := mscc["usedUnitContainer"].([]interface{}); ok {
										for _, usedUnitContainerInf := range usedUnitContainerList {
											if usedUnitContainer, ok := usedUnitContainerInf.(map[string]interface{}); ok {
												totalVol, _ := usedUnitContainer["totalVolume"].(float64)
												downlinkVol, _ := usedUnitContainer["downlinkVolume"].(float64)
												uplinkVol, _ := usedUnitContainer["uplinkVolume"].(float64)
												recordOpenTime, _ := cdr["recordOpeningTime"].(string)
												DataRgMap[ratingGroup] = &MsccInfoData {
													totalVolume: totalVol,
													totalUplink: uplinkVol,
													totalDownlink: downlinkVol,
													recordOpeningTime: recordOpenTime,
													arrIndexAtListOfMultipleUsageUnit: i,
												}
											}
										}
									}
							}
						}
					}
				}
			}
		}

        for j := 0; j < arrLen ; j++ {
            ChargeInfoMap := iChargeInformation.([]interface{})[i].(map[string]interface{})
            //getRatingGroup(ChargeInfoMap)
            if ChargeInfoMap["ratingIndication"] == nil {
                fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation.ratingIndication parmeter not found in cdr. continue with next chargeInformation")
				//increment "i" index->Send as it is
				i++
                continue
            }

            ratingIndication, ok := ChargeInfoMap["ratingIndication"].(map[string]interface{})
            if !ok {
                fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation.ratingIndication parmeter not found in cdr. continue with next chargeInformation")
				//increment "i" index->Send as it is
				i++
                continue
            }

			var (
				rg_pkg_prd string
			)

            rg, ok := ratingIndication["ratingGroup"].(float64)
            if !ok {
                fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation.ratingIndication.ratingGroup parmeter not found in cdr. continue with next chargeInformation")
				//increment "i" index->Send as it is
				i++
                continue
            }
			//To replace premiumRateType to NUC-CC when rateProfileId starts with 91001
			replacePremiumRateTypeToNuccc(ChargeInfoMap, action.MultiPackageConfig.RateProfilePrefixForReplacingToNuccc, transID)
			premiumType := getPremiumType(ChargeInfoMap, transID)
			
			if aggrIgnorePkgProductChangeFlag == false {
				pkg, ok := ratingIndication["packageId"].(string)
				if !ok {
                	fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation.ratingIndication.packageId parmeter not found in cdr. Ignoring")
            	}
				prd, ok := ratingIndication["productId"].(string)
				if !ok {
                	fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation.ratingIndication.productId parmeter not found in cdr. Ignoring")
            	}
				if action.MultiPackageConfig.FUPCopyPkgFromFirstOccranceOfNormalRate {	
					if premiumType == "FUP" && pkg == "" {
						//get package and product:
						if pkgInfo, ok := rgPkgMap[rg]; ok {
							pkg = pkgInfo.Package
							prd = pkgInfo.Product
						} else {
							//traverse array and fetch
							l := i+1
							for k := j+1; k < arrLen ; k++ {
								pkg, prd = getPackageProduct(iChargeInformation.([]interface{})[l].(map[string]interface{}), transID)
								l++
								if pkg != "" {
									rgPkgMap[rg] = RgPkgInfo{ 
											Package : pkg,
											Product	: prd,
									}
									break
								}
							}
						}
					} else {
						//Store package and product
						rgPkgMap[rg] = RgPkgInfo{ 
								Package : pkg,
								Product	: prd,
						}
					}
				}
				rg_pkg_prd = strconv.FormatFloat(rg, 'f', -1, 64) + "_" + pkg + "_" + prd
			} else {
				rg_pkg_prd = strconv.FormatFloat(rg, 'f', -1, 64)
			}

			if premiumType != ""  {
				if value, ok := action.MultiPackageConfig.SplitOnPremiumTypesMap[premiumType]; ok {
					rg_pkg_prd += value
				}
			}
			if action.MultiPackageConfig.SplitOnPulseChange {
				multiPkgInfo.premiumRateType = premiumType
				chargeInfoExtendedArrTmp := actionSplitOnPulseChange(&action, ChargeInfoMap, &multiPkgInfo, transID)
				if chargeInfoExtendedArrTmp != nil || len(chargeInfoExtendedArrTmp) != 0 {
					fmt.Println("chargeinfo splitted based on Pulse. ChargeInfoArr:", chargeInfoExtendedArrTmp)
					chargeInfoExtendedArr = append(chargeInfoExtendedArr, chargeInfoExtendedArrTmp...)
				}
			}
			
			if _, ok := action.MultiPackageConfig.RemoveChargeInfoFields.PremiumRateTypesMap[premiumType]; ok {
				if action.MultiPackageConfig.RemoveChargeInfoFields.DeleteListType == CGFM_DELETE_FROM_LIST {
					for _, deleteField := range action.MultiPackageConfig.RemoveChargeInfoFields.DeleteFields {
						if err, _ := DeleteNestedField(ChargeInfoMap, deleteField); err != nil {
							fmt.Println("CustomFunName:", action.FunctionName, ".Failed to delete ", deleteField)
						}
					}
				} else {
					ChargeInformationNew := make(map[string]interface{})
					for _, deleteField := range action.MultiPackageConfig.RemoveChargeInfoFields.DeleteFields {
						if err := CopyNestedFields( ChargeInfoMap, ChargeInformationNew, deleteField); err != nil {
							fmt.Println("CustomFunName:", action.FunctionName, ".Failed to copy to ", deleteField)
						} 
					}
					ChargeInfoMap = ChargeInformationNew
				}	
			}

            if val, ok := rgMap[rg_pkg_prd]; ok {
                //do array level aggregation
                //debit amount.
                //tax amount
                //balance change info - debit amount

                fmt.Println("CustomFunName:", action.FunctionName, ". Duplicate ratingGroup found at arrIndex:", i, "Previous index:", val, "RatingGroup:" , rg)

                //Calculation of debit amount
                //firstInstanceOfChargeInfoWithSameRG -> ficiwsrg
                ficiwsrg := iChargeInformation.([]interface{})[val].(map[string]interface{})

				if ficiwsrg["ratingIndication"] == nil {
					ficiwsrg["ratingIndication"] = ChargeInfoMap["ratingIndication"]
				}
				ficiwsrgPremiumType := getPremiumType(ficiwsrg, transID)
				if ficiwsrgPremiumType == "NUC-CC" && premiumType == "NetworkUsage" {
					if ratingIndication, ok := ChargeInfoMap["ratingIndication"]; ok {
						ficiwsrg["ratingIndication"] = ratingIndication
					}
					debitCash, chargeInfoFlag := ChargeInfoMap["debitCash"].(map[string]interface{})
					ficiwsrgDebitCash, ficiwsrgFlag := ficiwsrg["debitCash"].(map[string]interface{})
					var (
						lenOfFiciwsrgRateProfile int
						lenOfRateProfile 		 int 
						rateProfile 			 []interface{}
						ficiwsrgRateProfile 	 []interface{}
						ok						 bool
					)
					if rateProfile, ok = debitCash["rateProfile"].([]interface{}); ok {
						lenOfRateProfile = len(rateProfile)
					}
					if ficiwsrgRateProfile, ok = ficiwsrgDebitCash["rateProfile"].([]interface{}); ok {
						lenOfFiciwsrgRateProfile = len(ficiwsrgRateProfile)
					}
					if chargeInfoFlag && ficiwsrgFlag && lenOfFiciwsrgRateProfile > 0 && lenOfRateProfile > 0 {
						// copying rateProfileID , taxid, premiumRateType
						if rateProfileId, ok := rateProfile[0].(map[string]interface{})["rateProfileId"]; ok {
							ficiwsrgRateProfile[0].(map[string]interface{})["rateProfileId"] = rateProfileId
						}
						if taxId, ok := rateProfile[0].(map[string]interface{})["taxId"]; ok {
							ficiwsrgRateProfile[0].(map[string]interface{})["taxId"] = taxId
						}
						if premiumRateType, ok := rateProfile[0].(map[string]interface{})["premiumRateType"]; ok {
							ficiwsrgRateProfile[0].(map[string]interface{})["premiumRateType"] = premiumRateType
						}
						// copying discountProfileId
						discountProfile, currentChargeInfoFlag :=  debitCash["discountProfile"].(map[string]interface{})
						ficiwsrgdiscountProfile, ficiwsrgFlag :=  ficiwsrgDebitCash["discountProfile"].(map[string]interface{})
						if currentChargeInfoFlag && ficiwsrgFlag {
							if discountProfileId, ok := discountProfile["discountProfileId"]; ok {
								ficiwsrgdiscountProfile["discountProfileId"] = discountProfileId
							}
						}
						// copying unitSize, ratedUnit, ratedUsage
						if ratedUnit , ok := ChargeInfoMap["ratedUnit"]; ok {
							ficiwsrg["ratedUnit"] = ratedUnit
						}
						if unitSize , ok := ChargeInfoMap["unitSize"]; ok {
							ficiwsrg["unitSize"] = unitSize
						}
						if ratedUsage , ok := debitCash["ratedUsage"]; ok {
							ficiwsrgDebitCash["ratedUsage"] = ratedUsage
							aggrRateUsageFlg = false
						}
					}
				}
		
				if action.MultiPackageConfig.SplitOnPulseChange {
					if premiumType != "" {
						if _, ok := action.MultiPackageConfig.PremiumRateTypesLookup[premiumType] ; !ok {
							premiumType = ""
						}
					}
					if premiumType == "" {
						if prevUnitSize, ok := ficiwsrg["unitSize"].(string); ok {
							if currUnitSize, ok := ChargeInfoMap["unitSize"].(string); ok {
								if prevUnitSize == currUnitSize {
									if multiPkgInfo.serviceType == "Voice" {
										if cDuration, ok := ChargeInfoMap["duration"].(float64); ok {
											if duration, ok := ficiwsrg["duration"].(float64); ok {
												ficiwsrg["duration"] = (cDuration + duration)
											} else {
												ficiwsrg["duration"] = cDuration
											}
										}
									}
									if rot, ok := ChargeInfoMap["recordOpeningTime"]; ok {
										ficiwsrg["recordOpeningTime"] = rot
									}
								} else {
									i++
									continue
								}
							} 
						}
					}
				}

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
                                                fmt.Println("CustomFunName:", action.FunctionName, ". aggregating debitAllowance.debitAmount:")
                                        }
                                    ficiwsrg["debitAllowance"].(map[string]interface{})["debitAmount"] = (debitAmount1.(float64) + debitAmount2.(float64))
                                } else {
                                    if debugFlg {
                                            fmt.Println("CustomFunName:", action.FunctionName, ". debit amount not present in first instance of RG. updating debitAllowance.debitAmount")
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
						copyIncludedConnectionChargeAndTaxAmount(ChargeInfoMap, transID)
						aggregateIncludedAdditionalCharge(action, ficiwsrg, ChargeInfoMap, premiumType)
                        ficiwsrg["debitCash"] =  ChargeInfoMap["debitCash"]
                        //Balance Change info
                        balanceChangeInfoArr, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                        if f2 {
							CopyGCToFirstIndex(balanceChangeInfoArr, transID)
                            if len(balanceChangeInfoArr) > 0 {
                                if aggrBalChangeDebitArrFlg {
                                    balanceChangeInfoArr = aggrBalChangeDebitAmountArray(balanceChangeInfoArr, transID, debugFlg)
                                }
							}
						}
						ficiwsrg["debitCash"].(map[string]interface{})["balanceChangeInfo"] = balanceChangeInfoArr
						
                    } else {

                        if rateProfileIsArray {
							// added on 03-08-2023
							if currentRateProfile, ok := ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}); ok {
								if len(currentRateProfile) > 0 {
									if currentRateProfile[0].(map[string]interface{})["premiumRateType"] == nil {
										if ChargeInfoMap["unitSize"] != "" {
											ficiwsrg["unitSize"] , _ = ChargeInfoMap["unitSize"]
										}
										if ChargeInfoMap["ratedUnit"] != "" {
											ficiwsrg["ratedUnit"] , _ = ChargeInfoMap["ratedUnit"]
										}	
										if rateProfileID, ok := currentRateProfile[0].(map[string]interface{})["rateProfileId"]; ok {
											if ficiwsrgRateProfile, ok := ficiwsrg["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}); ok {
												if len(ficiwsrgRateProfile) > 0 {
													ficiwsrgRateProfile[0].(map[string]interface{})["rateProfileId"] = rateProfileID
												}
											}
										}
									}
								}
							}	
							copyIncludedConnectionChargeAndTaxAmount(ChargeInfoMap, transID)
							aggregateIncludedAdditionalCharge(action, ficiwsrg, ChargeInfoMap, premiumType)
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
                                                        fmt.Println("CustomFunName:", action.FunctionName, ". aggregating tax amount:")
                                                    }
                                                    rateProfileArr1[0].(map[string]interface{})["taxAmount"] = num1 + num2
                                                } else {
                                                    if debugFlg {
                                                        fmt.Println("CustomFunName:", action.FunctionName, "tax Amount value not present in first instance of RG. Updating value from current instance")
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
                                                    fmt.Println("CustomFunName:", action.FunctionName, ". aggregating tax amount:")
                                                }
                                                rateProfile1.(map[string]interface{})["taxAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println("CustomFunName:", action.FunctionName, "tax Amount value not present in first instance of RG. Updating value from current instance")
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
                                                    fmt.Println("CustomFunName:", action.FunctionName, ". aggregating discountProfilediscountAmount:")
                                                }
                                                discountProfile1.(map[string]interface{})["discountAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println("CustomFunName:", action.FunctionName, "discountProfilediscountAmount value not present in first instance of RG. Updating value from current instance")
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
                                        fmt.Println("CustomFunName:", action.FunctionName, ". aggregating ratedUsage")
                                       }
                                    ficiwsrg["debitCash"].(map[string]interface{})["ratedUsage"] = num1 + num2
                                   } else {
                                    if debugFlg {
                                        fmt.Println("CustomFunName:", action.FunctionName, "discountProfilediscountAmount value not present in first instance of RG. Updating value from current instance")
                                       }
                                    ficiwsrg["debitCash"].(map[string]interface{})["ratedUsage"] = num1
                                   }
                              }
                        }
						// reseting rated usage flag
						aggrRateUsageFlg = getConfig(action.Lookup, "AggregateRateUsage", true)

                        //Balance Change info
                        balanceChangeInfoArr, f2 := ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                        if f2 {
							CopyGCToFirstIndex(balanceChangeInfoArr, transID)
                            balanceChangeInfoArr1, f1 := ficiwsrg["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{})
                            if len(balanceChangeInfoArr) > 0 {
                                if aggrBalChangeDebitArrFlg {
                                    balanceChangeInfoArr = aggrBalChangeDebitAmountArray(balanceChangeInfoArr, transID, debugFlg)
                                }

                                if aggrBalChangeDebitFlg {
                                    if f1 && len(balanceChangeInfoArr1) > 0 {
										copyNewBalance( balanceChangeInfoArr[0].(map[string]interface{}), balanceChangeInfoArr1[0].(map[string]interface{}),transID)
                                        if num1, ok := balanceChangeInfoArr[0].(map[string]interface{})["debitAmount"].(float64); ok {
                                            if num2, ok := balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"].(float64); ok {
                                                if debugFlg {
                                                    fmt.Println("CustomFunName:", action.FunctionName, ". aggregating balanceChange.debitAmount:")
                                                }
                                                balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"] = num1 + num2
                                            } else {
                                                if debugFlg {
                                                    fmt.Println("CustomFunName:", action.FunctionName, "balanceChange.debitAmount value not present in first instance of RG. Updating value from current instance")
                                                }
                                                balanceChangeInfoArr1[0].(map[string]interface{})["debitAmount"] = num1
                                            }
                                        }
                                    } else {
                                        balanceChangeInfoArr1 = balanceChangeInfoArr
                                    }
                                }
                                if aggrBalChangeDebitArrFlg || aggrBalChangeDebitFlg{
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
						copyIncludedConnectionChargeAndTaxAmount(ChargeInfoMap, transID)
						aggregateIncludedAdditionalCharge(action, nil, ChargeInfoMap, premiumType)
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
							CopyGCToFirstIndex(balanceChangeInfoArr, transID)
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
		
		if action.MultiPackageConfig.SplitOnPulseChange {
			if chargeInfoExtendedArr != nil || len(chargeInfoExtendedArr) != 0 {
				ciArr, _ := iChargeInformation.([]interface{})
				ciArr = append(ciArr, chargeInfoExtendedArr...)
				cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = ciArr
			}
		}
    default:
        fmt.Println("CustomFunName:", action.FunctionName, ". recordExtensions.chargeInformation invalid dataType")
        if exitFromGppFlg {
            return errors.New("CustomFunName:" + action.FunctionName + "recordExtensions.chargeInformation nvalid dataType")
        } else {
            return nil
        }
        return nil
    }

    if debugFlg {
            fmt.Println("CustomFunName:", action.FunctionName, "CDR after multi package aggregation. CDR->:", cdr )
    }
	if serviceType, ok := cdr["basicService"].(string); ok && serviceType != "" && (serviceType == "Data" || serviceType == "MMS") {
		dataRemainingCalculation(action, cdr, DataRgMap, transID)
	}
	fmt.Println("CDR after duration calculations -->", cdr)
	// err := durationCalculator(cdr, transID)
	// if err != nil{
	// 	return errors.New("Error occured while calculating duration")
	// }
    return nil
}

func checkForNormalRateIfExist(chargeInfoArr []interface{}, index, lenofChargeInfoArr int, transID string) bool {
	if index == lenofChargeInfoArr -1 {
		return false
	}
	nextChargeInfo := chargeInfoArr[index+1].(map[string]interface{})
	premiumType := getPremiumType(nextChargeInfo, transID)
	if premiumType == "Connection" || premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "NUC-CC" {
		return checkForNormalRateIfExist(chargeInfoArr, index+1, lenofChargeInfoArr, transID)
	} else {
		return true
	}
	return true
}

func dataRemainingCalculation(action CustomFunction, cdr map[string]interface{}, dataRgMap map[float64]*MsccInfoData, transID string) {
	defer panicHandler("dataRemainingCalculation", transID)
	if chargeInfoArr, ok := cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}); ok {
		lenOfChargeInfoArr := len(chargeInfoArr)
		if lenOfChargeInfoArr > 1 {
			for i, ChargeInfoMap := range chargeInfoArr {
				var (
					ratingIndication map[string]interface{}
					ok bool
					rg float64
					dataInfo *MsccInfoData
				)
				if ratingIndication, ok = ChargeInfoMap.(map[string]interface{})["ratingIndication"].(map[string]interface{}); !ok {
					fmt.Println("RatingIndication Block is absent in the current chargeInformation, continuing to the next chargeInformation")
					continue
				}
				if rg, ok = ratingIndication["ratingGroup"].(float64); !ok {
					fmt.Println("RatingIndication Block does not contain RG, continuing to the next chargeInformation")
					continue
				}
				if dataInfo, ok = dataRgMap[rg]; !ok {
					continue
				}
				if ChargeInfoMap.(map[string]interface{})["debitCash"] != nil {
					premiumType := getPremiumType(ChargeInfoMap.(map[string]interface{}), transID)
					if premiumType == "Connection" || premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "NUC-CC" {
						continue
					}
				}
				if dataInfo.chargedFirstChargeInfo {
					ChargeInfoMap.(map[string]interface{})["duration"] = 0
					ChargeInfoMap.(map[string]interface{})["downlink"] = dataInfo.remainingVol
					ChargeInfoMap.(map[string]interface{})["uplink"] = 0
					if reservationTime, ok := ratingIndication["reservationTime"].(string); ok {
						if startTimeGerman, err := SubtractTime(reservationTime, action.MultiPackageConfig.TimeFormat, transID, int64(-dataInfo.prevDuration)); err == nil {
							ChargeInfoMap.(map[string]interface{})["startTimeGerman"] = startTimeGerman
						}
					}
					ChargeInfoMap.(map[string]interface{})["writeOff"] = "true"
				} else if dataInfo.totalVolume > 0 && !dataInfo.chargedFirstChargeInfo {
					normalRateExist := checkForNormalRateIfExist(chargeInfoArr, i, lenOfChargeInfoArr, transID)
					if !normalRateExist {
						continue
					}
					var (
						currentUsage float64
						duration float64
						err error
					)
					//calculate duration
					if reservationTime, ok := ratingIndication["reservationTime"].(string); ok {
						duration, err = SubtractTwoTimeInUtcFormat(dataInfo.recordOpeningTime, reservationTime)
						if err != nil {
							fmt.Println("Error occured while calculating duration", err)
						}
						ChargeInfoMap.(map[string]interface{})["duration"] = duration
						dataInfo.prevDuration = duration
					}
					if debitAllowance, present := ChargeInfoMap.(map[string]interface{})["debitAllowance"].(map[string]interface{}); present {
						if currentUsage, ok = debitAllowance["serviceUsage"].(float64); !ok {
							currentUsage, _ = debitAllowance["debitAmount"].(float64)
						}
					}
					if debitCash, ok := ChargeInfoMap.(map[string]interface{})["debitCash"].(map[string]interface{}); ok {
						if ratedUsage, present := debitCash["ratedUsage"].(float64); present {
							currentUsage += ratedUsage
						}
					}
					if adjustedUnitAmt, ok := ratingIndication["adjustedUnitFromExtChargeAmt"].(float64); ok {
						currentUsage += adjustedUnitAmt
					}
					if dataInfo.totalVolume != 0.0 && currentUsage > 0 && dataInfo.totalVolume >= currentUsage {
						dataInfo.remainingVol = dataInfo.totalVolume - currentUsage
						dataInfo.finalTotalVol = dataInfo.totalVolume - dataInfo.remainingVol
					}
					if dataInfo.totalDownlink != 0.0 && dataInfo.remainingVol >= 0 && dataInfo.totalDownlink >= dataInfo.remainingVol {
						dataInfo.finalDownlink = dataInfo.totalDownlink - dataInfo.remainingVol
						dataInfo.finalUplink = dataInfo.totalUplink
					} else if dataInfo.totalUplink != 0.0 && dataInfo.remainingVol >= 0  && dataInfo.totalUplink >= dataInfo.remainingVol {
						dataInfo.finalUplink = dataInfo.totalUplink - dataInfo.remainingVol
						dataInfo.finalDownlink = dataInfo.totalDownlink
					}
					dataInfo.chargedFirstChargeInfo = true
					if currentUsage > 0 {
						copyUpdatedTotalVolAndDownlink(cdr, dataRgMap, dataInfo.arrIndexAtListOfMultipleUsageUnit)	
					}
				}
			}
		} else {
			return
		}
	}
}

func copyUpdatedTotalVolAndDownlink(cdr map[string]interface{}, DataRgMap map[float64]*MsccInfoData, arrIndex int) {
	if listOfMultipleUnitUsage, ok := cdr["listOfMultipleUnitUsage"].([]interface{}); ok {
		if arrIndex < len(listOfMultipleUnitUsage) {
			if element, ok := listOfMultipleUnitUsage[arrIndex].(map[string]interface{}); ok {
				if ratingGroup, ok := element["ratingGroup"].(float64); ok {
					if dataInfo, ok := DataRgMap[ratingGroup]; ok {
						if usedUnitContainerList, ok := element["usedUnitContainer"].([]interface{}); ok {
							if len(usedUnitContainerList) > 0 {
								usedUnitContainerList[0].(map[string]interface{})["totalVolume"] = dataInfo.finalTotalVol
								usedUnitContainerList[0].(map[string]interface{})["downlinkVolume"] = dataInfo.finalDownlink
								usedUnitContainerList[0].(map[string]interface{})["uplinkVolume"] = dataInfo.finalUplink
							}
						}
					}
				}
			}
		}
	}
}

func SubtractTwoTimeInUtcFormat(firstTimeStr, secondTimeStr string) (float64, error) {
	firstTime, err := time.Parse(time.RFC3339, firstTimeStr)
	if err != nil {
		return 0, err
	}

	secondTime, err := time.Parse(time.RFC3339, secondTimeStr)
	if err != nil {
		return 0, err
	}

	diff := firstTime.Sub(secondTime).Seconds()
	return math.Floor(diff), nil
}

func copyIncludedConnectionChargeAndTaxAmount(ChargeInfoMap map[string]interface{}, transID string) {
	defer panicHandler("copyIncludedConnectionChargeAndTaxAmount", transID)
	if rateProfile, ok := ChargeInfoMap["debitCash"].(map[string]interface{})["rateProfile"].([]interface{}); ok {
		if len(rateProfile) > 0 {
			if val, ok := rateProfile[0].(map[string]interface{})["premiumRateType"]; ok && val == "Connection" {
				if balChangeInfoArr, ok := ChargeInfoMap["debitCash"].(map[string]interface{})["balanceChangeInfo"].([]interface{}); ok {
					if len(balChangeInfoArr) > 0 {
						if debitAmount, ok := balChangeInfoArr[0].(map[string]interface{})["debitAmount"]; ok {
							ChargeInfoMap["includedConnectionCharge"] = debitAmount
						}
						if taxAmount, ok := rateProfile[0].(map[string]interface{})["taxAmount"]; ok {
							ChargeInfoMap["includedConnectionChargeTaxAmount"] = taxAmount
						}
					}
				}
			}
		}
	}
}

// function to replace premium rate type to NUC-CC for rateProfileID starting with 91001
func replacePremiumRateTypeToNuccc(ChargeInfoMap map[string]interface{}, rateProfilePrefix, transID string){
	defer panicHandler("replacePremiumRateTypeToNuccc", transID)
	if debitCash , ok := ChargeInfoMap["debitCash"]; ok {
		if rateProfile, ok := debitCash.(map[string]interface{})["rateProfile"].([]interface{}) ; ok {
			if len(rateProfile) > 0 {
				if premiumRateType, ok := rateProfile[0].(map[string]interface{})["premiumRateType"].(string); ok {
					if val , ok := rateProfile[0].(map[string]interface{})["rateProfileId"].(string) ; ok && strings.HasPrefix(val, rateProfilePrefix) && premiumRateType == "Connection" {
						rateProfile[0].(map[string]interface{})["premiumRateType"] = "NUC-CC"
					}
				}
			}
		}
	}
}

// function to aggregate Included Additional charges from FUP, NUC, NUC-CC- and to determine the
func aggregateIncludedAdditionalCharge(action CustomFunction, ficiwsrg, ChargeInfoMap map[string]interface{}, premiumType string) {
	if ficiwsrg == nil || ficiwsrg["debitCash"] == nil {
		var debitAmountSum float64
		if debitCash, ok := ChargeInfoMap["debitCash"].(map[string]interface{}); ok {
			if premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "Connection" || premiumType == "NUC-CC" {
				if balChangeInfoArr, ok := debitCash["balanceChangeInfo"].([]interface{}); ok {
					for _, balChange := range balChangeInfoArr {
						if debitAmount, ok := balChange.(map[string]interface{})["debitAmount"].(float64); ok {
							debitAmountSum = debitAmountSum + debitAmount
						}
					}
					ChargeInfoMap["debitCash"].(map[string]interface{})["includedAdditionalCharges"] = debitAmountSum
				}
				if discountProfile, ok := debitCash["discountProfile"].(map[string]interface{}); ok {
					if val, ok := discountProfile["discountAmount"].(float64); ok {
						ChargeInfoMap["debitCash"].(map[string]interface{})["includedAdditionalCharges_discountAmount"] = val
					}
				}
				if rateProfile, ok := debitCash["rateProfile"].([]interface{}); ok {
					if len(rateProfile) > 0 {
						if taxAmt, ok := rateProfile[0].(map[string]interface{})["taxAmount"]; ok {
							ChargeInfoMap["debitCash"].(map[string]interface{})["includedAdditionalCharges_taxAmount"] = taxAmt
						}
					}
				}
			}
		}
	} else {
		var debitAmountSum float64
		if debitCash, ok := ChargeInfoMap["debitCash"].(map[string]interface{}); ok {
			if premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "Connection" || premiumType == "NUC-CC" {
				if balChangeInfoArr, ok := debitCash["balanceChangeInfo"].([]interface{}); ok {
					for _, balChange := range balChangeInfoArr {
						if debitAmount, ok := balChange.(map[string]interface{})["debitAmount"].(float64); ok {
							debitAmountSum = debitAmountSum + debitAmount
						}
					}
					if debitCashFi, ok := ficiwsrg["debitCash"].(map[string]interface{}); ok {
						if inclAddCharge, ok := debitCashFi["includedAdditionalCharges"].(float64); ok {
							ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges"] = inclAddCharge + debitAmountSum
						} else {
							ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges"] = debitAmountSum
						}
					}
				}
				if discountProfile, ok := debitCash["discountProfile"].(map[string]interface{}); ok {
					if val, ok := discountProfile["discountAmount"].(float64); ok {
						if debitCashFi, ok := ficiwsrg["debitCash"].(map[string]interface{}); ok {
							if inclAddChargeDiscountAmt, ok := debitCashFi["includedAdditionalCharges_discountAmount"].(float64); ok {
								ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges_discountAmount"] = inclAddChargeDiscountAmt + val
							} else {
								ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges_discountAmount"] = val
							}
						}
					}
				}
				if rateProfile, ok := debitCash["rateProfile"].([]interface{}); ok {
					if len(rateProfile) > 0 {
						if taxAmt, ok := rateProfile[0].(map[string]interface{})["taxAmount"].(float64); ok {
							if debitCashFi, ok := ficiwsrg["debitCash"].(map[string]interface{}); ok {
								if inclAddChargeTaxAmt, ok := debitCashFi["includedAdditionalCharges_taxAmount"].(float64); ok {
									ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges_taxAmount"] = inclAddChargeTaxAmt + taxAmt
								} else {
									ficiwsrg["debitCash"].(map[string]interface{})["includedAdditionalCharges_taxAmount"] = taxAmt
								}
							}
						}
					}
				}
			}
		}
	}
}

func getRatedUsageFromDebitAllowance(ciWrapper *ChargeInfoWrapper) (float64, error) {
	var (
		ratedUsage 	float64
		ok 			bool
	)
	if ciWrapper.debitAllowance != nil {
		if ratedUsage, ok = ciWrapper.debitAllowance["serviceUsage"].(float64); ok {
			return ratedUsage, nil
		} else if ratedUsage, ok = ciWrapper.debitAllowance["debitAmount"].(float64); ok {
			return ratedUsage, nil
		}
	}
	return 0.0, errors.New("Usage not found for debitAllowance")
}

func CopyGCToFirstIndex(balanceChangeInfoArr []interface{}, transID string) {
	blanceLen := len(balanceChangeInfoArr)
	var index int = -1
	for i := 0; i < blanceLen; i++ {
		balanceMap, ok := balanceChangeInfoArr[i].(map[string]interface{})
		if ok && balanceMap["balanceType"] == "GC" {
			index = i
			break
		}
	}
	if index > 0 {
		balanceChangeInfoArr[index], balanceChangeInfoArr[0] = balanceChangeInfoArr[0], balanceChangeInfoArr[index]
	}
	fmt.Println("balanceChangeInfoArr: ", balanceChangeInfoArr)

}
func aggrTaxAmountArray(rateProfileList []interface{}, transID string, debugFlg bool) []interface{} {
	var sum float64
	counter := 0

	if debugFlg {
		fmt.Println("TaxAmountArray Aggregation enabled", rateProfileList)
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
				fmt.Println("After tax amount arr aggregation. rateProfile:", rateProfileList)
			}
		}
	}
	return rateProfileList
}

func aggrBalChangeDebitAmountArray(balChangeList []interface{}, transID string, debugFlg bool) []interface{} {
	var sum float64
	counter := 0
	if debugFlg {
		fmt.Println("BalChange debit amount Array Aggregation enabled", balChangeList)
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
				fmt.Println("After balance Change debit amount arr aggregation. BalChangeLst:", balChangeList)
			}
		}
	}
	return balChangeList
}

func durationCalculator(cdr map[string]interface{}, transID string) error {
	if cdr["recordExtensions"] == nil {
		fmt.Println("RecordExtension is nil")
		return nil
	}

	ChargeInformation, ok := cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{})
	if !ok {
		fmt.Println("CustomFunName: DurationCalculator recordExtensions.chargeInformation parmeter not found in cdr")
		return nil
	}
	var newChargeInformationArr []interface{}
	for _, val := range ChargeInformation {
		if val != nil {
			newChargeInformationArr = append(newChargeInformationArr, val)
		}
	}
	totalTime, ok := cdr["listOfMultipleUnitUsage"].([]interface{})[0].(map[string]interface{})["usedUnitContainer"].([]interface{})[0].(map[string]interface{})["time"].(float64)
	if !ok {
		fmt.Println("Time Field is missing in the cdr, or this cdr is not a voice cdr")
		return nil
	}
	lenOfNewChargeInfoArr := len(newChargeInformationArr)
	for i := 0; i < lenOfNewChargeInfoArr; i++ {
		if newChargeInformationArr[i] == nil {
			continue
		}
		if i == lenOfNewChargeInfoArr-1 {
			if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitCash"]; ok {
				newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = totalTime
				break
			} else if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitAllowance"]; ok {
				newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = totalTime
				break
			}
		}
		if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitCash"]; ok {
			newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = newChargeInformationArr[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"]
			totalTime -= newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"].(float64)
		} else if _, ok := newChargeInformationArr[i].(map[string]interface{})["debitAllowance"]; ok {
			newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"] = newChargeInformationArr[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"]
			totalTime -= newChargeInformationArr[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["modifiedDuration"].(float64)
		}
	}
	cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = newChargeInformationArr
	return nil
}

// function which will aggregate all the chargeInformation with btpChangeFlag as "No"
// and finally we will have only two chargeInformations in chargeInformation array which will consist of
// chargeInformation with btpChangeFlag as "No" and chargeInformation with btpChangeFlag as "Yes"
func aggregateMultipleChargeInformationForBtpChange(cdr map[string]interface{}, btpChangeFlag *bool, transID string) {
	defer panicHandler("aggregateMultipleChargeInformationForBtpChange", transID)
	// Creating an array which will store chargeInformation with btpChangeFlag No and  chargeInformation with btpChangeFlag Yes
	// Once both chargeInformations are determined then aggregate all chargeinformations with btpChangeFlag as No
	finalChargeInforArray := make([]interface{}, 2)
	count := len((cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{}))
	for i := 0; i < count; i++ {
		if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"]; ok {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No" {
					// when the current chargeInformation is the first chargeInformation with btpChangeFLag as No (debitCash)
					if finalChargeInforArray[0] == nil {
						finalChargeInforArray[0] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						*btpChangeFlag = true
						continue
					} else {
						//aggregating into the existing chargeInformation when chargeInformation contains debitCash
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						*btpChangeFlag = true
						continue
					}
					// when the current chargeInformation is the first chargeInformation with btpChangeFLag as Yes (debitCash)
				} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes" {
					finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
					finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
					(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
					*btpChangeFlag = true
					continue
				}
			}
		} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitAllowance"]; ok {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No" {
					// when the current chargeInformation is the first chargeInformation with btpChangeFLag as No (debitAllowance)
					if finalChargeInforArray[0] == nil {
						finalChargeInforArray[0] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						*btpChangeFlag = true
						continue
					} else {
						//aggregating into the existing chargeInformation when chargeInformation contains debitAllowance
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						*btpChangeFlag = true
						continue
					}
				} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "Yes" {
					// when we have reached the last chargeInformation in the chargeInformation array then we need to subtract the aggregated modifiedDuration from CCTIME for finding the duration for the ChargeInformation with btpChangeCycle as yes
					if i == count-1 {
						finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
						finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
						(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
						*btpChangeFlag = true
						break
					}
					// when the current chargeInformation is the first chargeInformation with btpChangeFLag as Yes (debitAllowance)
					finalChargeInforArray[1] = (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i]
					finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = ((cdr["listOfMultipleUnitUsage"].([]interface{}))[0].(map[string]interface{})["usedUnitContainer"].([]interface{}))[0].(map[string]interface{})["time"].(float64)
					(cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] = nil
					*btpChangeFlag = true
					break
				}
			}

		}
	}
	for i := 0; i < count; i++ {
		if (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i] != nil {
			if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitCash"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No" {
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitCash"].(map[string]interface{})["ratedUsage"].(float64)
						finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[1].(map[string]interface{})["modifiedDuration"].(float64) - finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64)
					}
				}
			} else if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["debitAllowance"]; ok {
				if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"]; ok {
					if _, ok := (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string); ok && (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"]).([]interface{})[i].(map[string]interface{})["ratingIndication"].(map[string]interface{})["changeInBTPCycle"].(string) == "No" {
						finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"] = finalChargeInforArray[0].(map[string]interface{})["modifiedDuration"].(float64) + (cdr["recordExtensions"].(map[string]interface{})["chargeInformation"].([]interface{}))[i].(map[string]interface{})["debitAllowance"].(map[string]interface{})["debitAmount"].(float64)
					}
				}
			}
		}
	}
	if finalChargeInforArray[0] != nil && finalChargeInforArray[1] != nil {
		cdr["recordExtensions"].(map[string]interface{})["chargeInformation"] = finalChargeInforArray
	}
}

type EnumProfileType int

const (
	PROFILE_NOT_FOUND EnumProfileType = iota
	PROFILE_DA_QP_RP
	PROFILE_DA_RP
	PROFILE_DA_QP
	PROFILE_DA
	PROFILE_DA_DC_QP
	PROFILE_DA_DC
	PROFILE_DC_RP
	PROFILE_DC
)

//Combination:
// --> DebitAllowance WithQuota Profile + Rate profile
// --> DebitAllowance WithoutQuota Profile + Rate profile
// --> Rate profile
// --> DebitAllowance WithQuota Profile
// --> DebitAllowance WithoutQuota Profile

//Get rate profile
func getDebitCashProfile(chargeInformation map[string]interface{}, transID string) (map[string]interface{}, EnumProfileType) {
	defer panicHandler("getDebitCashProfile", transID)
	if debitCash, ok := chargeInformation["debitCash"].(map[string]interface{}); ok {
		if _, ok := debitCash["rateProfile"].([]interface{}); ok {
			return debitCash, PROFILE_DC_RP
		} else {
			//Only balance change info
			return debitCash, PROFILE_DC
		}
	}
	return nil, PROFILE_NOT_FOUND
}

//Get Quota profile
func getDebitAllowanceaProfile(chargeInformation map[string]interface{}, transID string) (map[string]interface{}, EnumProfileType) {
	defer panicHandler("getDebitAllowanceaProfile", transID)
	if debitAllowance, ok := chargeInformation["debitAllowance"].(map[string]interface{}); ok {
		if _, ok := debitAllowance["quotaProfile"]; ok {
			return debitAllowance, PROFILE_DA_QP
		} else {
			return debitAllowance, PROFILE_DA
		}
	}
	return nil, PROFILE_NOT_FOUND
}

//GetProfileType
func getProfile(chargeInformation map[string]interface{}, transID string) (map[string]interface{}, map[string]interface{}, EnumProfileType) {
	debitCash, rpProfileFlg := getDebitCashProfile(chargeInformation, transID)
	debitAllowance, daProfileFlg := getDebitAllowanceaProfile(chargeInformation, transID)

	switch rpProfileFlg {
	case PROFILE_DC_RP:
		switch daProfileFlg {
		case PROFILE_DA_QP:
			return debitCash, debitAllowance, PROFILE_DA_QP_RP
		case PROFILE_DA:
			return debitCash, debitAllowance, PROFILE_DA_RP
		case PROFILE_NOT_FOUND:
			return debitCash, nil, PROFILE_DC_RP
		}
	case PROFILE_DC:
		switch daProfileFlg {
		case PROFILE_DA_QP:
			return debitCash, debitAllowance, PROFILE_DA_DC_QP
		case PROFILE_DA:
			return debitCash, debitAllowance, PROFILE_DA_DC
		case PROFILE_NOT_FOUND:
			return debitCash, nil, PROFILE_DC
		}
	case PROFILE_NOT_FOUND:
		//PROFILE_DA_QP, PROFILE_DA, PROFILE_NOT_FOUND
		return debitCash, debitAllowance, daProfileFlg
	}
	return nil, nil, PROFILE_NOT_FOUND
}

//Check is pulse are different. if different need to split.
func isRateChange(action *CustomFunction, rateProfile map[string]interface{}) (bool, bool) {
	rateChange := false
	unitChange := false
	tierUsage, ok := rateProfile["tierUsage"].([]interface{})
	if len(tierUsage) <= 1 || !ok {
		//tierUsage or tierIndex not specified . Take 1st index of rate profile
		return rateChange, unitChange
	}

	tierUsageFirstSlab, ok := tierUsage[0].(map[string]interface{})
	if ok {
		pRange, ok := tierUsageFirstSlab["range"].(map[string]interface{})
		if !ok {
			return rateChange, unitChange
		}
		unitSize, _ := pRange["unitSize"].(float64)
		ratedUnit, _ := pRange["ratedUnit"].(string)
		for i := 1; i < len(tierUsage); i++ {
			cTierUsageFirstSlab, ok := tierUsage[i].(map[string]interface{})
			if !ok {
				continue
			}
			cRange, ok := cTierUsageFirstSlab["range"].(map[string]interface{})
			if !ok {
				continue
			}
			cUnitSize, _ := cRange["unitSize"].(float64)
			cRatedUnit, _ := cRange["ratedUnit"].(string)

			if cRatedUnit != ratedUnit {
				unitChange = true
				unitDet, udFlg1 := action.MultiPackageConfig.UnitLookup[ratedUnit]
				unitDet2, udFlg2 := action.MultiPackageConfig.UnitLookup[cRatedUnit]
				if udFlg1 && udFlg2 {
					unitSize = unitSize * unitDet.Value
					cUnitSize = cUnitSize * unitDet2.Value
					if unitSize != cUnitSize {
						rateChange = true
						break
					}
				} else {
					rateChange = true
				}
			}

			if unitSize != cUnitSize {
				rateChange = true
			}
		}
	}
	return rateChange, unitChange
}

//Action get pulse and unit.
func getPulseAndUnit(action *CustomFunction, tierUsageRange map[string]interface{}, unitChange bool) (float64, string) {
	prevPulse, _ := tierUsageRange["unitSize"].(float64)
	prevUnit, _ := tierUsageRange["ratedUnit"].(string)

	if unitChange {
		//convert to minute unit
		if unitDet, ok := action.MultiPackageConfig.UnitLookup[prevUnit]; ok {
			prevPulse = prevPulse * unitDet.Value
			return prevPulse, unitDet.MinuteUnit
		}
	}
	return prevPulse, prevUnit
}

//action to spilt CDR based on rate profile change.
func actionSplitOnRateProfileChange(ciWrapper *ChargeInfoWrapper, rateProfile map[string]interface{}, unitChange bool) []interface{} {
	chargeInfoExtendedArr := make([]interface{}, 0)
	if tierUsage, ok := rateProfile["tierUsage"].([]interface{}); ok {
		tierUsagePrev, ok := tierUsage[0].(map[string]interface{})
		if !ok {
			return nil
		}

		tierUsagePrevRange, ok := tierUsagePrev["range"].(map[string]interface{})
		if !ok {
			tierUsagePrevRange = make(map[string]interface{})
		}

		prevPulse, prevUnit := getPulseAndUnit(ciWrapper.action, tierUsagePrevRange, unitChange)
		var rangeDetails RateDetails
		addRatesWithCommonUnit(ciWrapper, tierUsagePrev, tierUsagePrevRange, &rangeDetails, prevPulse, prevUnit)

		for i := 1; i < len(tierUsage); i++ {
			tierUsageCurr, ok := tierUsage[i].(map[string]interface{})
			if !ok {
				continue
			}
			tierUsageCurrRange, ok := tierUsageCurr["range"].(map[string]interface{})
			if !ok {
				continue
			}
			currPulse, currUnit := getPulseAndUnit(ciWrapper.action, tierUsageCurrRange, unitChange)

			if currPulse == prevPulse {
				//Aggregate
				addRatesWithCommonUnit(ciWrapper, tierUsageCurr, tierUsageCurrRange, &rangeDetails, prevPulse, prevUnit)
				continue
			} else {
				//Calculate remaining balance, debit amount, starttime, tax amount, discount amount, ratedUsage
				chargeInformationTmp := CopyableMap(ciWrapper.chargeInformation).DeepCopy()
				actionCopyRateProfileParameter(ciWrapper, &rangeDetails)
				rangeDetails.reset()
				chargeInfoExtendedArr = append(chargeInfoExtendedArr, chargeInformationTmp)
				ciWrapper.chargeInformation = chargeInformationTmp
				if ciWrapper.debitCash != nil {
					if debitCash, ok := ciWrapper.chargeInformation["debitCash"].(map[string]interface{}); ok {
						ciWrapper.debitCash = debitCash
					}
				}
				addRatesWithCommonUnit(ciWrapper, tierUsageCurr, tierUsageCurrRange, &rangeDetails, currPulse, currUnit)
				prevPulse = currPulse
			}
		}
		actionCopyRateProfileParameter(ciWrapper, &rangeDetails)
	}
	return chargeInfoExtendedArr
}

//Rate details
type RateDetails struct {
	debitAmt   *float64
	taxAmt     *float64
	discount   *float64
	unitSize   *float64
	ratedUnit  string
	newBalance *float64
	ratedUsage *float64
}

func (r *RateDetails) reset() {
	r.debitAmt = nil
	r.taxAmt = nil
	r.discount = nil
	r.unitSize = nil
	r.ratedUnit = ""
	r.newBalance = nil
	r.ratedUsage = nil
}

//Copy rate profile parameters
func addRatesWithCommonUnit(ciWrapper *ChargeInfoWrapper, tierUsage map[string]interface{}, tierUsageRange map[string]interface{}, ratedetails *RateDetails, unitSize float64, ratedUnit string) {
	//Copy debit amount
	tierDebitAmt, ok := tierUsage["debitAmount"].(float64)
	if ok {
		if ratedetails.debitAmt == nil {
			ratedetails.debitAmt = &tierDebitAmt
		} else {
			*ratedetails.debitAmt += tierDebitAmt
		}
	}

	//Copy Discount
	tierDiscount, ok := tierUsage["discountAmount"].(float64)
	if ok {
		if ratedetails.discount == nil {
			ratedetails.discount = &tierDiscount
		} else {
			*ratedetails.discount += tierDiscount
		}
	}

	//Copy tax amount
	tierTaxAmt, ok := tierUsage["taxAmount"].(float64)
	if ok {
		if ratedetails.taxAmt == nil {
			ratedetails.taxAmt = &tierTaxAmt
		} else {
			*ratedetails.taxAmt += tierTaxAmt
		}
	}

	if ratedetails.unitSize == nil {
		ratedetails.unitSize = &unitSize
	}

	if ratedetails.newBalance == nil {
		ratedetails.newBalance = &ciWrapper.multiPkgInfo.newBalance
	}
	ratedetails.ratedUnit = ratedUnit

	//Copy tax amount
	ratedUsage, ok := tierUsage["ratedUsage"].(float64)
	if ok {
		if ratedetails.ratedUsage == nil {
			ratedetails.ratedUsage = &ratedUsage
		} else {
			*ratedetails.ratedUsage += ratedUsage
		}
	}
}

//Copy rate profile parameters
func actionCopyRateProfileParameter(ciWrapper *ChargeInfoWrapper, ratedetails *RateDetails) {
	//Copy debit amount
	if debitCash, ok := ciWrapper.chargeInformation["debitCash"].(map[string]interface{}); ok {
		if ratedetails.debitAmt != nil {
			if balChangeInfo, ok := debitCash["balanceChangeInfo"].([]interface{}); ok {
				for i := 0; i < len(balChangeInfo); i++ {
					if balChangeInfoMap, ok := balChangeInfo[i].(map[string]interface{}); ok {
						if balChangeInfoMap["balanceType"] == "GC" {
							balChangeInfoMap["debitAmount"] = *ratedetails.debitAmt
							//Reset newBalance
							if ratedetails.newBalance != nil {
								balChangeInfoMap["newBalance"] = *ratedetails.newBalance
							}
						}
					}
				}
			}
		}

		//Copy Discount
		discountInfo, FndFlg := debitCash["discountProfile"].(map[string]interface{})
		if FndFlg {
			if ratedetails.discount != nil {
				discountInfo["discountAmount"] = *ratedetails.discount
			} else {
				delete(discountInfo, "discountAmount")
			}
		}
	}

	if rateProfileArr, ok := ciWrapper.debitCash["rateProfile"].([]interface{}); ok {
		//Len check already done.
		if rateProfile, ok := rateProfileArr[0].(map[string]interface{}); ok {
			//Copy tax amount
			if ratedetails.taxAmt != nil {
				rateProfile["taxAmount"] = *ratedetails.taxAmt
			} else {
				delete(rateProfile, "taxAmount")
			}
		}
	}

	//set pulseUnit and size.
	if ratedetails.unitSize != nil {
		ciWrapper.chargeInformation["unitSize"] = strconv.Itoa(int(*ratedetails.unitSize)) + "/" + strconv.Itoa(int(*ratedetails.unitSize))
	}
	ciWrapper.chargeInformation["ratedUnit"] = ratedetails.ratedUnit
	if ratedetails.ratedUsage != nil {
		if ciWrapper.multiPkgInfo.serviceType == "Voice" {
			if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
				SetDurationAndTime(ciWrapper, mscc, *ratedetails.ratedUsage)
			}
		}
		if debitCash, ok := ciWrapper.chargeInformation["debitCash"].(map[string]interface{}); ok {
			debitCash["ratedUsage"] = *ratedetails.ratedUsage
		}
	}

}

type ChargeInfoWrapper struct {
	debitCash         map[string]interface{}
	debitAllowance    map[string]interface{}
	chargeInformation map[string]interface{}
	multiPkgInfo      *MultiPkgAdditionalInfo
	transID           string
	action            *CustomFunction
	ratingGroup       float64
	prevPulse         string
	prevUnit          string
}

//Function to get rating group
func getRatingGroup(chargeInformation map[string]interface{}, transID string) float64 {
	defer panicHandler("actionCustomFunction", transID)
	ratingGroup, _ := chargeInformation["ratingIndication"].(map[string]interface{})["ratingGroup"].(float64)
	return ratingGroup
}

//Split CDR on pulse change
func actionSplitOnPulseChange(action *CustomFunction, chargeInformation map[string]interface{}, multiPkgInfo *MultiPkgAdditionalInfo, transID string) []interface{} {

	var profileType EnumProfileType
	ciWrapper := &ChargeInfoWrapper{
		chargeInformation: chargeInformation,
		multiPkgInfo:      multiPkgInfo,
		transID:           transID,
		action:            action,
	}
	ciWrapper.ratingGroup = getRatingGroup(chargeInformation, transID)
	ciWrapper.debitCash, ciWrapper.debitAllowance, profileType = getProfile(chargeInformation, transID)

	if multiPkgInfo.premiumRateType != "" {
		fmt.Println("DerivePulse - Invoked. premiumRateType:", multiPkgInfo.premiumRateType, "ProfileType:", profileType, "ChargeInfo:", ciWrapper.chargeInformation)
		//In case of premium rate type populate pulse directly
		switch action.MultiPackageConfig.ProfilePriorityEnum {
		case QUOTA_PROFILE:
			switch profileType {
			case PROFILE_DA_QP, PROFILE_DA_DC_QP, PROFILE_DA_QP_RP:
				if pulse, unit, ok := getPulseUnitFromDAQP(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			case PROFILE_DA, PROFILE_DA_DC:
				if pulse, unit, ok := getPulseUnitFromDA(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			case PROFILE_DC_RP, PROFILE_DA_RP:
				if pulse, unit, ok := getPulseUnitFromRP(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			}
		case RATE_PROFILE:
			switch profileType {
			case PROFILE_DA_QP_RP, PROFILE_DA_RP, PROFILE_DC_RP:
				if pulse, unit, ok := getPulseUnitFromRP(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			case PROFILE_DA_QP, PROFILE_DA_DC_QP:
				if pulse, unit, ok := getPulseUnitFromDAQP(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			case PROFILE_DA, PROFILE_DA_DC:
				if pulse, unit, ok := getPulseUnitFromDA(ciWrapper); ok {
					setUnitSizeAndPulse(ciWrapper, pulse, unit)
				}
			}
		}
		return nil
	}
	fmt.Println("DerivePulse - Invoked. ProfileType:", profileType, "ChargeInfo:", ciWrapper.chargeInformation)

	switch profileType {
	case PROFILE_DA_QP, PROFILE_DA_DC_QP:
		//Consumed from Debit Allownce -> Pulse will be derived from QuotaProfile.
		//Debit cash contain balance change info
		return processCDRWithDebitAllowanceQP(ciWrapper)
	case PROFILE_DA, PROFILE_DA_DC:
		//Consumed from Debit Allownce
		//Debit cash contain balance change info
		return processCDRWithDebitAllowance(ciWrapper)
	case PROFILE_DC_RP:
		//Either debit allowance not present or part of other charge information
		return processCDRWithDebitCashRP(ciWrapper)
	case PROFILE_DC, PROFILE_NOT_FOUND:
		//No need to handle -> will be clubbed with previous charge information
		return nil
	case PROFILE_DA_QP_RP:
		//Consumed from debit allowance as well as RateProfile - For DebitAllowance Pulse will be derived from QuotaProfile.
		return processCDRWithDebitAllowanceRP(ciWrapper, false)
	case PROFILE_DA_RP:
		//Consumed from debit allowance as well as RateProfile
		return processCDRWithDebitAllowanceRP(ciWrapper, true)
	}
	return nil
}

//SetDuration and Time
func SetDurationAndTime(ciWrapper *ChargeInfoWrapper, mscc *MsccInfo, ratedUsage float64) {
	if ciWrapper.multiPkgInfo.serviceType == "Voice" {
		if ratingIndication, ok := ciWrapper.chargeInformation["ratingIndication"].(map[string]interface{}); ok {
			if adjustedUnitAmt, ok := ratingIndication["adjustedUnitFromExtChargeAmt"].(float64); ok {
				ratedUsage += adjustedUnitAmt
			}
		}
		ciWrapper.chargeInformation["duration"] = ratedUsage
		if mscc.voice.ccTime > ratedUsage {
			ciWrapper.chargeInformation["duration"] = ratedUsage
			mscc.voice.ccTime -= ratedUsage
		} else {
			ciWrapper.chargeInformation["duration"] = mscc.voice.ccTime
			mscc.voice.ccTime = 0
		}

		if recordOpeningTime, err := SubtractTime(ciWrapper.multiPkgInfo.recordOpeningTime, ciWrapper.action.MultiPackageConfig.TimeFormat, ciWrapper.transID, int64(mscc.voice.ccTime)); err == nil {
			ciWrapper.chargeInformation["recordOpeningTime"] = recordOpeningTime
		}
	}
}

//Function to get pulse and unit
func getPulseUnitFromDA(ciWrapper *ChargeInfoWrapper) (string, string, bool) {
	if ratedUnit, ok := ciWrapper.debitAllowance["ratedUsageUnit"].(string); ok {
		return "1/1", ratedUnit, true
	}
	return "1/1", "", false
}

//Function to get pulse and unit
func getPulseUnitFromDAQP(ciWrapper *ChargeInfoWrapper) (string, string, bool) {
	if quotaProfile, ok := ciWrapper.debitAllowance["quotaProfile"].(map[string]interface{}); ok {
		if rates, ok := quotaProfile["rates"].([]interface{}); ok {
			if len(rates) == 0 {
				return "", "", false
			}
			if len(rates) == 1 {
				return getPulseUnitFromRatesWithIndex(rates, 0)
			} else {
				return getPulseUnitFromRatesTelescopic(ciWrapper, rates)
			}
		}
	}
	return "", "", false
}

//Function to get pulse and unit
func getPulseUnitFromRP(ciWrapper *ChargeInfoWrapper) (string, string, bool) {
	if rateProfileList, ok := ciWrapper.debitCash["rateProfile"].([]interface{}); ok {
		if rateProfile, ok := rateProfileList[0].(map[string]interface{}); ok {
			if rates, ok := rateProfile["rates"].([]interface{}); ok {
				if len(rates) == 0 {
					return "", "", false
				}
				if len(rates) == 1 {
					return getPulseUnitFromRatesWithIndex(rates, 0)
				} else {
					return getPulseUnitFromRatesTelescopic(ciWrapper, rates)
				}
			}
		}
	}
	return "", "", false
}

//Function to get pulse and unit
func getPulseUnitFromRatesTelescopic(ciWrapper *ChargeInfoWrapper, rates []interface{}) (string, string, bool) {
	if ratesMap, ok := rates[0].(map[string]interface{}); ok {
		if ratesMap2, ok := rates[1].(map[string]interface{}); ok {
			unitSize, _ := ratesMap["unitSize"].(float64)
			unitSize2, _ := ratesMap2["unitSize"].(float64)

			if ratedUnit1, ok := ratesMap["ratedUnit"].(string); ok {
				if ratedUnit2, ok := ratesMap2["ratedUnit"].(string); ok {
					if ratedUnit1 == ratedUnit2 && unitSize == unitSize2 {
						return (strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))), ratedUnit1, true
					} else if ratedUnit1 == ratedUnit2 && unitSize != unitSize2 {
						return (strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize2))), ratedUnit1, true
					} else {
						unitDet, udFlg1 := ciWrapper.action.MultiPackageConfig.UnitLookup["ratedUnit1"]
						unitDet2, udFlg2 := ciWrapper.action.MultiPackageConfig.UnitLookup["ratedUnit2"]
						if udFlg1 && udFlg2 {
							unitSize = unitSize * unitDet.Value
							unitSize2 = unitSize2 * unitDet2.Value
							return (strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize2))), unitDet.MinuteUnit, true
						} else {
							return (strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))), ratedUnit1, true
						}
					}
				}
			}
		}
	}
	return "", "", false
}

//Function to get pulse and unit
func getPulseUnitFromRatesWithIndex(rates []interface{}, index int) (string, string, bool) {
	if ratesMap, ok := rates[index].(map[string]interface{}); ok {
		if unitSize, ok := ratesMap["unitSize"].(float64); ok {
			if ratedUnit, ok := ratesMap["ratedUnit"].(string); ok {
				return (strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))), ratedUnit, true
			}
		}
	}
	return "", "", false
}

//Function to set unit size and pulse
func setUnitSizeAndPulse(ciWrapper *ChargeInfoWrapper, pulse, unit string) {
	ciWrapper.prevPulse = pulse
	ciWrapper.prevUnit = unit
	ciWrapper.chargeInformation["unitSize"] = pulse
	ciWrapper.chargeInformation["ratedUnit"] = unit
}

//Function to get rate and pulse
func CopyPulseDurationFromRateProfile(ciWrapper *ChargeInfoWrapper, rates []interface{}, tierIndex int) {
	pulse, unit, ok := getPulseUnitFromRatesWithIndex(rates, tierIndex)
	if ok {
		setUnitSizeAndPulse(ciWrapper, pulse, unit)
	}

	ratedUsage, ok := ciWrapper.debitCash["ratedUsage"].(float64)
	if ok {
		ciWrapper.chargeInformation["ratedUsage"] = ratedUsage
		if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
			SetDurationAndTime(ciWrapper, mscc, ratedUsage)
		}
	}
}

//Function to get rate and pulse
func CopyPulseDurationFromRateProfileTeleScopic(ciWrapper *ChargeInfoWrapper, rates []interface{}) {
	pulse, unit, ok := getPulseUnitFromRatesTelescopic(ciWrapper, rates)
	if ok {
		setUnitSizeAndPulse(ciWrapper, pulse, unit)
	}

	ratedUsage, ok := ciWrapper.chargeInformation["ratedUsage"].(float64)
	if ok {
		ciWrapper.chargeInformation["ratedUsage"] = ratedUsage
		if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
			SetDurationAndTime(ciWrapper, mscc, ratedUsage)
		}
	}
}

//Function to get First Pulse Unit From Rate Profile
func getFirstPulseUnitFromRP(ciWrapper *ChargeInfoWrapper) (string, string, bool) {
	rateProfileList, ok := ciWrapper.debitCash["rateProfile"].([]interface{})
	if !ok || len(rateProfileList) == 0 {
		return "", "", false
	}

	rateProfile, ok := rateProfileList[0].(map[string]interface{})
	if !ok {
		return "", "", false
	}

	rateType, ok := rateProfile["ratingType"].(string)
	if !ok {
		rateType = "TELESCOPIC"
	}

	switch rateType {
	case "TIERED":
		if tierIndex, ok := rateProfile["tierIndex"].(float64); ok {
			//Used from Single Tier. Set PULSE.
			tierIndexInt := int(tierIndex) - 1
			if rates, ok := rateProfile["rates"].([]interface{}); ok {
				return getPulseUnitFromRatesWithIndex(rates, tierIndexInt)
			}
		}

		if tierUsage, ok := rateProfile["tierUsage"].([]interface{}); ok {
			if len(tierUsage) >= 1 {
				if tierUsageEle, ok := tierUsage[0].(map[string]interface{}); ok {
					if rates, ok := tierUsageEle["range"].(map[string]interface{}); ok {
						ratesArray := []interface{}{rates}
						return getPulseUnitFromRatesWithIndex(ratesArray, 0)
					}
				}
			}
		}

	case "TELESCOPIC":
		if rates, ok := rateProfile["rates"].([]interface{}); ok {
			if len(rates) == 1 {
				return getPulseUnitFromRatesWithIndex(rates, 0)
			} else if len(rates) > 1 {
				return getPulseUnitFromRatesTelescopic(ciWrapper, rates)
			}
		}
	}
	return "", "", false
}

//Function handle debit allowance request
func processCDRWithDebitAllowance(ciWrapper *ChargeInfoWrapper) []interface{} {
	//Scenario 1: CDR with only debit allowance charge information
	//Scenario 2: CDR with debit allowance and debit cash information - where debit cash contains only balance info
	//Scenario 3: CDR with debit allowance and debit cash information with rate profile -> both bucket used for charging

	//Normal rate charge information
	//First instance of RG Charge info is received.
	// Store reference so it will be easy to manipulate previous charge information
	pulse, unit, _ := getPulseUnitFromDA(ciWrapper)
	setUnitSizeAndPulse(ciWrapper, pulse, unit)
	ratedUsage, err := getRatedUsageFromDebitAllowance(ciWrapper)
	if err == nil {
		//ciWrapper.chargeInformation["ratedUsage"] = ratedUsage
		if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
			SetDurationAndTime(ciWrapper, mscc, ratedUsage)
		}
	}
	return nil
}

//Function handle debit allowance request
func processCDRWithDebitAllowanceQP(ciWrapper *ChargeInfoWrapper) []interface{} {
	//Scenario 1: CDR with only debit allowance charge information
	//Scenario 2: CDR with debit allowance and debit cash information - where debit cash contains only balance info
	//Scenario 3: CDR with debit allowance and debit cash information with rate profile -> both bucket used for charging

	if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
		//Normal rate charge information
		//First instance of RG Charge info is received.
		// Store reference so it will be easy to manipulate previous charge information
		pulse, unit, _ := getPulseUnitFromDAQP(ciWrapper)
		setUnitSizeAndPulse(ciWrapper, pulse, unit)
		ratedUsage, err := getRatedUsageFromDebitAllowance(ciWrapper)
		if err == nil {
			//ciWrapper.chargeInformation["ratedUsage"] = ratedUsage
			SetDurationAndTime(ciWrapper, mscc, ratedUsage)
		}
	}
	return nil
}

// Function to copy unitsize and rate
func CopyUnitSizeAndRate(chargeInformation map[string]interface{}, rates []interface{}, index int, transID string) {
	if len(rates) > index && index >= 0 {
		if ratesMap, ok := rates[index].(map[string]interface{}); ok {
			if unitSize, ok := ratesMap["unitSize"].(float64); ok {
				chargeInformation["unitSize"] = strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))
			}

			if ratedUnit, ok := ratesMap["ratedUnit"]; ok {
				chargeInformation["ratedUnit"] = ratedUnit
			}
		}
	}
}

// Derive Quota profile
func DeriveQuotaProfile(action *CustomFunction, chargeInformation map[string]interface{}, quotaProfile map[string]interface{}, transID string) {
	//ratingType is always tiered
	if rates, ok := quotaProfile["rates"].([]interface{}); ok {
		if len(rates) == 1 {
			CopyUnitSizeAndRate(chargeInformation, rates, 0, transID)
		} else if len(rates) > 1 {
			CopyUnitSizeAndRateTelescopic(action, chargeInformation, rates, transID)
		}
	}
}

// Function to copy unitsize and rate
func CopyUnitSizeAndRateTelescopic(action *CustomFunction, chargeInformation map[string]interface{}, rates []interface{}, transID string) {
	if ratesMap, ok := rates[0].(map[string]interface{}); ok {
		if ratesMap2, ok := rates[1].(map[string]interface{}); ok {
			unitSize, _ := ratesMap["unitSize"].(float64)
			unitSize2, _ := ratesMap2["unitSize"].(float64)

			if ratedUnit1, ok := ratesMap["ratedUnit"]; ok {
				if ratedUnit2, ok := ratesMap2["ratedUnit"]; ok {
					if ratedUnit1 == ratedUnit2 {
						chargeInformation["ratedUnit"] = ratedUnit1
						chargeInformation["unitSize"] = strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))
					} else {
						unitDet, udFlg1 := action.MultiPackageConfig.UnitLookup["ratedUnit1"]
						unitDet2, udFlg2 := action.MultiPackageConfig.UnitLookup["ratedUnit2"]
						if udFlg1 && udFlg2 {
							unitSize = unitSize * unitDet.Value
							unitSize2 = unitSize2 * unitDet2.Value
							chargeInformation["ratedUnit"] = unitDet.MinuteUnit
							chargeInformation["unitSize"] = strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize2))
						} else {
							chargeInformation["ratedUnit"] = ratedUnit1
							chargeInformation["unitSize"] = strconv.Itoa(int(unitSize)) + "/" + strconv.Itoa(int(unitSize))
						}
					}
				}
			}
		}
	}
}

//Function handle debit allowance request
func processCDRWithDebitAllowanceRP(ciWrapper *ChargeInfoWrapper, onlyDA bool) []interface{} {
	//Scenario 1: CDR with only debit allowance charge information
	//Scenario 2: CDR with debit allowance and debit cash information - where debit cash contains only balance info
	//Scenario 3: CDR with debit allowance and debit cash information with rate profile -> both bucket used for charging
	//Normal rate charge information
	//First instance of RG Charge info is received.
	// Store reference so it will be easy to manipulate previous charge information
	var pulse, unit string

	if onlyDA {
		pulse, unit, _ = getPulseUnitFromDA(ciWrapper)
	} else {
		pulse, unit, _ = getPulseUnitFromDAQP(ciWrapper)
	}

	rpPulse, rpUnit, _ := getFirstPulseUnitFromRP(ciWrapper)
	if pulse == rpPulse && unit == rpUnit {
		//Pulse is same
		ratedUsage, err := getRatedUsageFromDebitAllowance(ciWrapper)
		if err == nil {
			if mscc, ok := ciWrapper.multiPkgInfo.msccRGMap[ciWrapper.ratingGroup]; ok {
				if ciWrapper.multiPkgInfo.serviceType == "Voice" {
					if mscc.voice.ccTime > ratedUsage {
						mscc.voice.ccTime -= ratedUsage
					} else {
						mscc.voice.ccTime = 0
					}
				}
			}
		}
		return processCDRWithDebitCashRP(ciWrapper)
	} else {
		//Pulse is different
		chargeInfoExtendedArr := make([]interface{}, 0)
		chargeInformationTmp := CopyableMap(ciWrapper.chargeInformation).DeepCopy()
		delete(ciWrapper.chargeInformation, "debitCash")
		processCDRWithDebitAllowance(ciWrapper)
		chargeInfoExtendedArr = append(chargeInfoExtendedArr, chargeInformationTmp)

		ciWrapper.chargeInformation = chargeInformationTmp
		delete(ciWrapper.chargeInformation, "debitAllowance")
		chargeInfoExtendedArrTmp := processCDRWithDebitCashRP(ciWrapper)
		chargeInfoExtendedArr = append(chargeInfoExtendedArr, chargeInfoExtendedArrTmp...)
		return chargeInfoExtendedArr
	}
	return nil
}

//Function handle debit allowance request
func processCDRWithDebitCashRP(ciWrapper *ChargeInfoWrapper) []interface{} {
	rateProfileList, ok := ciWrapper.debitCash["rateProfile"].([]interface{})
	if !ok || len(rateProfileList) == 0 {
		return nil
	}

	rateProfile, ok := rateProfileList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rateType, ok := rateProfile["ratingType"].(string)
	if !ok {
		rateType = "TELESCOPIC"
	}

	switch rateType {
	case "TIERED":
		if tierIndex, ok := rateProfile["tierIndex"].(float64); ok {
			//Used from Single Tier. Set PULSE.
			tierIndexInt := int(tierIndex) - 1
			if rates, ok := rateProfile["rates"].([]interface{}); ok {
				CopyPulseDurationFromRateProfile(ciWrapper, rates, tierIndexInt)
				return nil
			}
		}

		rateChange, unitChange := isRateChange(ciWrapper.action, rateProfile)
		if rateChange == false {
			if tierUsage, ok := rateProfile["tierUsage"].([]interface{}); ok {
				if len(tierUsage) >= 1 {
					if tierUsageEle, ok := tierUsage[0].(map[string]interface{}); ok {
						if rates, ok := tierUsageEle["range"].(map[string]interface{}); ok {
							ratesArray := []interface{}{rates}
							CopyPulseDurationFromRateProfile(ciWrapper, ratesArray, 0)
						}
					}
				}
			}
			return nil
		}

		ChargeInformationExtended := actionSplitOnRateProfileChange(ciWrapper, rateProfile, unitChange)
		if ChargeInformationExtended == nil || len(ChargeInformationExtended) == 0 {
			return nil
		}
		return ChargeInformationExtended

	case "TELESCOPIC":
		if rates, ok := rateProfile["rates"].([]interface{}); ok {
			if len(rates) == 1 {
				CopyPulseDurationFromRateProfile(ciWrapper, rates, 0)
				return nil
			} else if len(rates) > 1 {
				CopyPulseDurationFromRateProfileTeleScopic(ciWrapper, rates)
				return nil
			}
		}
	}
	return nil
}

//function to subtract the given duration(in seconds) from the inputTime
func SubtractTime(inputTime, timeFormat, transID string, duration int64) (string, error) {
	actualFormatedTime := GetDynamicCNotationDateTimeSubstutions(timeFormat)
	tm, err := time.Parse(actualFormatedTime, inputTime)
	if err != nil {
		fmt.Println("failed to parse time value:", inputTime, " with format:", actualFormatedTime, ",error:", err)
		return "", errors.New("failed to parse time value:" + inputTime + " with format:" + actualFormatedTime + ",error:" + err.Error())
	}
	tm = tm.Add(time.Second * time.Duration(-duration))
	finalUtcTime := tm.Format(GetDynamicCNotationDateTimeSubstutions(timeFormat))
	return finalUtcTime, nil
}

func GetDynamicCNotationDateTimeSubstutions(input string) string {
	for key, value := range CNotationDateTimeSubstution {
		input = strings.Replace(input, key, value, -1)
	}
	return input
}

func GetCurrentCNotationDateTimeSubstutions(input string) string {
	t := time.Now()
	return t.Format(GetDynamicCNotationDateTimeSubstutions(input))
}

var CNotationDateTimeSubstution = map[string]string{
	"%YEAR%":      "2006",    //full year
	"%year%":      "06",      //short year
	"%MONTH%":     "January", //full char month
	"%month_nz%":  "01",      //num month with zero
	"%month_n%":   "1",       //num month with out zero
	"%month%":     "Jan",     //short char month
	"%WEEK_DAY%":  "Monday",  //day of week in full char
	"%week_day%":  "Mon",     //day of week in short char
	"%date_z%":    "02",      //date with zero
	"%date%":      "2",       //date with out zero
	"%hour_24%":   "15",      //hours in 24 hour format
	"%hour_12_z%": "03",      //hours in 12 hour format with zero
	"%hour_12%":   "3",       //hours in 12 hour format with out zero
	"%minute_z%":  "04",      //minute with zero
	"%minute%":    "4",       //minute with out zero
	"%second_z%":  "05",      //second with zero
	"%second%":    "5",       //second with out zero
}
