package customSet

//get longest prefix match for given string
//func GetMatch(setName, str string) (error, bool, *CustomSetStruct) {
//delete prefix data
//func DeletePrefix(setName, prefix string) (error, bool) {
//update prefix data
//func UpdatePrefixData(setName, prefix string, val *CustomSetStruct) error {
//get prefix data field Version
//func GetPrefixDataFieldVersion(setName, prefix string) (error, float64) {
//get prefix data field SecondParamFrom
//func GetPrefixDataFieldSecondParamFrom(setName, prefix string) (error, string) {
//get prefix data field SecondParamTo
//func GetPrefixDataFieldSecondParamTo(setName, prefix string) (error, string) {
//get prefix data field desc
//func GetPrefixDataFieldDesc(setName, prefix string) (error, string) {
//get prefix data
//func GetPrefixData(setName, prefix string) (error, *CustomSetStruct) {
//checks if prefix exists.
//func PrefixExists(setName, prefix string) (error, bool) {
//Init the custom set cache
//func Init(config string) error {
//add prefix data
//func AddPrefixData(setName, prefix string, val *CustomSetStruct) error {
//checks if Custom set exists.
//func CustomSetExists(setName string) bool {

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	ml "mwp-appcommon/mavnatslogger"

	cu "common.com/commonUtil"
	cb "common.com/couchbase"
	radix "github.com/armon/go-radix"
)

//data structure which maintains all configured radix data with custom set
var CustomSetCache = make(map[string]*radix.Tree)

//struct that maintains metadata for every prefix
type CustomSetStruct struct {
    Desc        string `json:"desc"`
    SecondParamFrom string `json:"secondParamFrom"`
    SecondParamTo string `json:"secondParamTo"`
    Version     float64 `json:"version"`
	Value string `json:"value"`
}

type CustomSetType struct {
	SetName 	   string `json:"SetName"`
	DocumentList [] string `json:"DocumentList"`
	Prefixes []struct{
		Prefix []string `json:"Prefix"`
		DataSet CustomSetStruct `json:"DataSet"`
	} `json:"Prefixes"`
}

type PrefixesData struct {
	Prefixes []struct{
		Prefix []string `json:"Prefix"`
		DataSet CustomSetStruct `json:"DataSet"`
	} `json:"Prefixes"`
}

type CustomSet struct {
	CustomSet struct {
		Sets []CustomSetType `json:"Sets"`
	} `json:"CustomSet"`
}

var CustomSetVar CustomSet
//Init the custom set cache
func Init(config string) error {

	startTime:=time.Now()
	defer PrintRTT(startTime, "CUSTOMSET")
	err := json.Unmarshal([]byte(config), &CustomSetVar)
	if err != nil {
		return err
	}

	for setNo,Set:=range CustomSetVar.CustomSet.Sets{
		if len (Set.DocumentList) !=0 {
			for _,docId:= range Set.DocumentList{
				prefixData, err:= GetPrefixsDataFromCB(docId)
				if err != nil{
					ml.MavLog(ml.ERROR, "", "Error While fetching document from CB for docId",docId)
					return err
				}
				CustomSetVar.CustomSet.Sets[setNo].Prefixes=append(CustomSetVar.CustomSet.Sets[setNo].Prefixes, prefixData.Prefixes...)
			}
		}
	}
	
	ml.MavLog(ml.INFO, "", "no of Sets: ", len(CustomSetVar.CustomSet.Sets))
	for i := 0; i < len(CustomSetVar.CustomSet.Sets); i++ {
		setStartTime:=time.Now()
		ml.MavLog(ml.INFO, "", "Set Name: ", CustomSetVar.CustomSet.Sets[i].SetName, ", no of Prefixes: ", len(CustomSetVar.CustomSet.Sets[i].Prefixes))
		r := radix.New()
		for j := 0; j < len(CustomSetVar.CustomSet.Sets[i].Prefixes); j++ {
			for k := 0 ; k < len(CustomSetVar.CustomSet.Sets[i].Prefixes[j].Prefix) ; k++{
				r.Insert(string(CustomSetVar.CustomSet.Sets[i].Prefixes[j].Prefix[k]), &CustomSetVar.CustomSet.Sets[i].Prefixes[j].DataSet)
			}
		}
		CustomSetCache[CustomSetVar.CustomSet.Sets[i].SetName] = r
		PrintRTT(setStartTime, "SET")
	}
	return nil
}

//checks if Custom set exists.
func CustomSetExists(setName string) bool {
	if _, ok := CustomSetCache[setName]; ok {
		return true
	} else {
		return false
	}
}

