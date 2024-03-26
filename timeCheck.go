package main

import (
	"fmt"
	"strings"
	"time"
)


// var CNotationDateTimeSubstution = map[string]string{
// 	"%YEAR%":      "2006",    //full year
// 	"%year%":      "06",      //short year
// 	"%MONTH%":     "January", //full char month
// 	"%month_nz%":  "01",      //num month with zero
// 	"%month_n%":   "1",       //num month with out zero
// 	"%month%":     "Jan",     //short char month
// 	"%WEEK_DAY%":  "Monday",  //day of week in full char
// 	"%week_day%":  "Mon",     //day of week in short char
// 	"%date_z%":    "02",      //date with zero
// 	"%date%":      "2",       //date with out zero
// 	"%hour_24%":   "15",      //hours in 24 hour format
// 	"%hour_12_z%": "03",      //hours in 12 hour format with zero
// 	"%hour_12%":   "3",       //hours in 12 hour format with out zero
// 	"%minute_z%":  "04",      //minute with zero
// 	"%minute%":    "4",       //minute with out zero
// 	"%second_z%":  "05",      //second with zero
// 	"%second%":    "5",       //second with out zero
// }

// func GetDynamicCNotationDateTimeSubstutions(input string) string {
// 	for key, value := range CNotationDateTimeSubstution {
// 		input = strings.Replace(input, key, value, -1)
// 	}
// 	return input
// }


func covertTimeToEpoch(inputTime string) int64 {
	if strings.ContainsAny(inputTime,"-T"){
		timeFormat := "%YEAR%-%month_nz%-%date_z%T%hour_24%:%minute_z%:%second_z%-07:00"
		actualFormatedTime := GetDynamicCNotationDateTimeSubstutions(timeFormat)
		tm, err := time.Parse(actualFormatedTime, inputTime)
		if err != nil {
			fmt.Println("Error occured while parsing the time")
		}
		return tm.Unix()
	}else{
		var charSign string
		var year , month , day , hour , minute , second , sign, zoneh, zonem, final string
		if len(inputTime) > 0 {
			year  = string(inputTime[1]) + string(inputTime[0])
			month = string(inputTime[3]) + string(inputTime[2])
			day   = string(inputTime[5]) + string(inputTime[4])
			hour  = string(inputTime[7]) + string(inputTime[6])
			minute = string(inputTime[9]) + string(inputTime[8])
			second = string(inputTime[11]) + string(inputTime[10])
			sign = string(inputTime[12]) + string(inputTime[13])
			zoneh = string(inputTime[15]) + string(inputTime[14])
			zonem = string(inputTime[17]) + string(inputTime[16])
		}
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
		return t.Unix()
	}
}
func main(){
	fmt.Println(covertTimeToEpoch("2023-03-28T09:19:05+02:00"))
}