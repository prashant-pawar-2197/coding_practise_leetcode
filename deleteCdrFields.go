package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/armon/go-radix"
)

const RAWCDR  = "rawCDR"
const (
	CGFM_RADIX_TRIE = iota
    CGFM_STRING_HAS_PREFIX 
)

type EnumAlgoType int

func ConvertAlgoTypeFromStringToEnum( typestr string ) EnumAlgoType {
	switch typestr {
		case "CGFM_RADIX_TRIE":
			return CGFM_RADIX_TRIE
		case "CGFM_STRING_HAS_PREFIX":
			return CGFM_STRING_HAS_PREFIX
	}
	return CGFM_RADIX_TRIE
}

type DeleteCdrFields struct {
	Tree    		*radix.Tree
	DeleteAlgoType 	EnumAlgoType
	DeleteListType 	EnumListType
	AlgoType		string `json:"algoType"`
	Location 		string `json:"location"`
	Fields 			[]string `json:"field"`
	ListType		string `json:"listType"`

}

type EnumListType int

const (
	CGFM_DELETE_FROM_LIST = iota
    CGFM_ALLOWED_IN_LIST 
)

func ConvertListTypeFromStringToEnum( typestr string ) EnumListType {
	switch typestr {
		case "CGFM_DELETE_FROM_LIST":
			return CGFM_DELETE_FROM_LIST
		case "CGFM_ALLOWED_IN_LIST":
			return CGFM_ALLOWED_IN_LIST
	}
	return CGFM_DELETE_FROM_LIST
}

func (d *DeleteCdrFields) UnmarshalJSON(data []byte) error {
    type DeleteCdrFields2 DeleteCdrFields
    if err := json.Unmarshal(data, (*DeleteCdrFields2)(d)); err != nil {
        return err
    }

	d.DeleteAlgoType = ConvertAlgoTypeFromStringToEnum(d.AlgoType)
	d.DeleteListType = ConvertListTypeFromStringToEnum(d.ListType)

	if d.DeleteAlgoType == CGFM_RADIX_TRIE {
		d.Tree = radix.New()
		for _, prefix := range d.Fields {
        	output := true
        	d.Tree.Insert(string(prefix), &output)
    	}
	}
    return nil
}

/* Deleting unnecessary fields from the cdrs */
func actionDeleteCdrFields(action DeleteCdrFields, cdr map[string]interface{}, transID string) error {
	switch action.Location {
	case RAWCDR:
		return action.DeleteFields(cdr, transID)
	default:
		fmt.Println("Unsupported Source Location:", action.Location, " configured.")
		return errors.New("Unsupported Source Location:" + action.Location + " configured.")
	}
	return nil
}


const WIDE_CARD_STAR = "*"

/* Deleting unnecessary fields from the cdrs */
func (d* DeleteCdrFields) DeleteFields(cdr map[string]interface{}, transID string) error {
	switch d.DeleteAlgoType {
	case CGFM_STRING_HAS_PREFIX:
		//Switch
		switch d.DeleteListType {
		case CGFM_DELETE_FROM_LIST:
			for _, field := range d.Fields {
				if strings.Contains(field, WIDE_CARD_STAR) {
					actualField := field[:len(field)-1]
					for cdrField := range cdr {
						if strings.Contains(cdrField, actualField) {
							delete(cdr, cdrField)
						}
					}
				} else {
					if _, ok := cdr[field]; ok {
						delete(cdr, field)
					}
				}
			}
		default:
			fmt.Println("Invalid operation", d.DeleteListType)

		}
	case CGFM_RADIX_TRIE:
		for key, _ := range cdr {
			_, _, exists := d.Tree.LongestPrefix(key);
			switch d.DeleteListType {
				case CGFM_DELETE_FROM_LIST:
					if exists {
						delete(cdr, key)
					}
				case CGFM_ALLOWED_IN_LIST:
					if !exists {
						delete(cdr, key)
					}
				default:
					fmt.Println("Invalid operation", d.DeleteListType)
			}
    	}
	}
	return nil
}

func main () {
	cdr := make(map[string]interface{})
	var action DeleteCdrFields
	readFileData, err := ioutil.ReadFile("deleteFieldSampleCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}

	readActionData, err := ioutil.ReadFile("deleteCdrAction.json")
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

	actionDeleteCdrFields(action , cdr, "transID")
	encodedData,_ := json.Marshal(cdr)
	ioutil.WriteFile("deleteCdrOutput",encodedData, 0644)
}