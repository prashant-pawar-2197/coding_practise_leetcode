package main

import (
	
	"fmt"
	"strings"
	"time"
)

func main(){
	var date string
	var charSign string
    var year , month , day , hour , minute , second , sign, zoneh, zonem, final string 
	date = "3250626045622B0000"
		//	0123456789
	if true {
		year  = string(date[1]) + string(date[0])
		month = string(date[3]) + string(date[2])
		day   = string(date[5]) + string(date[4])
		hour  = string(date[7]) + string(date[6])
		minute = string(date[9]) + string(date[8])
		second = string(date[11]) + string(date[10])
		sign = string(date[12]) + string(date[13])
		zoneh = string(date[15]) + string(date[14])
		zonem = string(date[17]) + string(date[16])
	}
	// The sign is ascii encoded
	if strings.EqualFold(sign,"2B") {
		charSign = "+"
	} else if strings.EqualFold(sign, "2D") {
		charSign = "-"
	} else {
		//ml.MavLog(ml.WARN, transID, "Invalid sign")
		fmt.Println("invalid sign")
	}
	final = "20" + year + "-" + month + "-" + day + "T" + hour + ":" + minute + ":" + second +
		charSign + zoneh + ":" + zonem
	// Parse the date string
	t, err := time.Parse(time.RFC3339Nano, final)
	if err != nil {
	//	ml.MavLog(ml.ERROR, transID, "Failed to parse timestamp")
		fmt.Println("failed to parse timestamp")
	}
	convertedTime:= t.UTC()
	fmt.Println(convertedTime)
	
}