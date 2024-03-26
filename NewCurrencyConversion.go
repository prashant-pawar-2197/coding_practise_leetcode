package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/shopspring/decimal"
)


type CurrencyConversionParam struct {
	Precision      int    `json:"precision"`
	DivisionFactor int    `json:"divisionFactor"`
	Operation      string `json:"operation"`
}

type CurrencyConversion struct {
	SourceFields      []string                           `json:"sourceFields"`
	DestinationFields []string                           `json:"destinationFields"`
	Key               []string                            `json:"key"`
	CurrencyLookup    map[string]CurrencyConversionParam `json:"currencyLookup"`
	Data              []struct {
		Key   string                  `json:"key"`
		Value CurrencyConversionParam `json:"value"`
	} `json:"data,omitempty"`
	DefaultPrecision	bool	`json:"defaultPrecision"`
}

func (s *CurrencyConversion) UnmarshalJSON(data []byte) error {
    type CurrencyLookup CurrencyConversion
    if err := json.Unmarshal(data, (*CurrencyLookup)(s)); err != nil {
        return err
    }
	if  s.CurrencyLookup == nil {
		s.CurrencyLookup = make(map[string]CurrencyConversionParam)
	}
	for j := 0; j < len(s.Data); j++ {
		s.CurrencyLookup[s.Data[j].Key] = s.Data[j].Value
	}
    return nil
}

func main() {
	var action CurrencyConversion
	var cdr map[string]interface{}
	readFileData, err := ioutil.ReadFile("sampleAction.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	//fmt.Println(action)
	readCDRData, err := ioutil.ReadFile("monetaryCDR.json")
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
	//fmt.Println(action.CurrencyLookup["Default"].Precision)
	//fmt.Println(key)

	var ( 
		ratio float64
        key string
        currencyConvDetails CurrencyConversionParam 
        currencyConvDetailsfnd bool
	)
    delimiter := "_"
        for _, v := range action.Key {
           key += cdr[v].(string)
           key += delimiter
        }
        if len(key) > 0{
			key = key[:len(key)-1]
		} 
        if len(key) == 0 { 
            if currencyConvDetails, currencyConvDetailsfnd = action.CurrencyLookup["Default"]; !currencyConvDetailsfnd {
               return 
			}
        }else {
            if currencyConvDetails , currencyConvDetailsfnd = action.CurrencyLookup[key]; !currencyConvDetailsfnd{
                return 
			}
        }
        ratio = math.Pow(10, float64(currencyConvDetails.Precision))
        for i := 0; i < len(action.SourceFields); i++ {
            field , ok := cdr[action.SourceFields[i]].(float64)
            if !ok {
                continue               
            }
            field = field/float64(currencyConvDetails.DivisionFactor)
            switch currencyConvDetails.Operation {
                case "ROUND":
                        cdr[action.DestinationFields[i]] = math.Round(field*ratio)/ratio
                case "CEIL":
                        cdr[action.DestinationFields[i]] = math.Ceil(field*ratio)/ratio
                case "FLOOR":
                        destValue, err := strconv.ParseFloat((decimal.NewFromFloat(field).RoundFloor(int32(currencyConvDetails.Precision)).String()), 64)
                        if err != nil {
                            continue
                        }
                        cdr[action.DestinationFields[i]] = destValue
            }  
        }      
		encodedData,_ := json.Marshal(cdr)
		ioutil.WriteFile("monetaryOutputData",encodedData, 0644)
}
