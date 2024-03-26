package main

import (
	"fmt"
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
)
var (
	collection *gocb.Collection
 	Bucket     *gocb.Bucket
)
func main() {
	username := os.Args[1]
	pass := os.Args[2]
	ipAddr := os.Args[3]
	bucketName := os.Args[4]
	fmt.Println("Username -->", username)
	fmt.Println("Password -->", pass)
	fmt.Println("IpAddr -->", ipAddr)
	fmt.Println("BucketName -->", bucketName)
    opts := gocb.ClusterOptions{Username: username, Password: pass}
	cluster, err := gocb.Connect(ipAddr, opts)
	if err != nil {
		fmt.Println("error in connection: ", err)
	}
	if err == nil {
		fmt.Println("connection success ")
	}
	Bucket = cluster.Bucket(bucketName)

	//Wait upto timout and checks database conn status
	if err = Bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		fmt.Println("Couchbase bootstrapError : %v", err)
	}
	col := Bucket.DefaultCollection()
	collection = col
	m1 := make(map[string]interface{})
	m1["Bharath"] = "SV"
	fmt.Println("Doc Inserted")
	resultwithOptions, err := collection.Upsert("BharathSampleDoc", &m1, &gocb.UpsertOptions{
		Timeout: 3 * time.Second,
		Expiry:  0 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resultwithOptions)
	fmt.Println("Doc fetched")
	updateGetResult, err := collection.Get("BharathSampleDoc", nil)
	if err != nil {
		panic(err)
	}
	var doc map[string]interface{}
	err = updateGetResult.Content(&doc)
	if err != nil {
		panic(err)
	}
	fmt.Println(doc)
}

/*
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
	m1 := make(map[string]interface{})
	m1["Prashant"] = "Pawar"
	sl := make([]interface{}, 2)
	sl = append(sl, 1)
	sl = append(sl, 2)
	m1["OtherData"] = sl
	fmt.Println("Doc Inserted")
	resultwithOptions, err := collection.Upsert("document-key-options", &m1, &gocb.UpsertOptions{
		Timeout: 3 * time.Second,
		Expiry:  0 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resultwithOptions)
	fmt.Println(time.Now())
	time.Sleep(10 * time.Second)
	var num int = 60
	fmt.Println("Expiry Set")
	touchResult, err := collection.Touch("document-key-options", time.Duration(num)*time.Second, &gocb.TouchOptions{
		Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		panic(err)
	}
	sl = append(sl, 3)
	m1["OtherData"] = sl
	time.Sleep(10 * time.Second)
	resultwithOptions, err = collection.Upsert("document-key-options", &m1, &gocb.UpsertOptions{
		Timeout: 3 * time.Second,
	})
	fmt.Println("Doc Upserted")
	fmt.Println(touchResult)
*/