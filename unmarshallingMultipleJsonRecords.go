package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

func sortByChargingSessIdentifierAndRecordSeqNumber(cdrMap []map[string]interface{})  {
	
	sort.Slice(cdrMap, func(i, j int) bool {
		var sortedByChargingSessIdentifier, sortedByRecordSequenceNumber bool
		if cdrMap[i]["chargingSessionIdentifier"] == nil || cdrMap[j]["chargingSessionIdentifier"] == nil{
			return true
		}
		sortedByChargingSessIdentifier = cdrMap[i]["chargingSessionIdentifier"].(string) < cdrMap[j]["chargingSessionIdentifier"].(string)

        // sort by lowest recordSequenceNumber
        if cdrMap[i]["chargingSessionIdentifier"] == cdrMap[j]["chargingSessionIdentifier"] {
            sortedByRecordSequenceNumber = cdrMap[i]["recordSequenceNumber"].(float64) < cdrMap[j]["recordSequenceNumber"].(float64)
            return sortedByRecordSequenceNumber
        }
        return sortedByChargingSessIdentifier
	})
	
}

func mainn(){
	start := time.Now()
	sampleData := make([]map[string]interface{},0,10)
	readFileData, err := ioutil.ReadFile("sampleLogData.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	//count , _ := cl.Count(readFileData)
	//fmt.Println(count)
	payload := strings.NewReader(string(readFileData))
		fscanner := bufio.NewScanner(payload)
		for fscanner.Scan(){
			var mapp map[string]interface{}
			err := json.Unmarshal(fscanner.Bytes(), &mapp)
			if err != nil {
				fmt.Printf("%s", err)
				continue
			}
			sampleData = append(sampleData,mapp)
		}

		//m := make(map[string]interface{})
		//m = sampleData[0]["header"].(map[string]interface{})
		fmt.Println( reflect.TypeOf(sampleData[0]["header"].(map[string]interface{})["eventTime"]))
		//sortByChargingSessIdentifierAndRecordSeqNumber(sampleData)

		file,err := os.OpenFile("ProcessedCDR.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		for i := 0; i < len(sampleData); i++ {
			encodedData,_ := json.Marshal(sampleData[i])
		if err != nil {
			  fmt.Println("Could not open ProcessedCDR")
		  return
		}
		_, err2 := file.WriteString(string(encodedData)+"\n")
		  if err2 != nil {
			  fmt.Println("Could not write text to example.txt")
		  }else{
		  fmt.Println("Operation successful! Text has been appended to example.txt")
		}
		defer file.Close()	
	}
	elapsed := time.Since(start)
    fmt.Printf("Unmarshalling + Marshalling + sorting is taking %s", elapsed)

}


/*
	errrr := ioutil.WriteFile("C:\\Users\\pawarpr\\Documents\\GoPractise and notes\\ProcessedCDR.json",[]byte(encodedData),0644)
	if errrr !=nil{
		fmt.Println("Couldn't write the data")
	}
*/