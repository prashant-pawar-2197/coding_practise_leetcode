package main 

import (
	"fmt"
)
type StoreSeqNum struct {
	SrcField 			string `json:"srcField"`
	FirstSeqNumField	string 	`json:"firstSeqNumField"`
	LastSeqNumField  	string	`json:"lastSeqNumField"`
	SeqNumListField		string	`json:"SeqNumListField"`
	SeqNumOpType		string	`json:"SeqNumOpType"`
	SubDocKeyToStoreMissingRecSeqNums	string		`json:"subDocKeyToStoreMissingRecSeqNums"`
	SessionKey							string		`json:"sessionKey"`
	SessionKeyKtab						string		`json:"sessionKeyKtab"`
	SessionInfoSubdocKey				string		`json:"sessionInfoSubdocKey"`
	SubDocKeyToStoreCountOfAggrCdrs		string		`json:"subDocKeyToStoreCountOfAggrCdrs"`
}

const (
	CGFM_SESSION_SUB_DOC_KEY		= "SessionInfo"
	CGFM_LAST_PROCESSED_SEQ_NUM		= "lastProcessedRecSeqNum"
)

func actionStoreSeqNum(action *StoreSeqNum, key string, cdr map[string]interface{}) error {
    UseKtab := true
	var (
        sessionKey      string
        sessionInfo     map[string]interface{}
        ok              bool
        recordSeqNum    float64
    )
    recordSeqNum, ok = cdr[action.SrcField].(float64)
	if !ok {
		fmt.Println("Field:", action.SrcField, " Not present in CDR")	
		return nil
	}
    // store any missing record sequence number
    // fetch sessionDoc key from cdr
    if sessionKey, ok = cdr[action.SessionKey].(string); ok {
        if UseKtab {
            sessionKey = action.SessionKeyKtab + sessionKey
        }
        err, subDoc := GetSubDoc(sessionKey, CGFM_SESSION_SUB_DOC_KEY)
        if err != nil {
            fmt.Println( "Key:", sessionKey,"SubdocKey", CGFM_SESSION_SUB_DOC_KEY, "failed to get sub document from DB. error:", err)
        }
        if sessionInfo, ok = subDoc.(map[string]interface{}); ok {
            if lastProcessedRecSeqNum, ok := sessionInfo[CGFM_LAST_PROCESSED_SEQ_NUM].(float64); ok {
                if (recordSeqNum - lastProcessedRecSeqNum) > 1 {
                    err, subDoc := GetSubDoc(key, action.SubDocKeyToStoreMissingRecSeqNums)
                    if err != nil {
                        fmt.Println( "Key:", key,"SubdocKey", action.SubDocKeyToStoreMissingRecSeqNums, "failed to get sub document from DB. error:", err)
                    }
                    missingRecSeqNumList := make([]float64, 0)
                    if subDoc != nil {
                        if missingRecSeqNumList , ok = subDoc.([]float64); ok {
                            for i := lastProcessedRecSeqNum+1; i < recordSeqNum; i++ {
                                missingRecSeqNumList = append(missingRecSeqNumList, i)
                            }
                        }
                    } else {
                        for i := lastProcessedRecSeqNum; i < recordSeqNum; i++ {
                            missingRecSeqNumList = append(missingRecSeqNumList, i)
                        }
                    }
                    err = icWrapper.IntransitCDRs.InsertSubDoc(key, action.SubDocKeyToStoreMissingRecSeqNums, missingRecSeqNumList)
                    if err != nil {
                        fmt.Println( "Key:", key, "Subdockey", action.SubDocKeyToStoreMissingRecSeqNums, "failed to insert document:,", missingRecSeqNumList, "error:", err)
                    }
                } else {
                    err, subDoc := icWrapper.IntransitCDRs.GetSubDoc(key, action.SubDocKeyToStoreMissingRecSeqNums)
                    if err != nil {
                        fmt.Println( "Key:", key,"SubdocKey", action.SubDocKeyToStoreMissingRecSeqNums, "failed to get sub document from DB. error:", err)
                    }
                    if subDoc != nil {
                        if missingSeqList , ok := subDoc.([]float64); ok {
                            var (
                                elemDeleted bool
                                newDeletedRecordSeqNumList []float64
                            )

                            for _, recordSeqNumberFromList := range missingSeqList {
                                if recordSeqNum == recordSeqNumberFromList {
                                    newDeletedRecordSeqNumList = deleteElementFloat64(missingSeqList, recordSeqNum)
                                    if !elemDeleted {
                                        elemDeleted = true
                                    }
                                }
                            }
                            if elemDeleted {
                                err = InsertSubDoc(key, action.SubDocKeyToStoreMissingRecSeqNums, newDeletedRecordSeqNumList)
                                if err != nil {
                                    fmt.Println( "Key:", key, "Subdockey", action.SubDocKeyToStoreMissingRecSeqNums, "failed to insert document:,", missingRecSeqNumList, "error:", err)
                                }
                            }
                        }
                    }
                }
            } else {
                fmt.Println( "Field:", lastProcessedRecSeqNum ," is missing in the sessionInfo subdoc, subdoc:", sessionInfo)
            }
        }
    }

	switch action.SeqNumOpTypeEnum {
		case gpp.CGFM_STORE_ALL_RECORD_SEQ_NUM_PER_SLOT:
			seqNumList := make([]interface{}, 0)
			if err, fieldExists, fieldValue := icWrapper.IntransitCDRs.GetDocField(key, action.SeqNumListField); err != nil {
				fmt.Println( "Key:", key, "Field:", action.SeqNumListField, "failed to get field. err:", err)
				seqNumList = append(seqNumList, recordSeqNum)
			} else if fieldExists == false {
				fmt.Println("Key:", key, "Field:", action.SeqNumListField, "not present in doc")
				seqNumList = append(seqNumList, recordSeqNum)
			} else {
				ok := false
				if seqNumList, ok = fieldValue.([]interface{}); ok {
					seqNumList = append(seqNumList, recordSeqNum)
				}
			}
		 	if err := icWrapper.IntransitCDRs.SetDocField(key, action.SeqNumListField, seqNumList); err != nil {
		 		fmt.Println( "Key:", key, "Field:", action.SeqNumListField, "failed to insert value to doc. err:", err)
				return err
			}
		case gpp.CGFM_STORE_FIRST_AND_LAST_RECORD_SEQUENCE_NUM:
			if err, fieldExists, _ := icWrapper.IntransitCDRs.GetDocField(key, action.FirstSeqNumField);  err != nil {
				fmt.Println( "Key:", key, "Field:", action.FirstSeqNumField, "failed to get field. err:", err)
				return err
			} else if fieldExists == false {
		 		if err := icWrapper.IntransitCDRs.SetDocField(key, action.FirstSeqNumField, recordSeqNum); err != nil {
		 			fmt.Println( "Key:", key, "Field:", action.FirstSeqNumField, "failed to insert value to doc. err:", err)
					return err
				}
			} 
			
		 	if err := icWrapper.IntransitCDRs.SetDocField(key, action.LastSeqNumField, recordSeqNum); err != nil {
		 		fmt.Println( "Key:", key, "Field:", action.LastSeqNumField, "failed to insert value to doc. err:", err)
				return err
			}
	}
	return nil
}

// for float64
func deleteElementFloat64(arr []float64, elem float64) []float64 {
	var newArr []float64
	for _, value := range arr {
		if value != elem {
			newArr = append(newArr, value)
		}
	}
	return newArr
}

func SetDocField(key, field string, val interface{}) error {
	err := InsertSubDoc(key, field, val)
	if err != nil {
		
	}
}
func main(){

}