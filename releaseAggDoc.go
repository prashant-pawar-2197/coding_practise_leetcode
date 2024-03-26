package main

import (
	"fmt"
)
func actionReleaseAggregatedDoc(action gpp.Action, key string, transID string, outCdrList interface{}, cdr map[string]interface{}, mapAccessInfo mc.CbMapAccessInfo) error {
    if action.RELEASEAGGREGATEDDOC.SubDocKey == "" {
        fmt.Println( "Subdockey is not configured hence not proceeding further")
		return errors.New("Subdockey is not configured hence not proceeding further")
    }
    if len(action.RELEASEAGGREGATEDDOC.Lookup) == 0 {
        fmt.Println( "Error occured while unmarshalling hence not proceeding further")
		return errors.New("Error occured while unmarshalling hence not proceeding further")
    }
    var (
        subDocKey string
        ok bool
        lookupKey string
    )
    if subDocKey , ok = cdr[action.RELEASEAGGREGATEDDOC.SubDocKey].(string); !ok {
        fmt.Println( "Subdockey is not present in the cdr hence not proceeding further")
		return errors.New("Subdockey is not present in the cdr hence not proceeding further")
    } 
    //Fetch the subdoc where the document keys would be stored
    subDoc := make(map[string]interface{})
    err, subDocIntf := CbCacheObj.GetSubDoc(key, subDocKey, mapAccessInfo)
    if err != nil {
        fmt.Println( "Failed to get Subdocument from couchbase")
        return errors.New("Failed to get Subdocument from couchbase")
    }
    subDoc = subDocIntf.(map[string]interface{})
    if lookupKey, ok = subDoc[action.RELEASEAGGREGATEDDOC.LookupKey].(string) ; !ok {
        fmt.Println( "LookupKey is not present in the field :" , action.RELEASEAGGREGATEDDOC.LookupKey , " in the subdoc")
        return errors.New("LookupKey is not present in the field :" + action.RELEASEAGGREGATEDDOC.LookupKey + " in the subdoc")
    }
    for key, val := range action.RELEASEAGGREGATEDDOC.Lookup {
        if strings.HasPrefix(lookupKey, key) {
            if docKey , ok := subDoc[val.CbKey].(string) ; ok {
                err, exist, doc := CbCacheObj.GetDocumentAll(docKey, mapAccessInfo)
                if err != nil {
                    fmt.Println( "Key:", docKey, "Failed to get document from couchbase")
                    return errors.New("Key: " + docKey + ", Failed to get document from couchbase")
                }
                if exist == false {
                    fmt.Println( "Key:", docKey, "Document not present in couchbase")
                    return nil
                }
                doc[CGFM_TAG_ID]= mapAccessInfo.TagId
                doc["SCTNAME"] = val.SctName
                gpp.RwLock.RLock()
                err = gpp.ExecuteServices(doc, outCdrList, val.SctName, transID)
                gpp.RwLock.RUnlock()
                if err != nil {
                    ml.MavLog(ml.ERROR, transID, "Failed to execute the services")
                    return errors.New("Failed to execute the services")
                }
            }
        }
    }
    return nil
}