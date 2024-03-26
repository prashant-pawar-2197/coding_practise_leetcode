package main

import (
	"fmt"
	"strings"
	"time"
)

var CNotationDateTimeSubstution = map[string]string{
	"%YEAR%":      "2006",    //full year
	"%year%":      "06",      //short year
	"%MONTH%":     "January", //full char month
	"%month_nz%":  "01",      //num month with zero
	"%month_n%":   "1",       //num month with out zero
	"%month%":     "Jan",     //short char month
	"%WEEK_DAY%":  "Monday",  //day of week in full char
	"%week_day%":  "Mon",     //day of week in short char
	"%date_z%":    "02",      //date with zero
	"%date%":      "2",       //date with out zero
	"%hour_24%":   "15",      //hours in 24 hour format
	"%hour_12_z%": "03",      //hours in 12 hour format with zero
	"%hour_12%":   "3",       //hours in 12 hour format with out zero
	"%minute_z%":  "04",      //minute with zero
	"%minute%":    "4",       //minute with out zero
	"%second_z%":  "05",      //second with zero
	"%second%":    "5",       //second with out zero
}

func GetDynamicCNotationDateTimeSubstutions(input string) string {
	for key, value := range CNotationDateTimeSubstution {
		input = strings.Replace(input, key, value, -1)
	}
	return input
}

func main() {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Println("err occ")
	}
	actualFormatedTime := GetDynamicCNotationDateTimeSubstutions("%YEAR%%month_nz%%date_z%%hour_24%%minute_z%%second_z%")
	fieldValue := "20080520140000"
	tm, err := time.Parse(actualFormatedTime, fieldValue)
	if err != nil {
		fmt.Println(err)
	}
	tm = tm.In(location)
	fmt.Println(tm.UTC())
}
