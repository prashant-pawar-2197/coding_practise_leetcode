package main

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"regexp"
)

func main() {
	/*
	strArr := []string{ "SAT.1\n", 
						"SAT.1\a", 
						"SAT.1\b", 
						"SAT.1\f",
						"SAT.1\r", 
						"SAT.1\t", 
						"SAT.1\v", 
						"SAT.1", 
						"9999999\n",
						"o2 Team",
						"Whatsapp",
						"name@email.com\n",
				}
				*/
	//pattern  := "(^[a-z0-9A-Z.]+[^\\naf\b])"
	//pattern  := "(^[a-z0-9A-Z.]+[^\bnfrt])" working
	//pattern := "(^[a-z0-9A-Z.@ ]+)([^\bnfrt])?"
	/* 
	input := 	SAT.1\n8\t5
				SAT.185
				SAT\nT.1\n8\t5
				SAT.1\n8\n5
				SAT.1\t\n
				SAT.1\n\n
				SAT.1\n\nabc
				Whatsapp
				o2 Team
				9999999\n
				name@email.com\n

	pattern := "([\b\n\f\r\t]+)"

	*/

	strArr := []string{
				"SAT.1\n8\t5",
				"SAT.185",
				"SAT\nT.1\n8\t5",
				"SAT.1\n8\n5",
				"SAT.1\t\n",
				"SAT.1\n\n",
				"SAT.1\n\nabc",
				"Whatsapp",
				"o2 Team",
				"9999999\n",
				"name@email.com\n",
	}
	pattern := "([\b\n\f\r\t]+)"
	//pattern := (^[a-z0-9A-Z.@ ]+)([^\bnfrt])?"
	//field := "SAT.1\n8\t5"
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Err occured")
	}
	/*
	mnc := re.ReplaceAllString(field, "~")
	fmt.Println("Output -->", mnc)
	*/
	for _, val := range strArr {
		fieldValue := val
		mnc := re.ReplaceAllString(fieldValue, "")
		fmt.Println("Output -->", mnc)
	}
	
}