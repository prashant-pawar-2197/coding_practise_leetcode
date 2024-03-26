package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
)

// type CopyableMap   map[string]interface{}
// type CopyableSlice []interface{}

// // DeepCopy will create a deep copy of this map. The depth of this
// // copy is all inclusive. Both maps and slices will be considered when
// // making the copy.
// func (m CopyableMap) DeepCopy() map[string]interface{} {
//     result := map[string]interface{}{}

//     for k,v := range m {
//         // Handle maps
//         mapvalue,isMap := v.(map[string]interface{})
//         if isMap {
//             result[k] = CopyableMap(mapvalue).DeepCopy()
//             continue
//         }
		
//         // Handle slices
//         slicevalue,isSlice := v.([]interface{})
//         if isSlice {
//             result[k] = CopyableSlice(slicevalue).DeepCopy()
//             continue
//         }

//         result[k] = v
//     }

//     return result
// }

// func (s CopyableSlice) DeepCopy() []interface{} {
//     result := []interface{}{}

//     for _,v := range s {
//         // Handle maps
//         mapvalue,isMap := v.(map[string]interface{})
//         if isMap {
//             result = append(result, CopyableMap(mapvalue).DeepCopy())
//             continue
//         }

//         // Handle slices
//         slicevalue,isSlice := v.([]interface{})
//         if isSlice {
//             result = append(result, CopyableSlice(slicevalue).DeepCopy())
//             continue
//         }
//         result = append(result, v)
//     }

//     return result
// }

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

type AccessNetworkInfo struct{
	AccessNetworkInformation string		`json:"accessNetworkInformation"`
	AccessChangeTime int64		`json:"accessChangeTime"`
	Direction string	
}

var AccessNetworkInfoArr []AccessNetworkInfo = make([]AccessNetworkInfo, 0)



func fetchAsbcCdrs(correlationDoc map[string]interface{} ,AsbcCdrsMap map[string]map[string]interface{}, direction string){
	for _, asbcCdr := range correlationDoc["ReferenceDocs"].([]interface{}) {
		val, ok := asbcCdr.(map[string]interface{})["ServiceType"].(string)
		if ok{
			if strings.HasSuffix(val, direction){
				docRes , err := collection.Get(asbcCdr.(map[string]interface{})["DocKey"].(string), nil)
				if err != nil {
					fmt.Println("Error occured while fetching asbcCdr")
				}
				asbcCdr := make(map[string]interface{})
				docRes.Content(&asbcCdr)
				AsbcCdrsMap[direction] = asbcCdr
			}else{
				continue
			}
		}else{
			fmt.Println("Correlation document is not correctly configured")
		}
	}
}

// func covertTimeToEpoch(inputTime string) int64 {
// 	if strings.ContainsAny(inputTime,"-T"){
// 		timeFormat := "%YEAR%-%month_nz%-%date_z%T%hour_24%:%minute_z%:%second_z%Z"
// 		actualFormatedTime := GetDynamicCNotationDateTimeSubstutions(timeFormat)
// 		tm, err := time.Parse(actualFormatedTime, inputTime)
// 		if err != nil {
// 			fmt.Println("Error occured while parsing the time")
// 		}
// 		return tm.Unix()
// 	}else{
// 		var charSign string
// 		var year , month , day , hour , minute , second , sign, zoneh, zonem, final string
// 		if len(inputTime) > 0 {
// 			year  = string(inputTime[1]) + string(inputTime[0])
// 			month = string(inputTime[3]) + string(inputTime[2])
// 			day   = string(inputTime[5]) + string(inputTime[4])
// 			hour  = string(inputTime[7]) + string(inputTime[6])
// 			minute = string(inputTime[9]) + string(inputTime[8])
// 			second = string(inputTime[11]) + string(inputTime[10])
// 			sign = string(inputTime[12]) + string(inputTime[13])
// 			zoneh = string(inputTime[15]) + string(inputTime[14])
// 			zonem = string(inputTime[17]) + string(inputTime[16])
// 		}
// 		if strings.EqualFold(sign,"2B") {
// 			charSign = "+"
// 		} else if strings.EqualFold(sign, "2D") {
// 			charSign = "-"
// 		} else {
// 			//ml.MavLog(ml.WARN, transID, "Invalid sign")
// 			fmt.Println("invalid sign")
// 		}
// 		final = "20" + year + "-" + month + "-" + day + "T" + hour + ":" + minute + ":" + second +
// 			charSign + zoneh + ":" + zonem
// 		// Parse the date string
// 		t, err := time.Parse(time.RFC3339Nano, final)
// 		if err != nil {
// 		//	ml.MavLog(ml.ERROR, transID, "Failed to parse timestamp")
// 			fmt.Println("failed to parse timestamp")
// 		}
// 		return t.Unix()
// 	}
// }

