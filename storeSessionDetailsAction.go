package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/couchbase/gocb/v2"
)

var (
	collection *gocb.Collection
	 Bucket     *gocb.Bucket 
)

func main(){
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

	cdr1 := make(map[string]interface{})
	var action1 StoreSessionDetails
	readFileData1, err := ioutil.ReadFile("sampleCDRForSuDocOperations.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	cdr2 := make(map[string]interface{})
	readFileData2, err := ioutil.ReadFile("sampleCDRForSuDocOperations2.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	readActionData1, err := ioutil.ReadFile("sampleActionSubDocOperations.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData1, &cdr1)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	err = json.Unmarshal(readFileData2, &cdr2)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	err = json.Unmarshal(readActionData1, &action1)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	actionStoreSessionDetails(&action1, "key", cdr1, nil)
	actionStoreSessionDetails(&action1, "key", cdr2, nil)
}

type StoreSessionDetails struct {
	SubDocKey							string		`json:"subDocKey"`	// subdockey against which this info will be stored
	SessionDocKeySourceFields			[]struct {
		Source 			string 	`json:"source"`
		KeyType			string	`json:"keyType"`
	}												`json:"sessionDocKeySourceFields"`	// will contain value such as SessionKey , RGKey, Level2Dockey
	SessionDocKeysDestinationSubDocKey	string		`json:"sessionDocKeysDestinationSubDocKey"`
	ListOfInitialValueFields 			[]struct {	// in this if value is already populated then will not overwrite
		Source 			string `json:"source"`		// for initialCDCID
		Destination 	string `json:"destination"`
	}												`json:"listOfInitialValueFields"`
	ListOfLatestValueFields 			[]struct {	// in this if value is already populated then it will overwrite the previous one
		Source 			string `json:"source"`		// for current cdcid and lastProcessedRecSeqNum
		Destination 	string `json:"destination"`
	}												`json:"listOfLatestValueFields"`
	Err									error
	LenOfSessionDocKeySourceFields	int
	LenOfListOfInitialValueFields	int
	LenOfListOfLatestValueFields	int
}

func (d *StoreSessionDetails) UnmarshalJSON(data []byte) error {
	type StoreSessionDetailsV1 StoreSessionDetails
    if err := json.Unmarshal(data, (*StoreSessionDetailsV1)(d)); err != nil {
        return err
    }
	if d.SubDocKey == "" {
		d.Err = errors.New("SubdocKey is not configured")
	}
	if d.SessionDocKeysDestinationSubDocKey == "" {
		d.Err = errors.New("SessionDocKeysDestinationSubDocKey is not configured")
	}
	d.LenOfSessionDocKeySourceFields = len(d.SessionDocKeySourceFields)
	if d.LenOfSessionDocKeySourceFields == 0 {
		d.Err = errors.New("SessionDocKeySourceFields is empty")
	}
	d.LenOfListOfInitialValueFields = len(d.ListOfInitialValueFields)
	if d.LenOfListOfInitialValueFields == 0 {
		d.Err = errors.New("ListOfInitialValueFields is empty")
	}
	d.LenOfListOfLatestValueFields = len(d.ListOfLatestValueFields)
	if d.LenOfListOfLatestValueFields == 0 {
		d.Err = errors.New("ListOfLatestValueFields is empty")
	}
    return nil
}

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

func InsertSubDoc(key, path string, doc interface{}) error {
	muOps := []gocb.MutateInSpec{
		gocb.UpsertSpec(path, &doc, &gocb.UpsertSpecOptions{CreatePath: true}),
	}
	_, err := collection.MutateIn(key, muOps, &gocb.MutateInOptions{})
	return err
}

func GetSubDoc(key, path string) (error, interface{}) {
	
	ops := []gocb.LookupInSpec{
		gocb.GetSpec(path, &gocb.GetSpecOptions{}),
	}
	couchOp := func(key, path string) (interface{}, error) {
		return collection.LookupIn(key, ops, &gocb.LookupInOptions{})
	}
	response, err := couchOp(key, path)
	if err != nil {
		fmt.Println("GetSubDoc", "Failed to access content for subdoc for key:"+key, err)
		return err, nil
	}
	crdlResponse := response.(*gocb.LookupInResult)
	var val interface{}
	crdlResponse.ContentAt(0, &val)
	return nil, val
}

func actionStoreSessionDetails(action *StoreSessionDetails, key string, cdr map[string]interface{}, icWrapper interface{}) error {
    if action.Err != nil {
        fmt.Println(action.Err)
        return action.Err
    }
    var keyMap map[string]interface{}
    // fetch subDoc from session doc
    err, doc := GetSubDoc(key, action.SubDocKey)
    if err != nil {
        fmt.Println("Failed to check if sub doc exists for key:", key, "and sub doc path:", action.SubDocKey, ",error:", err)
    }
    if doc == nil {
       doc = make(map[string]interface{})
    }
    if subDoc, ok := doc.(map[string]interface{}); ok {
        if keyMap, ok = subDoc[action.SessionDocKeysDestinationSubDocKey].(map[string]interface{}); !ok {
            keyMap = make(map[string]interface{})
        }
        for i := 0; i < action.LenOfSessionDocKeySourceFields; i++ {
            if fieldVal, present := cdr[action.SessionDocKeySourceFields[i].Source].(string); present && fieldVal != "" {
                keyMap[fieldVal] = action.SessionDocKeySourceFields[i].KeyType
            }
        }
        subDoc[action.SessionDocKeysDestinationSubDocKey] = keyMap
        for i := 0; i < action.LenOfListOfInitialValueFields; i++ {
            if fieldVal, present := cdr[action.ListOfInitialValueFields[i].Source]; present {
                switch fieldVal.(type) {
                case string:
                    if _, ok := subDoc[action.ListOfInitialValueFields[i].Destination]; !ok && fieldVal != "" {
                        subDoc[action.ListOfInitialValueFields[i].Destination] = fieldVal
                    }
                case float64, int:
                    if _, ok := subDoc[action.ListOfInitialValueFields[i].Destination]; !ok {
                        subDoc[action.ListOfInitialValueFields[i].Destination] = fieldVal
                    }
                }
            }
        }
        for i := 0; i < action.LenOfListOfLatestValueFields; i++ {
            if fieldVal, present := cdr[action.ListOfLatestValueFields[i].Source]; present {
                switch fieldVal.(type) {
                case string:
                    if fieldVal != "" {
                        subDoc[action.ListOfLatestValueFields[i].Destination] = fieldVal
                    } 
                case float64, int:
                    subDoc[action.ListOfLatestValueFields[i].Destination] = fieldVal
                }
            }
        }
        err := InsertSubDoc(key, action.SubDocKey, subDoc)
        if err != nil {
            fmt.Println("Failed to insert subdoc:", action.SubDocKey, "with key:", key, ", error:", err)
            return errors.New("Failed to insert subdoc with subDockey :" + action.SubDocKey + " and key :" + key)
        }
        fmt.Println("Successfully stored cdr in couchbase", subDoc)
    } else {
        fmt.Println("Error occured while typecasting subdoc in session doc")
        return errors.New("Error occured while typecasting subdoc in session doc")
    }
    return nil
}
func InsertDoc(key string, value map[string]interface{}) error {
	
	upsertOptions := &gocb.UpsertOptions{}
	upsertOptions.Expiry = time.Duration(0) * time.Second
	
	couchOp := func(key string, value map[string]interface{}) (interface{}, error) {
		return collection.Upsert(key, &value, upsertOptions)
	}
	_, err := couchOp(key, value)
	if err != nil {
		return err
	}
	return err
}