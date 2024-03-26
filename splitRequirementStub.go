package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/couchbase/gocb/v2"
)
type CommitCdrOnSplit struct{
	ParamList []ParamDetail	`json:paramlist`
}

// type ParamDetail struct {
//     FieldName string	`json:fieldName`
//     IsMandatary bool	`json:isMandatory`
// }

// var (
// 	collection *gocb.Collection
// 	Bucket *gocb.Bucket
// )
var inputCDR = map[string]interface{}{"@timestamp":"2022-08-18T11:03:03.924640524Z","pduSessionChargingInformation_pduSessionInformation_ratType":"5g","CBkeyseperator":":","IMEI":"imeisv-33353835303530383030323032313046","IMSI":"202020000000001","MSISDN":"14699162513","PID":"proc_pid_0","SCTNAME":"SCT_VOICE_CHF_MVNO1","SessionKey":"14699162513:af1.TMDEP-15643.mavenir.com.mavenir.com;1096298391;58081","basicService":"Voice","bssSubscriberID":"123e4567-e89b-12d3-a456-2020200001","causeForRecordClosing":"PartialCDRSessionEstablished","chargedDurationUnit":"sec","chargingSessionIdentifier":"af1.TMDEP-15643.mavenir.com.mavenir.com;1096298391;58081","direction":"O","documentID":"_af1.TMDEP-15643.mavenir.com.mavenir.com;1096298391;58081_0_2022-07-07T06:38:05+00:00_b4c1b04d-4c14-42da-ae17-ab58a7e3d8c5_SCT_VOICE_CHF_MVNO1","duration":0,"endTimeGerman":"2022-07-07T08:38:05+02:00","imsInfo_accessNetworkInformation_0":"3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=310310138f0080801","imsInfo_calledPartyAddress":"sip:001869;phone-context=mavenir.com@msg.lab.t-mobile.com;user=phone","imsInfo_callingPartyAddress_0":"sip:14699162513@mavenir.com","imsInfo_eventType_sipMethod":"INVITE","imsInfo_iMSChargingIdentifier":"ocatf000.sip.lab.t-mobile.com-1521-66590-266619","imsInfo_interOperatorIdentifier_0_originatingIOI":"32345","imsInfo_nodeFunctionality":6,"imsInfo_numberPortabilityRoutingInformation":"Routing","imsInfo_requestedPartyAddress_0":"sip:19144392231;phone-context=mavenir.com@msg.lab.t-mobile.com;user=phone","imsInfo_roleofNode":0,"imsInfo_userSessionID":"000050D7F79C-3288-1689e700-48-5aa9a25e-59ddf","initialRecordOpeningTime":"2022-07-07T06:38:05+00:00","keyseperator":"_","localRecordSequenceNumber":"43d8046b-2913-4fa1-9019-8d5d2cfa4e35::595","locationinfo":"310310138f0080801","networkProvider":"310310","networkSessionID":"ocatf000.sip.lab.t-mobile.com-1521-66590-266619","recordExtensions_extChargingInformation_callDestinationType":"CDT_INTERNATIONAL","recordExtensions_extChargingInformation_callLeg":"CL_MO","recordExtensions_extChargingInformation_calledPartyTerminalType":"CPTT_FIXED_LINE","recordExtensions_extChargingInformation_forwardedFlag":"FL_FORWARDED","recordExtensions_extChargingInformation_intlUsageZone_0":"AUCKLAND","recordExtensions_extChargingInformation_networkFlag":"NF_OFFNET","recordExtensions_extChargingInformation_roamType":"RT_NO_ROAM","recordExtensions_extChargingInformation_roamUsageZone_0":"PARIS","recordExtensions_requestType":"ONLINE_SESSION","recordExtensions_spInformation_spId":"541b778b2d524bb4a3d059fb0fbd12ee","recordExtensions_spInformation_spName":"T-Mobile","recordExtensions_subscriberInformation_billingAcctId":"77008800","recordExtensions_subscriberInformation_imsi":"202020000000001","recordExtensions_subscriberInformation_msisdn":"14699162513","recordExtensions_subscriberInformation_state":"ACTIVE","recordExtensions_subscriberInformation_subCosID":"PrepaidBaseClass","recordExtensions_subscriberInformation_subscriberId":"123e4567-e89b-12d3-a456-2020200001","recordExtensions_subscriberInformation_tariffPlan":"PKGDIA0000002","recordExtensions_subscriberInformation_taxExemption":true,"recordExtensions_subscriberInformation_type":"Prepaid","recordExtensions_transactionType":"CREATE","recordOpeningTime":"2022-07-07T06:38:05+00:00","recordSequenceNumber":0,"recordType":"D","subscriberIdentifier":"imsi-202020000000001","tariff":"PKGDIA0000002"}
var list CommitCdrOnSplit
var stuff map[string]interface{}