func covertTimeToUtc(inputTime int64) string {
	UtcTime := time.Unix(inputTime, 0)
	location,err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Println("Error occured while determining the location")
	}
	UtcTime = UtcTime.In(location)
	return UtcTime.Format(GetDynamicCNotationDateTimeSubstutions("%YEAR%-%month_nz%-%date_z%T%hour_24%:%minute_z%:%second_z%-07:00"))
}

func delete_empty (s []AccessNetworkInfo) []AccessNetworkInfo {
	var r []AccessNetworkInfo
	for _, networkInfo := range s {
		if !reflect.ValueOf(networkInfo).IsZero() {
			r = append(r, networkInfo)
		}
	}
	return r
}

func sortBasedOnTime(){ 
	fmt.Println(AccessNetworkInfoArr)
	sort.SliceStable(AccessNetworkInfoArr, func(i, j int) bool {
		return AccessNetworkInfoArr[i].AccessChangeTime < AccessNetworkInfoArr[j].AccessChangeTime
	})
	fmt.Println(AccessNetworkInfoArr)
}

var (
	collection *gocb.Collection
	Bucket     *gocb.Bucket
)
func SplitRatedCDRAccToASBC(correlationDoc map[string]interface{}){
	// map to store ASBC CDR
	//var AsbcCdrsMap map[string]map[string]interface{} =  make(map[string]map[string]interface{})
	commitDocs := correlationDoc["CommitDocs"].([]interface{})
	if len(commitDocs) > 0{
		//traverse each RatedCDR and split based on ASBC
		for _, doc := range commitDocs {
			docRes , err := collection.Get(doc.(map[string]interface{})["DocKey"].(string), nil)
			if err != nil {
				fmt.Println("Error occured while fetching the cdr")
				continue
			}
			var cdr CopyableMap
			var outCdrList []map[string]interface{} = make([]map[string]interface{}, 0)
			docRes.Content(&cdr)
			cdr["startTimeGerman"] = covertTimeToEpoch(cdr["startTimeGerman"].(string))
			cdr["endTimeGerman"] = covertTimeToEpoch(cdr["endTimeGerman"].(string))
			//cdrsMap[doc.(map[string]interface{})["ServiceType"].(string)] = cdr
			// if cdr["recordExtensions_extChargingInformation_networkFlag"] == "NF_OFFNET"{
			// 	if strings.HasSuffix(doc.(map[string]interface{})["ServiceType"].(string), "O") {
			// 		fetchAsbcCdrs(correlationDoc, AsbcCdrsMap,"O")
			// 	}else if strings.HasSuffix(doc.(map[string]interface{})["ServiceType"].(string), "T"){
			// 		fetchAsbcCdrs(correlationDoc, AsbcCdrsMap,"T")
			// 	}
			// }else{
				for _, asbcDoc := range correlationDoc["ReferenceDocs"].([]interface{})  {
					docRes , err := collection.Get(asbcDoc.(map[string]interface{})["DocKey"].(string), nil)
					if err != nil {
						fmt.Println("Error occured while fetching asbcCdr")
					}
					asbcCdr := make(map[string]interface{})
					docRes.Content(&asbcCdr)
					counter := 0	
					var ok bool
					for{
						var networkInfo AccessNetworkInfo
						plmn := strings.Replace("pCSCFRecord_list-Of-AccessNetworkInfoChange_$i_accessNetworkInformation", "$i", strconv.Itoa(counter), -1)
						networkInfo.AccessNetworkInformation , ok = asbcCdr[plmn].(string)
						if ok {
							timeStamp := strings.Replace("pCSCFRecord_list-Of-AccessNetworkInfoChange_$i_accessChangeTime", "$i", strconv.Itoa(counter), -1)
							AccessChangeTime , ok := asbcCdr[timeStamp].(string)
							if !ok {
								fmt.Println("AccessChangeTime is absent")
								break
							}
							networkInfo.AccessChangeTime = covertTimeToEpoch(AccessChangeTime)
							if networkInfo.AccessChangeTime < cdr["startTimeGerman"].(int64){
								continue
							}
						}else{
							fmt.Println("AccessNetworkInformation is absent")
							break
						}
						if asbcCdr["pCSCFRecord_role-of-Node"] == "originating"{
							networkInfo.Direction = "O"
						}else{
							networkInfo.Direction = "T"
						}
						AccessNetworkInfoArr = append(AccessNetworkInfoArr, networkInfo)
						counter++
					}
				}
//			}
			sortBasedOnTime()
			newAccessNetworkInfoArr := delete_empty(AccessNetworkInfoArr)
			lenOfAccessNetworkInfoArr := len(newAccessNetworkInfoArr) 
			var i int = 0
			originalCDR := cdr.DeepCopy()
			for i < lenOfAccessNetworkInfoArr {
				if _, ok := cdr["startTimeGerman"] ; ok{
					// if reflect.TypeOf(val).Kind() == reflect.String {
					// 	startTime =  covertTimeToEpoch(val.(string))
					// }else {
					// 	startTime = val.(int64)
					// }
					//startTime = cdrStartTime
					endTime := newAccessNetworkInfoArr[i].AccessChangeTime
					cdr["duration"] = endTime-cdr["startTimeGerman"].(int64)
					cdr["startTimeGerman"] = cdr["startTimeGerman"].(int64)
					cdr["endTimeGerman"] = endTime
					if i == 0 {
						cdr["partialIndicator"] = "F"
					}else{
						cdr["partialIndicator"] = "I"
					}
					spilttedCDR := cdr.DeepCopy()
					outCdrList = append(outCdrList, spilttedCDR)
					cdr["startTimeGerman"] = endTime
					if strings.HasSuffix(doc.(map[string]interface{})["ServiceType"].(string), "O"){
						if newAccessNetworkInfoArr[i].Direction == "O"{
							cdr["networkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
						}else if newAccessNetworkInfoArr[i].Direction == "T"{
							cdr["secondPartyNetworkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
						}
					}else{
						if newAccessNetworkInfoArr[i].Direction == "T"{
							cdr["networkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
						}else if newAccessNetworkInfoArr[i].Direction == "O"{
							cdr["secondPartyNetworkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
						}
					}
				}else{
					fmt.Println("StartTime is missing in the CDR")
					// Should We exit if startTime is not present in the RATED CDR
				}
				i++
			}
			sessionEndtime := originalCDR["endTimeGerman"].(int64)
			if cdr["startTimeGerman"].(int64) < sessionEndtime {
				startTime :=   cdr["startTimeGerman"].(int64)
				duration := sessionEndtime-startTime
				cdr["duration"] = duration
				cdr["partialIndicator"] = "L"
				cdr["startTimeGerman"] = startTime
				cdr["endTimeGerman"] = sessionEndtime
				spilttedCDR := cdr.DeepCopy()
				outCdrList = append(outCdrList, spilttedCDR)
			}
			// originalCDR["startTimeGerman"] = covertTimeToEpoch(originalCDR["startTimeGerman"].(string))
			// originalCDR["endTimeGerman"] = covertTimeToEpoch(originalCDR["endTimeGerman"].(string))
			outCdrList = append(outCdrList, originalCDR)
			for _, cdr := range outCdrList {
				if cdr == nil{
					continue
				}
				cdr["startTimeGerman"] =  covertTimeToUtc(cdr["startTimeGerman"].(int64))
				cdr["endTimeGerman"] = covertTimeToUtc(cdr["endTimeGerman"].(int64))
				// write to kafka after this...
			}
			encodedData,_ := json.Marshal(outCdrList)
			ioutil.WriteFile("isbcCodecOutput",encodedData, 0644)
		}
	}else{
		fmt.Println("flush T/S cdr")
	}
}

func main() {
	opts := gocb.ClusterOptions{Username: "root", Password: "mavenir"}
	cluster, err := gocb.Connect("127.0.0.1:8091", opts)
	if err != nil {
		fmt.Println("error in connection: ", err)
	}
	if err == nil {
		fmt.Println("connection success ")
	}
	Bucket = cluster.Bucket("sessiondb_mediation")

	//Wait upto timout and checks database conn status
	if err = Bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		fmt.Println("Couchbase bootstrapError : %v", err)
	}
	col := Bucket.DefaultCollection()
	collection = col

	//to store T/S cdrs
//	var cdrsMap map[string]map[string]interface{} = make(map[string]map[string]interface{})

	//fetch Correlation document
	existRes, err := collection.Exists("correlationKey", nil)
	if err != nil {
		fmt.Println("error occured while checking of doc exists for key")
	}
	if existRes.Exists(){
		docRes, err := collection.Get("correlationKey", nil)
		if err != nil {
			fmt.Println("error occured while fetching document for key")
			return
		}
		correlationDoc := make(map[string]interface{})
		docRes.Content(&correlationDoc)
		SplitRatedCDRAccToASBC(correlationDoc)
	}
}