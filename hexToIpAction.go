package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

type HexToIpConversion struct {
	Source string `json:"source"`
	IpType string `json:"ipType"`
	Destination string `json:"destination"`
}

func actionHexToIpConversion(cdr map[string]interface{}, action HexToIpConversion) error {
	if action.Source == "" {
		fmt.Println("Source field is not configured")
		return errors.New("Source field is not configured")
	}
	if action.Destination == "" {
		fmt.Println("Destination field is not configured")
		return errors.New("Destination field is not configured")
	}
	if action.IpType == "" {
		fmt.Println("IpType field is not configured")
		return errors.New("IpType field is not configured")
	}
	var (
		fieldVal string
		fieldAvailable bool
	)
	switch action.IpType {
	case "IPV4":
		if fieldVal, fieldAvailable = cdr[action.Source].(string) ; !fieldAvailable || fieldVal == "" {
			fmt.Println("Source field ", action.Source, " is not present or is empty")
			return errors.New("Source field " + action.Source + " is not present or is empty")
		} else {
			num, err := strconv.ParseUint(fieldVal, 16, 32)
			if err != nil {
				fmt.Println("Error occured while parsing hex value :", fieldVal)
				return errors.New("Error occured while parsing hex value :" + fieldVal)			
			}
			// Convert the integer to an IP address
			ip := net.IPv4(byte(num>>24), byte(num>>16), byte(num>>8), byte(num))
			cdr[action.Destination] = ip.String()
			fmt.Println(cdr[action.Destination])
		}
	default: 
		fmt.Println("Unsupported IP Type : ", action.IpType, " configured.")
		return errors.New("Unsupported IP Type : " + action.IpType + " configured.")
	}
	return nil
}

func main() {
	var action HexToIpConversion
	var cdr map[string]interface{}
	readFileData, err := ioutil.ReadFile("sampleActionHexToIp.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	//fmt.Println(action)
	readCDRData, err := ioutil.ReadFile("pgwCDRHextoIp.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	err = json.Unmarshal(readCDRData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	actionHexToIpConversion(cdr, action)
}
