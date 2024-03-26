package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
)

// perform operations on a field inside the given document which is an array...
type SubDocOperations struct {
	OperationType		string		`json:"operationType"` // Can be INSERT or DELETE
	SubDocKey			string		`json:"subDocKey"`
	ListOfFields		[]string	`json:"listOfFields"`
	SubDocFieldLength 	string		`json:"subDocFieldLength"`
	Err 				error
	LenOfListOfFields	int
}

type Flag struct {
    InfoLogFlag     bool
}

var cu Flag

var (
	 	collection *gocb.Collection
	  	Bucket     *gocb.Bucket 
)

func (d *SubDocOperations) UnmarshalJSON(data []byte) error {
    cu.InfoLogFlag = true
	type SubDocOperationsV1 SubDocOperations
    if err := json.Unmarshal(data, (*SubDocOperationsV1)(d)); err != nil {
        return err
    }
	if d.OperationType == "" {
		d.Err = errors.New("OperationType is not configured")
	}
	if d.SubDocKey == "" {
		d.Err = errors.New("SubDocKey is not configured")
	}
	lenOfListFields := len(d.ListOfFields)
	if lenOfListFields == 0 {
		d.Err = errors.New("ListOfFields is empty")
	} else {
		d.LenOfListOfFields = lenOfListFields
	}
	if d.SubDocFieldLength == "" {
		d.Err = errors.New("SubDocFieldLength is not configured")
	}
    return nil
}

