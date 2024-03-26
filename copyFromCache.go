package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/couchbase/gocb/v2"
)

type CopyFieldsFromCache struct {
	SubDocKey string	`json:"subDocKey"`
	RgKey string	`json:"rgKey"`
	FallbackType string	`json:"fallbackType"`
	Source		[]string	`json:"source"`
	SrcLen		int
	Destination	[]string	`json:"destination"`
	DestLen		int
}

func (c *CopyFieldsFromCache) UnmarshalJSON(data []byte) error {
    type CopyFieldsFromCache2 CopyFieldsFromCache
    if err := json.Unmarshal(data, (*CopyFieldsFromCache2)(c)); err != nil {
        return err
    }
	c.SrcLen 	= len(c.Source)
	c.DestLen 	= len(c.Destination) 
    return nil
}

func copyField(action CopyFieldsFromCache, cdr, doc1, doc2, doc3 map[string]interface{}) error {
	for i, srcfield := range action.Source {
		if fvalue, ok := doc1[srcfield]; ok {
			cdr[action.Destination[i]] = fvalue
			continue
		}else if fvalue, ok := doc2[srcfield]; ok{
			cdr[action.Destination[i]] = fvalue
			continue
		}else if fvalue, ok := doc3[srcfield]; ok{
			cdr[action.Destination[i]] = fvalue
			continue
		}else{
			fmt.Println("Field is not present in all the documents")
			return errors.New("Field is not present in all the documents")
		}
	}
	return nil
}
func main(){
	var action CopyFieldsFromCache
	readFileData, err := ioutil.ReadFile("sampleActionCopyFromCache.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)	
	}
	cdr := make(map[string]interface{})
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
	cdr["sessionKey"] = "sessionKey"
	cdr["rgKey"] = "rgKey"
	
	value := "RG->SESSION_RG->SESSION"
	switch value {
		case "RG->SESSION_RG->SESSION":
			if _, ok := cdr[action.RgKey]; !ok {
				return 
			}
			if _, ok := cdr[action.SubDocKey]; ok {
				return
			}
			rgdoc := make(map[string]interface{})
			response, err := collection.Get(cdr[action.RgKey].(string), nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response.Content(&rgdoc)
			sessionRgDoc := make(map[string]interface{})
			response1, err := collection.Get(cdr[action.SubDocKey].(string), nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response1.Content(&sessionRgDoc)
			sessionDoc := make(map[string]interface{})
			response2, err := collection.Get(cdr["key"].(string), nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response2.Content(&sessionDoc)
	
			for i, srcfield := range action.Source {
				if fvalue, ok := rgdoc[srcfield]; ok {
					cdr[action.Destination[i]] = fvalue
					continue
				}else if fvalue, ok := sessionRgDoc[srcfield]; ok{
					cdr[action.Destination[i]] = fvalue
					continue
				}else if fvalue, ok := sessionDoc[srcfield]; ok{
					cdr[action.Destination[i]] = fvalue
					continue
				}else{
					fmt.Println("Field is not present in all the documents")
				}
			}
		case "SESSION_RG->SESSION":
			sessionRgDoc := make(map[string]interface{})
			response1, err := collection.Get(cdr[action.SubDocKey].(string), nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response1.Content(&sessionRgDoc)
			sessionDoc := make(map[string]interface{})
			response2, err := collection.Get(cdr["key"].(string), nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response2.Content(&sessionDoc)
			for i, srcfield := range action.Source {
				if fvalue, ok := sessionRgDoc[srcfield]; ok {
					cdr[action.Destination[i]] = fvalue
					continue
				} else if fvalue, ok := sessionDoc[srcfield]; ok {
					cdr[action.Destination[i]] = fvalue
					continue
				} else {
					fmt.Println("Field is not present in all the documents")
				}
			}
		case "CURRENTRULEKEY":
			doc := make(map[string]interface{})
			response, err := collection.Get("key", nil)
			if err != nil {
				fmt.Println("Failed to get document from couchbase")
			}
			response.Content(&doc)
			for i := 0; i < action.SrcLen; i++ {
				if fvalue, ok := doc[action.Source[i]]; ok {
					cdr[action.Destination[i]] = fvalue
				}else {
					fmt.Println("Field --", action.Source[i] , " is not present, moving on to the next field" )
				}
			}
		default : fmt.Println("Invalid configuration for fallback type")

			
	}
}