//checks if prefix exists.
func PrefixExists(setName, prefix string) (error, bool) {
	found := CustomSetExists(setName)
	if found {
		_, exists := CustomSetCache[setName].Get(prefix)
		return nil, exists
	} else {
		return errors.New("setName:" + setName + " does not exist."), false
	}
}

//get prefix data
func GetPrefixData(setName, prefix string) (error, interface{}) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, nil
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), nil
	}

	val, _ := CustomSetCache[setName].Get(prefix)
	return nil, val
}

//get prefix data field desc
func GetPrefixDataFieldDesc(setName, prefix string) (error, string) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, ""
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), ""
	}

	val, _ := CustomSetCache[setName].Get(prefix)
	return nil, val.(*CustomSetStruct).Desc
}

//get prefix data field SecondParam
func GetPrefixDataFieldSecondParamFrom(setName, prefix string) (error, string) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, ""
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), ""
	}

	val, _ := CustomSetCache[setName].Get(prefix)
	return nil, val.(*CustomSetStruct).SecondParamFrom
}

func GetPrefixDataFieldSecondParamTo(setName, prefix string) (error, string) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, ""
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), ""
	}

	val, _ := CustomSetCache[setName].Get(prefix)
	return nil, val.(*CustomSetStruct).SecondParamTo
}

//get prefix data field Version
func GetPrefixDataFieldVersion(setName, prefix string) (error, float64) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, -1
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), 0
	}

	val, _ := CustomSetCache[setName].Get(prefix)
	return nil, val.(*CustomSetStruct).Version
}

//update prefix data
func UpdatePrefixData(setName, prefix string, val *CustomSetStruct) error {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName)
	}

	CustomSetCache[setName].Insert(prefix, val)
	return nil
}

//delete prefix data
func DeletePrefix(setName, prefix string) (error, bool) {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err, false
	}

	if !exists {
		return errors.New("prefix:" + prefix + " does not exist in setName:" + setName), false
	}

	_, flag := CustomSetCache[setName].Delete(prefix)
	return nil, flag 
}

//get longest prefix match for given string
func GetMatch(setName, str string) (error, bool, *CustomSetStruct) {
	if _, ok := CustomSetCache[setName]; ok {
		_, val1, found := CustomSetCache[setName].LongestPrefix(str)
		if found {
			// value := val1.(*CustomSetStruct)
			return nil, true, val1.(*CustomSetStruct)
		}else{
			return errors.New("str:" + str + " under setName:" + setName + " did not match to any prefix..."), false, nil
		}
	} else {
		return errors.New("setName:" + setName + " not found in cache..."), false, nil
	}
}

//add prefix data
func AddPrefixData(setName, prefix string, val *CustomSetStruct) error {
	err, exists := PrefixExists(setName, prefix)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("prefix:" + prefix + " already exist in setName:" + setName)
	}

	CustomSetCache[setName].Insert(prefix, val)
	return nil
}

func UpdateCustomSetConfig(){
	var newCustomSetCache = make(map[string]*radix.Tree)
	ml.MavLog(ml.INFO, "", "Length of CustomSet: ", len(CustomSetVar.CustomSet.Sets))
	for i := 0; i < len(CustomSetVar.CustomSet.Sets); i++ {
		ml.MavLog(ml.INFO, "", "Set Name: ", CustomSetVar.CustomSet.Sets[i].SetName, ", No of Prefixes: ", len(CustomSetVar.CustomSet.Sets[i].Prefixes))
		r := radix.New()
		for j := 0; j < len(CustomSetVar.CustomSet.Sets[i].Prefixes); j++ {
			for k := 0 ; k < len(CustomSetVar.CustomSet.Sets[i].Prefixes[j].Prefix) ; k++{
				r.Insert(string(CustomSetVar.CustomSet.Sets[i].Prefixes[j].Prefix[k]), &CustomSetVar.CustomSet.Sets[i].Prefixes[j].DataSet)
			}
		}
		newCustomSetCache[CustomSetVar.CustomSet.Sets[i].SetName] = r
	}

	CustomSetCache = newCustomSetCache
}

