package main

import (
	//"encoding/json"
	//	"errors"
	//"errors"
	"fmt"
	"reflect"

	//	"math"
	//	"reflect"
	"time"

	"github.com/couchbase/gocb/v2"
)

var (
	collection *gocb.Collection
	Bucket     *gocb.Bucket
	// insertOpts *gocb.InsertOptions
	// DupCheckDocTTL time.Duration
	// dupcheckMap = make(map[string]interface{})
)

func GetDocumentAll(key string) (error, bool, map[string]interface{}) {

	_, exists := DocumentExist(key)
	if exists {
		couchOp := func(key string) (interface{}, error) {
			return collection.Get(key, nil)
		}
		response, err := couchOp(key)
		if err != nil {
			fmt.Println("GetDocumentAll", "error getting value for key:"+key, err)
			return err, false, nil
		}
		crdlResponse := response.(*gocb.GetResult)
		docMap := make(map[string]interface{})
		crdlResponse.Content(&docMap)
		return nil, true, docMap
	} else {
		fmt.Println("GetDocumentAll", "document dosn't exist in couchbase for key:", key)
		return nil, false, nil
	}
}

func DocumentExist(key string) (error, bool) {
	couchOp := func(key string) (interface{}, error) {
		return collection.Exists(key, nil)
	}
	resp, err := couchOp(key)
	if err != nil {
		fmt.Println("DocumentExist", "error checking value exists for key:"+key, err)
		return err, false
	}
	crdlResponse := resp.(*gocb.ExistsResult)
	return nil, crdlResponse.Exists()
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
	key := "restoreDoc"

	err, _, doc := GetDocumentAll(key)
	if err != nil {
		fmt.Println(err)
	}
	var (
		exist 	bool
		preDocs	= make(map[string]interface{})	
	)
	fmt.Println(reflect.TypeOf(doc["PREDocs"]))
	if preDocs , exist = doc["PREDocs"].(map[string]interface{}); exist {
		if cdrDocKeysArr, ok := preDocs["0"].([]interface{}); ok {
			cdrDocKeysArr = append(cdrDocKeysArr, "cdrkey")
			fmt.Println(cdrDocKeysArr...)
		}
	}
	//path := "click"
	// ops := []gocb.LookupInSpec{
	// 	gocb.GetSpec(path, &gocb.GetSpecOptions{}),
	// }
	// ops := []gocb.LookupInSpec{
	// 	gocb.ExistsSpec(path, &gocb.ExistsSpecOptions{}),
	// }

	// result, err := collection.LookupIn(key, ops, &gocb.LookupInOptions{})
	// if err != nil {
	// 	fmt.Println("Error occured :", err)
	// 	if errors.Is(err, gocb.ErrDocumentNotFound) {
	// 		fmt.Println("&&&&&& Error occured :", err)
	// 	}
	// }
	// var val map[string]interface{}
	// fmt.Println(result)
	// err = result.ContentAt(0, &val)
	// if err != nil {
	// 	if errors.Is(err, gocb.ErrPathNotFound){
	// 		fmt.Println("Error occured ---->>>>>:",err)
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// }
	// fmt.Println("Value------>",val)
	/*
	pings, err := Bucket.Ping(&gocb.PingOptions{
		ReportID:     "couchbase_healthcheck_report",
		ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeKeyValue},
	})
	if err != nil {
		fmt.Println("Cluster Ping failed During ", err)
	}
	pingReport := pings.Services[gocb.ServiceTypeKeyValue]
	fmt.Println("total Node count in cluster ", len(pingReport))

	fmt.Println(reflect.TypeOf(len(pingReport)))
	//Update max node failure Threshold count as per the configured percent
	// int(len(pingReport) * (cdbCfg.MaxNodeFailureThresholdPerc / 100))
	var a int = 5
	var per int = 80
	fmt.Println("val",  int(math.Floor(float64(a) * (float64(per)*1.0)/100)))
	fmt.Println(reflect.TypeOf(a * int(math.Floor((80*1.0)/100))))
	*/
}
