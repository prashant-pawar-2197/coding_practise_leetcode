package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	cdr := make(map[string]interface{})
	readFileData, err := ioutil.ReadFile("cacheSampleData.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	fmt.Println(cdr["TIMER_ID_FIELD"].(map[string]interface{})["short_timer"])
}