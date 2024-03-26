package main

import (
	"fmt"
	"strconv"
)

func formSubDocKey(cdr map[string]interface{}, fields []string, transID string) string{
    key := ""
	delimeter := ""

    for index, field := range fields {
        if fvalue, ok := cdr[field]; ok {
            switch fvalue.(type) {
            case string:
                key += delimeter + fvalue.(string) 
            case int:
                key += delimeter + strconv.Itoa(fvalue.(int))
            case float64:
                key += delimeter + strconv.Itoa(int(fvalue.(float64)))
            default:
              //  ml.MavLog(ml.ERROR, transID, "Unsupported data type while forming CB key")
                return ""
            }
        } else {
           // ml.MavLog(ml.ERROR, transID, "No field", field, "in the cdr")
            return ""
        }
		if index == 0 {
			delimeter = "."
		}
    }
	return key
}
func main (){
	cdr := make(map[string]interface{})
	cdr["RGKey"] = "4915566102607:cmazpcmvctas01diam01.ims.mnc023.mcc262.3gppnetwork.org;1911823522;71390:5001"
	var arr []string = []string{"RGKey"}
	fmt.Println(formSubDocKey(cdr,arr,""))
}
