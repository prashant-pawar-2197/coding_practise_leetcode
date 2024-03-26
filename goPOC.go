package main

import (
	//"encoding/json"
	"fmt"
	"strings"
	//	"reflect"
	//"strconv"
	//"strings"
	//"time"
	//	ck "com.Go/checkpack"
	//"reflect"
)

type Person struct {
	name            string
	age             int
	email           int
	favouriteColors []string // non-comparable field
	details 		Details
}

type Details struct {
	gender			string
}
func printt(s string){
	fmt.Println(s)
}

// func deleteElement(arr []interface{}, elem interface{}) []interface{} {
// 	var newArr []interface{}
// 	for _, value := range arr {
// 		if !reflect.DeepEqual(value, elem) {
// 			newArr = append(newArr, value)
// 		}
// 	}
// 	return newArr
// }

type EnumDocType		int
type EnumRestoreMsgType int

const (
	CGFM_INVALID_DOC_TYPE EnumDocType = iota
	CGFM_SESSION
	CGFM_RESTORE
)

type DocWrapper struct {
	Type 	EnumDocType `json:"type,omitempty"`
	Key		string 		`json:"key,omitempty"`
}

type CancelTimerRequest struct {
	//SessionKey - Can be multiple in case of data
	//DataIC - Session
	//Data - Session
	DocInfo		[]DocWrapper `json:"docInfo,omitempty"`
}

const (
	CGFM_INVALID_TYPE EnumRestoreMsgType = iota
	CGFM_CANCEL_TIMER_REQ
	CGFM_CANCEL_TIMER_RESP
	CGFM_CANCEL_TIMER_AND_PULL_DOCS_REQ
	CGFM_CANCEL_TIMER_AND_PULL_DOCS_RESP
)

type RestoreReqMsg struct {
	MsgType 	EnumRestoreMsgType 	`json:"msgType,omitempty"`
	Msg			[]byte				`json:"msg,omitempty"`
}

func main() {
	/*
	l1 := DocWrapper{Type: CGFM_SESSION, Key: "kkk"}
	list := make([]DocWrapper, 0)
	list = append(list, l1)
	cancelReq := CancelTimerRequest{DocInfo: list}
	
	dataInByte, _ := json.Marshal(cancelReq)
	req := RestoreReqMsg{MsgType: CGFM_CANCEL_TIMER_REQ, Msg: dataInByte}
	newdatabyte, _ := json.Marshal(req)

	fmt.Println(string(newdatabyte))
	var newRequest RestoreReqMsg
	err := json.Unmarshal(newdatabyte, &newRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(newRequest)
	}
	var CancelReq CancelTimerRequest
	err = json.Unmarshal(newRequest.Msg, &CancelReq)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(CancelReq)
	}
	*/

	fieldVal := "SESSION::4915566580775:cmazpcmvctas01diam02.ims.mnc023.mcc262.3gppnetwork.org;781616307;460764_222222222"
	key := fieldVal[(strings.IndexByte(fieldVal, ':') + 2):]
	//keyList := strings.Split(fieldVal, "::")
	//for _, key := range keyList {
		fmt.Println(key)
	//}

	tempDoc := make(map[string]interface{})
		if tempDoc == nil {
			fmt.Println("Temp doc is empty")
		} else {
			fmt.Println("Not empty")
		}
	
}
