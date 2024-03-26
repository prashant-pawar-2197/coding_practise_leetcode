package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/*
map[PREDocs:map[4.11214502e+08:[4915025246211:prashantpawar_pr_PRE_2]]
SCTNAME:SCT_DOC_RESTORATION_DATA
SessionKey:4915025246211:prashantpawar_pr
TIMER_ID_FIELD_V2:map[Restoration_Timer:map[timerData:CGFM_CDC_1_4915025246211:prashantpawar_pr|SCT_DOC_RESTORATION_DATA|Restoration_Timer|prashantpawar_pr|false|false|1699524548
timerID:2]]
*/
func main() {
premap := make(map[string][]string)
keyArr := make([]string, 0)
keyArr = append(keyArr, "4915025246211:prashantpawar_pr_PRE_2")
str := strconv.FormatInt(411214502, 10)
premap[str] = keyArr
restorMap := make(map[string]interface{})
restorMap["PREDocs"] = premap
restorMap["SCTNAME"] = "SCT_DOC_RESTORATION_DATA"
restorMap["SessionKey"] = "4915025246211:prashantpawar_pr"
byteArr , err := json.Marshal(restorMap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(byteArr))
	}
}