func main() {
	file, err := ioutil.ReadFile("sampleJsonData.json")
	if err != nil {
		fmt.Println("error occured while reading the file")
	}
	
	_ = json.Unmarshal(file, &list)
	//fmt.Println(list.ParamList[0].FieldName)
	
	opts := gocb.ClusterOptions{Username: "root",Password: "mavenir"}
	cluster, err := gocb.Connect("127.0.0.1:8091", opts)
	if err != nil {
		fmt.Println("error in connection: ", err)
	}
	if err == nil {
		fmt.Println("connection success ")
	}
	Bucket = cluster.Bucket("sessiondb_mediation")

	//Wait upto timout and checks database conn status
	if err = Bucket.WaitUntilReady(5 * time.Second, nil); err != nil{
		fmt.Println("Couchbase bootstrapError : %v", err)
	}
	col := Bucket.DefaultCollection()
	collection = col
	elem,_ := collection.Get("rgkey",nil)
	
	elem.Content(&stuff)

	 _ , ok := stuff["SPLITKEY"]
	if !ok {
		var splitKeyToInsert = make(map[string]interface{})
		for _, v := range list.ParamList {
			splitKeyToInsert[v.FieldName] = inputCDR[v.FieldName]
		}
		fmt.Println(splitKeyToInsert)
		stuff["SPLITKEY"] = splitKeyToInsert
		_, err := collection.Upsert("rgkey", &stuff, &gocb.UpsertOptions{Timeout: 3 * time.Second})
    	fmt.Println(err)
		/*
		err = dbLayer.Layer.InsertDoc(key, doc)
                if err != nil {
                    ml.MavLog(ml.ERROR, transID, "error while inserting doc to couch base for key:", key, " error:", err, " doc:", doc)
                    continue
                }
		*/
		//EXIT THIS FUNCTION NOW AS THE SPLITKEY WAS NOT PRESENT EARLIER
	}
	
	//checking cdr fields with couchbase
	for _, v := range list.ParamList {
		if inputCDR[v.FieldName] == nil && v.IsMandatary == false {
			//assume no change	
			continue
		}
		if inputCDR[v.FieldName] != nil {
			if inputCDR[v.FieldName] != stuff["SPLITKEY"].(map[string]interface{})[v.FieldName]{
				fmt.Println(v.FieldName, " <--- Its not same")
				invokeSct()
			}else{
				fmt.Println(v.FieldName, "ITS SAME")
				continue
			}
		}
	}
}

func invokeSct(){
	fmt.Println("SCT INVOKED")
	updateKey("rgkey")
}

func updateKey(key string){
	var splitKeyToInsert = make(map[string]interface{})
		for _, v := range list.ParamList {
			splitKeyToInsert[v.FieldName] = inputCDR[v.FieldName]
		}
		fmt.Println(splitKeyToInsert)
		stuff["SPLITKEY"] = splitKeyToInsert
		_, err := collection.Upsert("rgkey", &stuff, &gocb.UpsertOptions{Timeout: 3 * time.Second})
    	fmt.Println(err)
}


	//jsonStr, err := json.Marshal(stuff)
	//fmt.Println(elem)
	//fmt.Println(string(jsonStr))