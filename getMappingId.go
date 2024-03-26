package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func actionGetMappingID(action GetMappingId, cdr map[string]interface{}, transID string) error {
	var fieldValue, spfield string
	fieldAvailable := false
	foundCoId := false
	//mandatoryField := false
	//fetching the service type from the action configured
	ServType := action.ServiceType

	if (strings.EqualFold(ServType, "sms")) {
		//in case of sms type service
		for i, value := range action.Source_1 {
			if i == 0 {
				spfield, fieldAvailable = cdr[value].(string)
				if (!fieldAvailable || spfield == "") {
					break
				}
				fieldValue = spfield + fieldValue
				foundCoId = true
				continue
			} else {
				spfield, fieldAvailable = cdr[value].(string)
				if (!fieldAvailable || spfield == "") {
					fmt.Println("field ", value , " not found in cdr")
					return errors.New("field "+ value + " not found in cdr")
				}
				fLen := len(spfield)
				if fLen > 5 {
					fieldValue = fieldValue + spfield[fLen-5:]
				} else {
					fieldValue = fieldValue + spfield
				}
			}
		}

		//if ccCorelation id is not configured in action
		if !foundCoId {
			for _, value := range action.Source_2 {
				spfield, fieldAvailable = cdr[value].(string)
				if !fieldAvailable || spfield == "" {
					fmt.Println("field ", value , " not found in cdr")
					return errors.New("field "+ value + " not found in cdr")
				}
				fieldValue = spfield
			}
		}
	} else if ((strings.EqualFold(ServType, "mms")) || (strings.EqualFold(ServType, "voice")) ||
		(strings.EqualFold(ServType, "data"))) {
		//in other type of service
		for _, value := range action.Source_1 {
			spfield, fieldAvailable = cdr[value].(string)
			if !fieldAvailable || spfield == "" {
				fmt.Println("field ", value , " not found in cdr")
				return errors.New("field "+ value + " not found in cdr")
			}
			fieldValue = spfield
		}
	} else if (strings.EqualFold(ServType, "edr")) {
		//in edr type of service
		for _, value := range action.Source_1 {
			spfield, fieldAvailable = cdr[value].(string)
			if !fieldAvailable || spfield == "" {
				break
			}
			fieldValue = spfield
		}

		if fieldValue == "" {
			for _, value := range action.Source_2 {
				spfield, fieldAvailable = cdr[value].(string)
				if !fieldAvailable || spfield == "" {
					fmt.Println("field ", value , " not found in cdr")
					return errors.New("field "+ value + " not found in cdr")
				}
				fieldValue = spfield
			}
		}
	} else {
		//invalid service type
		fmt.Println("Unsupported service:", ServType, " configured.")
		return errors.New("Unsupported service:" + ServType + " configured.")
	}

	cdr[action.Destination[1]] = fieldValue

	var idHashValue string
	if len(fieldValue) > 32 {
		//hashing the value before populating since the length is greater than 32
		byteFieldValue := []byte(fieldValue)
		hashData := md5.Sum(byteFieldValue)
		idHashValue = hex.EncodeToString(hashData[:])
	} else {
		// no need of hashing the value since the length is less than or equal to 32
		idHashValue = fieldValue
	}

	cdr[action.Destination[0]] = idHashValue
	fmt.Println( "input network session id is generated as : ", cdr[action.Destination[1]], " output network session id is generated as : ", cdr[action.Destination[0]])

	return nil
}

type GetMappingId struct {
	ServiceType  string   `json:"serviceType"`
	Source_1     []string   `json:"source_1,omitempty"`
	Source_2     []string   `json:"source_2,omitempty"`
	Destination  []string   `json:"destination"`
	SaveFlag     bool    `json:"saveFlag"`
}

func main() {
	cdr := make(map[string]interface{})
	var action GetMappingId
	readFileData, err := ioutil.ReadFile("getMappingSdr.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	readActionData, err := ioutil.ReadFile("getMappingsdrAction.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readActionData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	actionGetMappingID(action,cdr,"123")
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("finalOutputPackageSplitCDR",encodedData, 0644)
}