package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)
type CustomSet struct {
	SetName 	   string `json:"SetName"`
	Prefixes []struct{
		Prefix []string `json:"Prefix"`
		Value string `json:"Value"`
	} `json:"Prefixes"`
}
func main() {
	const fileName = "C://Users//pawarpr//OneDrive - Mavenir Systems, Inc//Documents//GoPractise and notes//sampleCSV.csv"
	var cs  = CustomSet {Prefixes :[]struct{Prefix []string "json:\"Prefix\""; Value string "json:\"Value\""}{{Prefix: make([]string,0,20) , Value: "int mobile" },{Prefix: make([]string,0,20) , Value: "int premium" },{Prefix: make([]string,0,20) , Value: "int fixed" },{Prefix: make([]string,0,20) , Value: "int satellite" },{Prefix: make([]string,0,20) , Value: "int service" }}} 
	cs.SetName = "someLookup"
	re, _ := regexp.Compile("(^[0]+)([0-9]+)")
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if record[2] == "int mobile"{
			numWithoutLeadingZeros := re.ReplaceAllString(record[1], "$2")
			cs.Prefixes[0].Prefix = append(cs.Prefixes[0].Prefix,numWithoutLeadingZeros)
		}
		if record[2] == "int premium"{
			numWithoutLeadingZeros := re.ReplaceAllString(record[1], "$2")
			cs.Prefixes[1].Prefix = append(cs.Prefixes[1].Prefix, numWithoutLeadingZeros)
		}
		if record[2] == "int fixed"{
			numWithoutLeadingZeros := re.ReplaceAllString(record[1], "$2")
			cs.Prefixes[2].Prefix = append(cs.Prefixes[2].Prefix, numWithoutLeadingZeros)
		}
		if record[2] == "int satellite"{
			numWithoutLeadingZeros := re.ReplaceAllString(record[1], "$2")
			cs.Prefixes[3].Prefix = append(cs.Prefixes[3].Prefix, numWithoutLeadingZeros)
		}
		if record[2] == "int service"{
			numWithoutLeadingZeros := re.ReplaceAllString(record[1], "$2")
			cs.Prefixes[4].Prefix = append(cs.Prefixes[4].Prefix, numWithoutLeadingZeros)
		}
	}
	u, err := json.Marshal(cs)
	if err != nil {
		panic(err)
	}
	errrr := ioutil.WriteFile("./newProcessedCSV.json",[]byte(u),0644)
        if errrr !=nil{
    	     fmt.Println("Couldn't write the data")
        }
}