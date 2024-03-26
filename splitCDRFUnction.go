package main

import (
	"encoding/json"
//	"errors"
	"fmt"
	"io/ioutil"
	parser "ruleEngine/parser"
	"sort"
	"strconv"
	"strings"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func main(){
	cdr := make(map[string]interface{})
	var action SplitCdr
	readFileData, err := ioutil.ReadFile("splitcdrSampleData.json")
	if err != nil {
		fmt.Println("Error occured")
	}

	readActionData, err := ioutil.ReadFile("splitaction.json")
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
	splitCdr(action, cdr, "", "transid")
	encodedData, _ := json.Marshal(cdr)
	ioutil.WriteFile("pulseSplitOutput", encodedData, 0644)
}


type SplitCdr struct {
	ArrayNameList[] struct {
		ArrayName string `json: "arrayName"`
    } `json:"arrayNameList"`
	Condition string `json:"condition,omitempty"`
	MsgInterface MsgInterfaceConfig `json:"msgInterface"`
	SctList []string `json:"sctList"`
	SctLookup    map[string]string
	SctLookupKey	string	`json:"sctLookupKey"`
	ChangeSctName bool	`json:"changeSctName"`
	Data              []struct {
		Key   string	`json:"key"`
		Value string	`json:"value"`
	} `json:"data"`
	TransactionTypeChangeKeyFields []string `json:"transactionTypeChangeKeyFields"`
	ChangePrevOccuranceTxnTypeToUpdateOnSameRGOnRelease bool `json:"ChangePrevOccuranceTxnTypeToUpdateOnSameRGOnRelease"`
}

type MsgInterfaceConfig struct {
	OutInterface	string `json:"outInterface"`
	OutInterfaceType	EnumOutputInterface `json:"outInterfaceType,omitempty"`
	Es 		EsConfig `json:"es,omitempty"`
	Kafka	KafkaConfig `json:"kafka,omitempty"`
}

type EnumOutputInterface int

const (
	INVALID_OUTPUT_INTERFACE EnumOutputInterface = iota
    KAFKA 
	ELASTIC_SEARCH
	GRPC
	CIM
	
	FILE
)
const (
	RHSvalueConstStart = "value("
    RHSvalueConstEnd = ")"
)

type EsConfig struct {
    DocID                   []DocID `json:"docID"`
    IndexName               string `json:"indexName"`
    DuplicateIndexName      string `json:"duplicateIndexName"`
    SavetoKafkaFlag         bool `json:"saveToKafka"`
    CdrTimeField            string `json:"cdrTimeField"`
    CdrTimeFieldFormat      string `json:"cdrTimeFieldFormat"`
    CdrTimeGoTimeFmt        string
    IndexPrefixTimeFormat   string `json:"indexPrefixTimeFormat"`
    IndexPrefixGoTimeFmt    string
}

type DocID struct {
	Field string `json:"field"`
}

type KafkaConfig struct {
	TopicName string `json:"topicName"`
}
const (
    LEFT_OUT = "leftOuts"
)

func splitCdr(action SplitCdr, cdrMap map[string]interface{}, outCdrList interface{}, transID string) error {
	// if outInterface != action.MsgInterface.OutInterfaceType {
	// 	fmt.Println("Invalid configuration. Supported OutInterface:", OutputInterfaceConvertEnumToString(outInterface), " Configured OutInterface:" + OutputInterfaceConvertEnumToString(action.MsgInterface.OutInterfaceType))
	// 	return errors.New("Invalid configuration. Supported OutInterface: "+ OutputInterfaceConvertEnumToString(outInterface) + " Configured OutInterface:" + OutputInterfaceConvertEnumToString(action.MsgInterface.OutInterfaceType) )
	// }	
		
	switch(action.MsgInterface.OutInterfaceType){
		case KAFKA:
		case ELASTIC_SEARCH:
		default:
//			fmt.Println("Invalid configuration. OutInterface:", action.MsgInterface.OutInterface)
//			return errors.New("Invalid configuration. OutInterface:" + action.MsgInterface.OutInterface)
	}

	//this map used to maintain segragated data from cdr based on configured and unconfigured array.
	tempMap := make(map[string]map[int]map[string]interface{})

	//this is 1 time loop through cdr elements to segrigate configured and unconfigured arrays/elements
	for cdrKey, cdrValue := range cdrMap {
		//used to find if element is part of configured array or not
		keyFound := false
		//loop through all configured arrays for given cdr element.
		for _, arrayName := range action.ArrayNameList {
			//if cdr element is match to configured array

			//if match, _ := regexp.MatchString(arrayName.ArrayName+"_"+"\\d+_", cdrKey); match {
			if strings.HasPrefix(cdrKey, arrayName.ArrayName) {
				//get the index value from array element
				indexValue := getIndex(cdrKey, arrayName.ArrayName)
				if  indexValue > -1 {
					//indexValue, _ := strconv.Atoi(cdrKey[len(arrayName.ArrayName+"_"):][:strings.Index(cdrKey[len(arrayName.ArrayName+"_"):], "_")])

					//create internal structure of tempMap if not already exists
					if tempMap[arrayName.ArrayName] == nil {
						tempMap[arrayName.ArrayName] = make(map[int]map[string]interface{})
					}
					if tempMap[arrayName.ArrayName][indexValue] == nil {
						tempMap[arrayName.ArrayName][indexValue] = make(map[string]interface{})
					}
					//populate map according.
					tempMap[arrayName.ArrayName][indexValue][cdrKey] = cdrValue
					//set it saying that we found a match.
					keyFound = true
				}
			}
		}
		//for given element after going through all the config arrays no match found, then consider it as left out and maintain it seprately.
		if !keyFound {
			if tempMap[LEFT_OUT] == nil {
				tempMap[LEFT_OUT] = make(map[int]map[string]interface{})
				tempMap[LEFT_OUT][0] = make(map[string]interface{})
			}
			tempMap[LEFT_OUT][0][cdrKey] = cdrValue
		}
	}

	//just print for reference sake.
	fmt.Println("Split map post decoding:", tempMap)

	//these variables help to create array of combined new cdrs.
	outPutMap := make(map[int]interface{})
	outPutCounter := 0

	//this is anonymous will do recursion and recreate output new cdrs.
	var combineFunc func(configIndex int, newCdr map[string]interface{})

	combineFunc = func(configIndex int, newCdr map[string]interface{}) {
		if newCdr == nil {
			newCdr = make(map[string]interface{})
		}
		if configIndex < len(action.ArrayNameList) {
			for i := 0; i < len(tempMap[action.ArrayNameList[configIndex].ArrayName]); i++ {
				newCdr1 := make(map[string]interface{})
				for k1, v1 := range newCdr {
					newCdr1[k1] = v1
				}
				for k, v := range tempMap[action.ArrayNameList[configIndex].ArrayName][i] {
					key := formKey(k, action.ArrayNameList[configIndex].ArrayName)
					if key != "" {
						newCdr1[key] = v
					}
					//using regex removing index information from new cdr
					/*
					m := regexp.MustCompile("^" + action.ArrayNameList[configIndex].ArrayName + "_(\\d*)_(.*)$")
					newCdr1[m.ReplaceAllString(k, action.ArrayNameList[configIndex].ArrayName+"_$2")] = v
					*/
				}
				combineFunc(configIndex+1, newCdr1)
			}
		} else {
			outPutMap[outPutCounter] = newCdr
			outPutCounter++
		}
	}

	//calling anonymous func to start recreating output cdrs.
	combineFunc(0, nil)

	for key, element := range outPutMap {
		splitcdr := element.(map[string]interface{})
		for k, v := range tempMap[LEFT_OUT][0] {
				splitcdr[k] = v
		}
		if retVal := EvaluateIfCondition(action.Condition, splitcdr, transID); !retVal {
			delete(outPutMap, key)
		}
	}

	msize := len(outPutMap);
	counter := 1;
	//this loop will help in creating final output.

	keys := make([]int, 0, len(outPutMap))
	for k := range outPutMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	if action.ChangePrevOccuranceTxnTypeToUpdateOnSameRGOnRelease {
		if cdrMap["recordExtensions_transactionType"] == "RELEASE" {
			rgMap := make(map[string]bool)
			for i := (len(keys)-1); i >= 0; i-- {
				splitcdr := outPutMap[keys[i]].(map[string]interface{})
				if len(action.TransactionTypeChangeKeyFields) > 0 {
					key := formKeyUsingCdrFields(splitcdr, action.TransactionTypeChangeKeyFields, transID)
					if _, rgPresent := rgMap[key]; rgPresent {
						//RG Combination present. Change transaction type to UPDATE
						splitcdr["recordExtensions_transactionType"] = "UPDATE"
						splitcdr["listOfMultipleUnitUsage_usedUnitContainer_0_triggers_0_triggerType"] = "VALIDITY_TIME"
					} else {
						//RG not present. Insert
						rgMap[key] = true
					}
				}
			}
		}
	}
	splitCounter := 0
	for i, k := range keys {
		splitcdr := outPutMap[k].(map[string]interface{})

		enrichSplitFlag := false
		if premiumType, ok := splitcdr["recordExtensions_chargeInformation_debitCash_rateProfile_0_premiumRateType"].(string); ok {
			if premiumType == "Connection" || premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "NUC-CC" {
				// do nothing
			} else {
				enrichSplitFlag = true
			}
		} else {
			enrichSplitFlag = true
		}

		if enrichSplitFlag {
			if msize == 1 {
				splitcdr["firstSplitCdr"] = true
				splitcdr["lastSplitCdr"] = true
			} else {
				if counter == 1 {
					splitcdr["firstSplitCdr"] = true
				} else {
					splitcdr["firstSplitCdr"] = false
				}
				if i+1 == msize {
					splitcdr["lastSplitCdr"] = true
				} else {
					// next is the last splitcdr
					flag := checkForNormalRateIfExistAfterPRP(outPutMap, i, msize, keys)
					if flag {
						splitcdr["lastSplitCdr"] = false
					} else {
						splitcdr["lastSplitCdr"] = true
					}
				}
			}
	
			counter++
		}

		if action.ChangeSctName {
			if val, ok := splitcdr[action.SctLookupKey].(string); ok {
				if sctname, ok := action.SctLookup[val]; ok {
					splitcdr["SCTNAME"] = sctname
				} else {
					fmt.Println("PremiumRateType : ", val , " is not configured")
				}
			} else {
				fmt.Println("Field : ",action.SctLookupKey , " is not present in the splitted cdr")
			}
		}

		switch(action.MsgInterface.OutInterfaceType) {
			case KAFKA:
				//Error log Handled inside the below function. Continue with next split CDR
				//gpp.SavetoKafkaInterface(action.MsgInterface.Kafka, splitcdr, outCdrList, transID)
				fmt.Println("Saving cdr to kafka", splitcdr)
				encodedData, _ := json.Marshal(splitcdr)
				ioutil.WriteFile("splitCDROutput"+strconv.Itoa(splitCounter), encodedData, 0644)
				splitCounter++
			case ELASTIC_SEARCH:
				// docID := gpp.DeriveDocumentID(action.MsgInterface.Es.DocID, splitcdr)
				// if len(docID) > 0 {
				// 	splitcdr["documentID"] = docID
				// }
				// gpp.SavetoEsInterface(action.MsgInterface.Es, splitcdr, outCdrList, docID, transID)
		}
		// for i, sct := range action.SctList {
		// 	splitcdr["SCTNAME"] = sct 
		// 	newTransId := splitcdr["transID"].(string) + "_s" + strconv.Itoa(i+1)
		// 	splitcdr["transID"] = newTransId
		// 	err := loaderInvokeSCT(splitcdr, outCdrList, sct, newTransId)
		// 	if err != nil {
		// 		fmt.Println("Failed to execute the services for SCT ", sct)
		// 	}			
		// }
	}
	return nil
}

func GetPremiumType(chargeInfo map[string]interface{}, transID string) (premiumType string) {
	defer PanicHandler("getPremiumType", transID)
	premiumType, _ = chargeInfo["debitCash"].(map[string]interface{})["rateProfile"].([]interface{})[0].(map[string]interface{})["premiumRateType"].(string)
	return premiumType
}

func checkForNormalRateIfExistAfterPRP(outPutMap map[int]interface{}, i, msize int, keys []int) bool {
	if i == msize-1 {
		return false
	}
	nextSplitcdr := outPutMap[keys[i+1]].(map[string]interface{})
	if premiumType, ok := nextSplitcdr["recordExtensions_chargeInformation_debitCash_rateProfile_0_premiumRateType"].(string); ok {
		if premiumType == "Connection" || premiumType == "FUP" || premiumType == "NetworkUsage" || premiumType == "NUC-CC" {
			return checkForNormalRateIfExistAfterPRP(outPutMap, i+1, msize, keys)
		} else {
			return true
		}
	}
	return true
}

func loaderInvokeSCT(cdr map[string]interface{}, outCdrList interface{}, sctName ,transID string) error {
	// err := gpp.ExecuteServices(cdr, outCdrList, sctName, transID)
	// if err != nil {
	// 	fmt.Println("Failed to execute the services")
	// 	return errors.New("Failed to execute the services")
	// }
	// return nil
	fmt.Println("Invoking GPP for SCT", sctName)
	return nil
}

func formKey(cdrKey , cdrPrefix string) string {
    cdrPrefixLen := len(cdrPrefix) + 1 // Need to skip _
	cdrKeyLen := len(cdrKey)
	if cdrPrefixLen > cdrKeyLen {
		return ""
	}

    index := strings.Index(cdrKey[cdrPrefixLen:], "_")
	if index == -1 {
		return ""
	}

    s := []rune(cdrKey)
	index += cdrPrefixLen + 1
	res := append(s[0:cdrPrefixLen], s[index:]...)
	return string(res)
}

func formKeyUsingCdrFields (cdr map[string]interface{}, keyArr []string, transID string) (formedKey string) {
	var dstValue string
	for _, field := range keyArr {
		if cdr[field] != nil {
			switch cdr[field].(type) {
			case string:
				dstValue = cdr[field].(string)
			case float64:
				floatVal := cdr[field].(float64)
				dstValue = strconv.FormatFloat(floatVal, 'f', -1, 64)
			case int:
				dstValue = strconv.Itoa(cdr[field].(int))
			case bool:
				dstValue = strconv.FormatBool(cdr[field].(bool))
			default:
				fmt.Println("Unsupported source field")
			}
			formedKey += dstValue
		}
	}
	return formedKey
}

func getIndex(cdrKey , cdrPrefix string) int {
    cdrPrefixLen := len(cdrPrefix) + 1 // Need to skip _
	cdrKeyLen := len(cdrKey)
	if cdrPrefixLen > cdrKeyLen {
		return -1
	}

	index := strings.Index(cdrKey[cdrPrefixLen:], "_")
	if index == -1 {
		return index
	}

	indexValue, err := strconv.Atoi(cdrKey[cdrPrefixLen:][:index])
	if err != nil {
		return -1
	}
	return indexValue
}

func EvaluateIfCondition(ifConditionStmt string, cdr map[string]interface{}, transID string) (bool) {

    if ifConditionStmt == "" {
        return true
    }

    condition := ifConditionStmt

    for {
        startIndex := strings.Index(condition, RHSvalueConstStart)
        if startIndex == -1 {
            break
        }

        field := condition[startIndex+len(RHSvalueConstStart) : startIndex+strings.Index(condition[startIndex:], RHSvalueConstEnd)]

        fieldValue := ""

        if _, ok := cdr[field]; !ok {
            fmt.Println("configured field:", field, "not available in cdr")
            return false
        }

        switch cdr[field].(type) {
            case int64:
                    fieldValue = strconv.FormatInt(cdr[field].(int64), 10)
            case int:
                fieldValue = strconv.Itoa(cdr[field].(int))
            case float64:
                fieldValue = fmt.Sprintf("%f", cdr[field].(float64))
            case string:
                fieldValue = "\"" + cdr[field].(string) + "\""
            case bool:
                if cdr[field].(bool) {
                    fieldValue = "true"
                } else {
                    fieldValue = "false"
                }
            default:
                    fmt.Println("unsupported field type in cdr for field:", field)
                    return false
        }
        condition = condition[0:startIndex] + fieldValue + condition[strings.Index(condition, RHSvalueConstEnd)+1:]
    }
    //check result from rule engine
    result := InvokeRuleEngine(condition, cdr, transID)

    fmt.Println("Result after executing query:", condition, " Result:" , result)
    return result
}

func InvokeRuleEngine(rule string, Cdr map[string]interface{}, transID string) bool {

	//	fmt.Println("CustomWhen", rule)
		evRe, err := parser.NewEvaluator(rule)
		if err != nil {
			fmt.Println("CustomWhen:",rule , "Error in evalating Rule, error:", err)
			return false
		}
		resRe, errRe := evRe.Process(Cdr)
		if errRe != nil {
			fmt.Println("CustomWhen:",rule , "Error in procesing the Rule, error:", errRe)
			return false
		}
		fmt.Println("CustomWhen:",rule , "RE Resp:", resRe)
	
		return resRe
	}

type Evaluator struct {
	lastDebugErr error

	rule string
	tree antlr.ParseTree

	testHookPanic func()
}


func PanicHandler(funcName, transID string) {
	if err := recover(); err != nil {
		fmt.Println("An Error occured in function --", funcName, " error is --", err)
	}
}