func actionSubDocOperations(action *SubDocOperations, key string, cdr map[string]interface{}) error {
    if action.Err != nil {
        fmt.Println(action.Err.Error())
        return action.Err
    }
    var (
        listOfEmptyFields string
        delimiter string = ","
    )
    switch action.OperationType {
    case "INSERT":
        updateGetResult, err := collection.Get("document-key", nil)
        cbDoc := make(map[string]interface{}, 0)
        err = updateGetResult.Content(&cbDoc)
        if err != nil {
            panic(err)
        }
        doc, _ := cbDoc[action.SubDocKey].([]interface{})
        if doc == nil {
            var SubDoc []interface{}
            for i := 0; i < action.LenOfListOfFields; i++ {
                if val, ok := cdr[action.ListOfFields[i]]; ok {
                    SubDoc = append(SubDoc, val)
                } else {
                    listOfEmptyFields = listOfEmptyFields + delimiter + action.ListOfFields[i]
                    continue
                }
            }
            if listOfEmptyFields != "" {
				if cu.InfoLogFlag {
                	fmt.Println("List Of fields which were empty in the cdr :", listOfEmptyFields)
            	}
            }
            lenOfAppendedDoc := len(SubDoc)
            if lenOfAppendedDoc > 0 {
                cbDoc[action.SubDocKey] = SubDoc
                cdr[action.SubDocFieldLength] = lenOfAppendedDoc
                _, err := collection.Upsert("document-key", &cbDoc, &gocb.UpsertOptions{
                    Timeout: 3 * time.Second,
                })
                if err != nil {
                    fmt.Println(err)
                }
            }
            return nil
        } else {
            originalLenOfSubDoc := len(doc)
            for i := 0; i < action.LenOfListOfFields; i++ {
                if val, ok := cdr[action.ListOfFields[i]]; ok {
                    for _, elem := range doc {
                        switch elem.(type) {
                        case string:
                            if reflect.TypeOf(val).Kind() ==  reflect.String {
                                if !strings.EqualFold(elem.(string), val.(string)) {
                                    doc = append(doc, val)
                                } else {
									if cu.InfoLogFlag {
                                    	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                                	}
                                }
                            }
                        case int:
                            if reflect.TypeOf(val).Kind() ==  reflect.Int {
                                if elem.(int) != val.(int) {
                                    doc = append(doc, val)
                                } else {
									if cu.InfoLogFlag {
                                    	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                                	}
                                }
                            }
                        case float64:
                            if reflect.TypeOf(val).Kind() ==  reflect.Float64 {
                                if elem.(float64) != val.(float64) {
                                    doc = append(doc, val)
                                } else {
									if cu.InfoLogFlag {
                                    	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                                	}
                                }
                            }
                        default:
                            fmt.Println("Unsupported data type while forming CB key")
                        }
                    }
                } else {
                    listOfEmptyFields = listOfEmptyFields + delimiter + action.ListOfFields[i]
                    continue
                }
            }
            // if addition happened then subdoc length will increase and only then we'll insert the subdoc
            lenOfAppendedDoc :=  len(doc)
            if originalLenOfSubDoc < lenOfAppendedDoc {
                cbDoc[action.SubDocKey] = doc
                cdr[action.SubDocFieldLength] = lenOfAppendedDoc
                _, err := collection.Upsert("document-key", &cbDoc, &gocb.UpsertOptions{
                    Timeout: 3 * time.Second,
                })
                if err != nil {
                    fmt.Println(err)
                }
            }
        }
    case "DELETE":
        updateGetResult, err := collection.Get("document-key", nil)
        cbDoc := make(map[string]interface{}, 0)
        err = updateGetResult.Content(&cbDoc)
        if err != nil {
            panic(err)
        }
        doc, _ := cbDoc[action.SubDocKey].([]interface{})
        originalLenOfSubDoc := len(doc)
        if originalLenOfSubDoc == 0 {
            fmt.Println("Subdoc list is already empty")
            cdr[action.SubDocFieldLength] = 0
            return errors.New("Subdoc list is already empty")
        }
        for i := 0; i < action.LenOfListOfFields; i++ {
            if val, ok := cdr[action.ListOfFields[i]]; ok {
                for _, elem := range doc {
                    switch elem.(type) {
                    case string:
                        if reflect.TypeOf(val).Kind() ==  reflect.String {
                            if strings.EqualFold(elem.(string), val.(string)) {
                                doc = deleteElement(doc, val)
                            } else {
								if cu.InfoLogFlag {
                                	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                            	}
                            }
                        }
                    case int:
                        if reflect.TypeOf(val).Kind() ==  reflect.Int {
                            if elem.(int) == val.(int) {
                                doc = deleteElement(doc, val)
                            } else {
								if cu.InfoLogFlag {
                                	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                            	}
                            }
                        }
                    case float64:
                        if reflect.TypeOf(val).Kind() ==  reflect.Float64 {
                            if elem.(float64) == val.(float64) {
                                doc = deleteElement(doc, val)
                            } else {
								if cu.InfoLogFlag {
                                	fmt.Println("Value : ", val, " already present in the subdoc, hence not adding") 
                            	}
                            }
                        }
                    default:
                        fmt.Println("Unsupported data type while forming CB key")
                    }
                }
            } else {
                listOfEmptyFields = listOfEmptyFields + delimiter + action.ListOfFields[i]
                continue
            }
        }
        // if addition happened then subdoc length will increase and only then we'll insert the subdoc
        lenOfAppendedDoc := len(doc)
        if originalLenOfSubDoc > lenOfAppendedDoc {
            cbDoc[action.SubDocKey] = doc
            cdr[action.SubDocFieldLength] = lenOfAppendedDoc
            _, err := collection.Upsert("document-key", &cbDoc, &gocb.UpsertOptions{
                Timeout: 3 * time.Second,
            })
            if err != nil {
                fmt.Println(err)
            }
        }  
    }
    return nil
}

func deleteElement(arr []interface{}, elem interface{}) []interface{} {
	var newArr []interface{}
	for _, value := range arr {
		if !reflect.DeepEqual(value, elem) {
			newArr = append(newArr, value)
		}
	}
	return newArr
}

func main()  {
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

	cdr := make(map[string]interface{})
	var action SubDocOperations
	readFileData, err := ioutil.ReadFile("sampleCDRForSuDocOperations.json")
	if err != nil {
		fmt.Println("Error occured")
	}

	readActionData, err := ioutil.ReadFile("sampleActionSubDocOperations.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	err = json.Unmarshal(readFileData, &cdr)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	err = json.Unmarshal(readActionData, &action)
	if err != nil {
		fmt.Println("Error occured while unmarshalling", err)
	}
	actionSubDocOperations(&action,"document-key",cdr)
    cdr["RGKey"] = 2.0
    actionSubDocOperations(&action,"document-key",cdr)
    cdr["RGKey"] = 1.0
    actionSubDocOperations(&action,"document-key",cdr)
    actionSubDocOperations(&action,"document-key",cdr)
    action.OperationType = "DELETE"
    actionSubDocOperations(&action,"document-key",cdr)
	encodedData, _ := json.Marshal(cdr)
	ioutil.WriteFile("pulseSplitOutput", encodedData, 0644)

}