func UpdateCustomSetFromDB(config string) error {
	var customSetTemp CustomSet
	var newCustomSetCacheTemp = make(map[string]*radix.Tree)
	var mu = &sync.Mutex{}
	err := json.Unmarshal(([]byte(config)), &customSetTemp)
	if err != nil {
		ml.MavLog(ml.INFO, "", "Error occured while unmarshalling customSet ", err)
		return err
	}

	ml.MavLog(ml.INFO, "", "Length of CustomSet: ", len(customSetTemp.CustomSet.Sets))
	for i := 0; i < len(customSetTemp.CustomSet.Sets); i++ {
		ml.MavLog(ml.INFO, "", "Set Name: ", customSetTemp.CustomSet.Sets[i].SetName, ", No of Prefixes: ", len(customSetTemp.CustomSet.Sets[i].Prefixes))
		r := radix.New()
		for j := 0; j < len(customSetTemp.CustomSet.Sets[i].Prefixes); j++ {
			for k := 0 ; k < len(customSetTemp.CustomSet.Sets[i].Prefixes[j].Prefix) ; k++{
				r.Insert(string(customSetTemp.CustomSet.Sets[i].Prefixes[j].Prefix[k]), &customSetTemp.CustomSet.Sets[i].Prefixes[j].DataSet)
			}
		}
		newCustomSetCacheTemp[customSetTemp.CustomSet.Sets[i].SetName] = r
	}

	mu.Lock()
	CustomSetVar = customSetTemp
	CustomSetCache = newCustomSetCacheTemp
	mu.Unlock()
	return nil
}

//function to read customset subdocuments
func GetPrefixsDataFromCB(docId string) (PrefixesData, error){
	var prefixData PrefixesData
	prefixByteData,err := cb.GetDocContentInBytes(docId)
	if err != nil {
		return prefixData,err
	}
	err = json.Unmarshal(prefixByteData, &prefixData)
	if err != nil {
		return prefixData,err
	}
	return prefixData,nil
}

//function to print RTT
func PrintRTT(startTime time.Time, msg string){
	ml.MavLog(ml.INFO, "", msg, " loading RTT :", time.Now().Sub(startTime))
}

//function to reload/update customSet on chache update
//return error if any
func ReloadCustomSetFromDB(customSetTemp CustomSet, docId string, isFirstTime bool) error {
	if isFirstTime {
		cu.StoreData(CustomSetVar, docId)
	}
	var newCustomSetCacheTemp = make(map[string]*radix.Tree)
	var mu = &sync.Mutex{}
	for setNo,Set:=range customSetTemp.CustomSet.Sets{
		if len (Set.DocumentList) !=0 {
			for _,docId:= range Set.DocumentList{
				prefixData, err:= GetPrefixsDataFromCB(docId)
				if err != nil{
					ml.MavLog(ml.ERROR, "", "Error While fetching document from CB for docId: ",docId)
					return err
				}
				customSetTemp.CustomSet.Sets[setNo].Prefixes=append(customSetTemp.CustomSet.Sets[setNo].Prefixes, prefixData.Prefixes...)
			}
		}
	}
	
	ml.MavLog(ml.INFO, "", "", "no of Sets: ", len(customSetTemp.CustomSet.Sets))
	for i := 0; i < len(customSetTemp.CustomSet.Sets); i++ {
		ml.MavLog(ml.INFO, "", "Set Name: ", customSetTemp.CustomSet.Sets[i].SetName, ", No of Prefixes: ", len(customSetTemp.CustomSet.Sets[i].Prefixes))
		r := radix.New()
		for j := 0; j < len(customSetTemp.CustomSet.Sets[i].Prefixes); j++ {
			for k := 0 ; k < len(customSetTemp.CustomSet.Sets[i].Prefixes[j].Prefix) ; k++{
				r.Insert(string(customSetTemp.CustomSet.Sets[i].Prefixes[j].Prefix[k]), &customSetTemp.CustomSet.Sets[i].Prefixes[j].DataSet)
			}
		}
		newCustomSetCacheTemp[customSetTemp.CustomSet.Sets[i].SetName] = r
	}

	mu.Lock()
	CustomSetVar = customSetTemp
	CustomSetCache = newCustomSetCacheTemp
	mu.Unlock()
	cu.StoreData(CustomSetVar, docId)
	ml.MavLog(ml.DEBUG, "", "CustomSet config update successfull on cache update")
	return nil
}

//function to reload/update single customSet on chache update
//return error if any
func ReloadCustomSetDataFromDB(customSetTemp CustomSetType, customSetName string, isFirstTime bool) error {
	var mu = &sync.Mutex{}

	ml.MavLog(ml.INFO, "ReloadCustomSetDataFromDB", "Set Name: ", customSetTemp.SetName, ", No of Prefixes: ", len(customSetTemp.Prefixes))
	r := radix.New()
	for j := 0; j < len(customSetTemp.Prefixes); j++ {
		for k := 0 ; k < len(customSetTemp.Prefixes[j].Prefix) ; k++{
			r.Insert(string(customSetTemp.Prefixes[j].Prefix[k]), &customSetTemp.Prefixes[j].DataSet)
		}
	}

	mu.Lock()
	CustomSetCache[customSetTemp.SetName] = r
	mu.Unlock()

	ml.MavLog(ml.DEBUG, "ReloadCustomSetDataFromDB", "CustomSet",  customSetTemp.SetName, " update successfull on cache update")
	return nil
}