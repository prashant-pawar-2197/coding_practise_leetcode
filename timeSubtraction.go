package main

import (
	"fmt"
	"math"
	// "strings"
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

// func subtractTime(inputTime, timeFormat string, duration int64) (finalUtcTime string) {
// 	actualFormatedTime := GetDynamicCNotationDateTimeSubstutions(timeFormat)
// 	tm, err := time.Parse(actualFormatedTime, inputTime)
// 	if err != nil {
// 		// ml.MavLog(ml.INFO, transID, "failed to parse time value:", inputTime, " with format:", actualFormatedTime, ",error:", err)
// 		// return errors.New("failed to parse time value:" + inputTime + " with format:" + actualFormatedTime + ",error:" + err.Error())
// 	}
// 	tm = tm.Add(time.Second * time.Duration(-duration))
// 	finalUtcTime = tm.Format(GetDynamicCNotationDateTimeSubstutions(timeFormat))
// 	return finalUtcTime
// }

func SubtractTwoTimeInUtcFormatt(firstTimeStr, secondTimeStr string) (float64, error) {
	firstTime, err := time.Parse(time.RFC3339, firstTimeStr)
	if err != nil {
		return 0, err
	}

	secondTime, err := time.Parse(time.RFC3339, secondTimeStr)
	if err != nil {
		return 0, err
	}

	diff := firstTime.Sub(secondTime).Seconds()
	return math.Floor(diff), nil
}
func main() {
	fmt.Println(SubtractTwoTimeInUtcFormatt("", "2023-09-18T11:55:47Z"))
}

