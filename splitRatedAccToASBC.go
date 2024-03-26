package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type InitialNetworkDetails struct {
    ServiceRequestTime_O int64
    ServiceDeliveryStartTimeStamp_O int64
    ServiceDeliveryEndTimeStamp_O int64
    ServiceRequestTime_T int64
    ServiceDeliveryStartTimeStamp_T int64
    ServiceDeliveryEndTimeStamp_T int64
}

func calculateTotal4GWifiTime (initialNetworkInfo InitialNetworkDetails, direction string, accAccessNetworkInfo []AccessNetworkInfo, transID string) (int64) {

    var total4GWiFiTime int64 = 0
    var netType string
    var serviceDeliveryEndTime, serviceDeliveryStartTime int64
    if direction == "O"{
        serviceDeliveryEndTime= initialNetworkInfo.ServiceDeliveryEndTimeStamp_O
        serviceDeliveryStartTime= initialNetworkInfo.ServiceDeliveryStartTimeStamp_O
    }else{
        serviceDeliveryEndTime= initialNetworkInfo.ServiceDeliveryEndTimeStamp_T
        serviceDeliveryStartTime= initialNetworkInfo.ServiceDeliveryStartTimeStamp_T
    }
   fmt.Println(len(accAccessNetworkInfo), "InitialNetworkDetails", initialNetworkInfo, "accAccessNetworkInfo",accAccessNetworkInfo)
    if len(accAccessNetworkInfo)>0{
        for i,ani:=range accAccessNetworkInfo{
            if i > 0{
                if (netType == "4G" || netType == "WIFI") && ani.AccessChangeTime > serviceDeliveryStartTime {
                    duration := ani.AccessChangeTime - serviceDeliveryStartTime
                    total4GWiFiTime = (total4GWiFiTime + duration)
                    serviceDeliveryStartTime=ani.AccessChangeTime
                }
            }else{
                netType=ani.NetworkType
            }
            if i == (len(accAccessNetworkInfo)-1){
                if (ani.NetworkType == "4G" || ani.NetworkType == "WIFI") && (serviceDeliveryEndTime > ani.AccessChangeTime){
                    duration := serviceDeliveryEndTime - ani.AccessChangeTime
                    total4GWiFiTime = (total4GWiFiTime + duration)
                }
            }
        }
    }else{
        total4GWiFiTime=serviceDeliveryEndTime - serviceDeliveryStartTime
    }
    return total4GWiFiTime
}
type SplitRatedCdrAccToAsbc struct {
	SctName	string  `json:sctName"`
	ServiceKey string	`json:serviceKey"`
	CorrelationDockeys  []struct { // to find correlation key as per field value after matching key and servicekey
		Key   string  `json:"key"`
		Fields	[]string  `json:"value"`
	} `json:"correlationDockeys"`
	CorrelationDockeysMap map[string][]string
	NetTypeRegexPattern *regexp.Regexp
	NetProviderRegexPattern *regexp.Regexp
	NetworkTypePattern string	`json:networkTypePattern"`
	NetworkTypeCustomSetName string	`json:networkTypeCustomSetName"`
	NetworkProviderPattern string `json:networkProviderPattern"`
	NetworkProviderCustomSetName string	`json:networkProviderCustomSetName"`
	NetTypeGroup string	`json:netTypeGroup"`
	NetProviderGroup string	`json:netProviderGroup"`
	RtpReceivedPrecision int `json:rtpReceivedPrecision"`
	RtpSentPrecision int `json:rtpSentPrecision"`
	AvgPackateSizePrecision int `json:avgPackateSizePrecision"`
}

func actionSplitRatedCdrAccToASBC(action *SplitRatedCdrAccToAsbc, key string, cdr map[string]interface{}) error{
    
    var (
        AccessNetworkInfoArr []AccessNetworkInfo = make([]AccessNetworkInfo, 0)
        correlationDoc map[string]interface{} = make(map[string]interface{})
        err error
        exist bool
        splitOccurred bool
        partialIndicatorCounter int = 0
        networkType string
        networkProvider string
        rtpDetails RtpDetails
        totalTimeAsbc4GWifi_O int64
        totalTimeAsbc4GWifi_T int64
        initialNetworkInfo InitialNetworkDetails
        AccessNetworkInfoArr_O []AccessNetworkInfo = make([]AccessNetworkInfo, 0)
        AccessNetworkInfoArr_T []AccessNetworkInfo = make([]AccessNetworkInfo, 0)
    )
    var kTab string
	if cdr["correlationServiceKey"] == action.ServiceKey {
        correlationDoc = cdr
		if val, ok := cdr[CORRELATION_CB_KEY].(string) ; ok {
			defer deleteCorrelationDoc(val, icWrapper)
         }
    }else {
        var correlationKey string
        if val, ok := cdr["correlationServiceKey"].(string) ; ok {
            if fieldsArr, ok := action.CorrelationDockeysMap[val]; ok {
                correlationKey = formDocKey(cdr, "_", fieldsArr, icWrapper.IntransitCDRs.TransID)
                 if util.UseKtab && kTab != "" {
                     correlationKey = kTab+KEYSET_SEPERATOR+KEYSET_SEPERATOR+correlationKey
                }
            }

            err, exist, correlationDoc = icWrapper.IntransitCDRs.GetDocumentAll(correlationKey)
             if util.UseKtab && kTab != "" {
                correlationDoc[KTAB] = kTab
             }
            if err != nil {
				fmt.Println("Key:", correlationKey, "Error occured while fetching the correlation document, err:", err)          
				return errors.New("Error occured while fetching the correlation document")
			}
            if !exist {
                fmt.Println("Key:", correlationKey, "Correlation Document doesnot exist in couchbase")
                return errors.New("Key:" + correlationKey + "Correlation Document doesnot exist in couchbase")   
            }
			 defer deleteCorrelationDoc(correlationKey, icWrapper)
        }else{
            fmt.Println("correlationServiceKey is missing in the cdr")
            return errors.New("correlationServiceKey is missing in the cdr")
        }
    }
    
    defer DeleteReferenceDocs(correlationDoc, icWrapper, COMMIT_DOCS)
    defer DeleteReferenceDocs(correlationDoc, icWrapper, REFERENCE_DOCS)
    defer DeleteReferenceDocs(correlationDoc, icWrapper, DROP_DOCS)

    commitDocs, ok := correlationDoc[COMMIT_DOCS].([]interface{})
    if !ok{
        fmt.Println("No rated cdr present in the correlation document")
        return errors.New("No rated cdr present in the correlation document")    
    }
	if len(commitDocs) > 0{
        //fetching ASBC record first to avoid double iteration later incase of ONNET
        if referDocs , ok := correlationDoc["ReferenceDocs"].([]interface{}); ok && len(referDocs) > 0 {
            for _, asbcDoc := range referDocs {
                if val , ok := asbcDoc.(map[string]interface{})["DocKey"].(string) ; ok {
                    err, exist , docRes := icWrapper.IntransitCDRs.GetDocumentAll(val)
                    if err != nil {
                        fmt.Println("Key:", val, "Error occured while fetching asbcCdr, err:", err)
                        continue
                    }
                    if !exist {
                        fmt.Println("Key:", val, "Document doesnot exist in couchbase")
                        continue  
                    }
                    counter := 0	
                    var ok bool
                    initialNetworkInfo, err = GetInitialNetworkInfo(docRes, icWrapper.IntransitCDRs.TransID)
                    if docRes["pCSCFRecord_role-of-Node"] == "originating"{
                        rtpS := strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-callingPartyPacketStatistic_$i_totalRtpSent", "$i", strconv.Itoa(counter), -1)
                        if rtpSent,ok:= docRes[rtpS];ok{
                            rtpDetails.TotalRtpSent_O=rtpSent.(float64)
                            rtpDetails.isRtpDataPresent_O=true
                        }
                        rtpR := strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-callingPartyPacketStatistic_$i_totalRtpReceived", "$i", strconv.Itoa(counter), -1)
                        if rtpReceived,ok:= docRes[rtpR];ok{
                            rtpDetails.TotalRtpReceived_O=rtpReceived.(float64)
                        }
                        avgPktSize:=strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-callingPartyPacketStatistic_$i_averagePacketSize", "$i", strconv.Itoa(counter), -1)
                        if avgPacketSize,ok:= docRes[avgPktSize];ok{
                            rtpDetails.AveragePacketSize_O=avgPacketSize.(float64)
                        }
                    }else{
                        rtpS := strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-calledPartyPacketStatistic_$i_totalRtpSent", "$i", strconv.Itoa(counter), -1)
                        if rtpSent,ok:= docRes[rtpS];ok{
                            rtpDetails.TotalRtpSent_T=rtpSent.(float64)
                            rtpDetails.isRtpDataPresent_T=true
                        }
                        rtpR := strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-calledPartyPacketStatistic_$i_totalRtpReceived", "$i", strconv.Itoa(counter), -1)
                        if rtpReceived,ok:= docRes[rtpR];ok{
                            rtpDetails.TotalRtpReceived_T=rtpReceived.(float64)
                        }
                        avgPktSize:=strings.Replace("pCSCFRecord_recordExtensions_$i_list-of-calledPartyPacketStatistic_$i_averagePacketSize", "$i", strconv.Itoa(counter), -1)
                        if avgPacketSize,ok:= docRes[avgPktSize];ok{
                            rtpDetails.AveragePacketSize_T=avgPacketSize.(float64)
                        }
                    }
                    for {
                        var networkInfo AccessNetworkInfo
                        plmn := strings.Replace("pCSCFRecord_list-Of-AccessNetworkInfoChange_$i_accessNetworkInformation", "$i", strconv.Itoa(counter), -1)
                        networkInfo.AccessNetworkInformation , ok = docRes[plmn].(string)
                        if ok {
                            var networkTypeTemp, networkProviderTemp  string
                            if action.NetTypeRegexPattern.MatchString(networkInfo.AccessNetworkInformation) {
                                netTypeTemp := action.NetTypeRegexPattern.ReplaceAllString(networkInfo.AccessNetworkInformation, action.NetTypeGroup)
                                err, _ , doc := customSet.GetMatch(action.NetworkTypeCustomSetName, netTypeTemp)
                                if err != nil {
                                    ml.MavLog(ml.WARN, icWrapper.IntransitCDRs.TransID, "error from custom set GetMatch for setName:", action.NetworkTypeCustomSetName, " and cdr Field:",
                                    netTypeTemp, " error:", err)
                                }
                                if doc != nil {
                                    networkTypeTemp = doc.Value
                                    networkInfo.NetworkType=networkTypeTemp
                                }
                            }
                            if action.NetProviderRegexPattern.MatchString(networkInfo.AccessNetworkInformation){
                                netProviderTemp := action.NetProviderRegexPattern.ReplaceAllString(networkInfo.AccessNetworkInformation, action.NetProviderGroup)
                                err, _ , doc := customSet.GetMatch(action.NetworkProviderCustomSetName, netProviderTemp)
                                if err != nil {
                                    ml.MavLog(ml.WARN, icWrapper.IntransitCDRs.TransID, "error from custom set GetMatch for setName:", action.NetworkProviderCustomSetName, " and cdr Field:",
                                    netProviderTemp, " error:", err)
                                }
                                if doc != nil {
                                    networkProviderTemp = doc.Value
                                } else {
                                    if len(netProviderTemp) > 6 {
                                        networkProviderTemp = netProviderTemp[0:6]
                                    }
                                }
                            }
                            if counter != 0 && ((networkTypeTemp == networkType && networkProviderTemp == networkProvider) || (networkTypeTemp != "WIFI" && networkProviderTemp == "")) {
                                counter++
                                continue
                            }
                            timeStamp := strings.Replace("pCSCFRecord_list-Of-AccessNetworkInfoChange_$i_accessChangeTime", "$i", strconv.Itoa(counter), -1)
                            AccessChangeTime , ok := docRes[timeStamp].(string)
                            if !ok {
                                fmt.Println("AccessChangeTime is absent")
                                break
                            }
                            val := convertTimeToEpoch(AccessChangeTime, icWrapper.IntransitCDRs.TransID)
                            if val == 0 {
                                fmt.Println("Error occured while converting time to epoch, err is", err)
                                counter++
                                continue 
                            }
                            networkInfo.AccessChangeTime = val
                            if docRes["pCSCFRecord_role-of-Node"] == "originating"{
                                networkInfo.Direction = "O"
                                AccessNetworkInfoArr_O=append(AccessNetworkInfoArr_O, networkInfo)
                            }else{
                                networkInfo.Direction = "T"
                                AccessNetworkInfoArr_T=append(AccessNetworkInfoArr_T, networkInfo)
                            }
                            AccessNetworkInfoArr = append(AccessNetworkInfoArr, networkInfo)
                            networkType = networkTypeTemp
                            networkProvider = networkProviderTemp
                        }else{
                            if counter == 0 {
                            	fmt.Println("Access network change info not exist in cdr")
                            } else {
                                fmt.Println("No more access network change info in the CDR")
                            }
                            break
                        }
                        counter++
                    }
                }else{
                    fmt.Println("ASBC CDR's document key is not present , hence document not accessible")
                    for _, doc := range commitDocs {
                        if val, ok := doc.(map[string]interface{})["DocKey"].(string); ok {
                            err, exist , cdr := icWrapper.IntransitCDRs.GetDocumentAll(val)
                            if err != nil {
                                fmt.Println("Key:", val, "Error occured while fetching the Rated CDR")
                                continue
                            }
                            if !exist {
                                fmt.Println("Rated CDR with key ",val," doesnot exist in couchbase")
                                continue  
                            }
							icWrapper.GppInfo.Append("InvokingNewSct", "", "==>")
                            err = gpp.ExecuteServices(cdr, icWrapper, action.SctName, icWrapper.GppInfo, icWrapper.IntransitCDRs.TransID)
							icWrapper.GppInfo.Append("", action.SctName, "<==END|")
                            if err != nil {
                                fmt.Println("failed to execute the services")
                            }
                        }
                    }
                    return errors.New("ASBC CDR's document key is not present , hence document not accessible")
                }
            }
        } else {
            // Flush T/S cdr when no ASBC cdr is received with NetworkChangeInfo
            fmt.Println("Flushing T/S cdr as no ASBC cdr is received")
            for _, doc := range commitDocs {
                if val, ok := doc.(map[string]interface{})["DocKey"].(string); ok {
                    err, exist , cdr := icWrapper.IntransitCDRs.GetDocumentAll(val)
                    if err != nil {
                        fmt.Println("Key:", val, "Error occured while fetching the Rated CDR")
                        continue
                    }
                    if !exist {
                        fmt.Println("Rated CDR with key ",val," doesnot exist in couchbase")
                        continue  
                    }
					icWrapper.GppInfo.Append("InvokingNewSct", "", "==>")
                    err = gpp.ExecuteServices(cdr, icWrapper, action.SctName, icWrapper.GppInfo, icWrapper.IntransitCDRs.TransID)
					icWrapper.GppInfo.Append("", action.SctName, "<==END|")
                    if err != nil {
                        fmt.Println("failed to execute the services")
                    }
                }
            }
            return errors.New("Flushing T/S cdr as no ASBC cdr is received")
        }
        lenOfAccessNetworkArr :=  len(AccessNetworkInfoArr)
        if lenOfAccessNetworkArr > 0 {
            sortBasedOnTime(AccessNetworkInfoArr , icWrapper.IntransitCDRs.TransID) 
            totalTimeAsbc4GWifi_O = calculateTotal4GWifiTime(initialNetworkInfo , "O", AccessNetworkInfoArr_O, icWrapper.IntransitCDRs.TransID)
            totalTimeAsbc4GWifi_T = calculateTotal4GWifiTime(initialNetworkInfo , "T", AccessNetworkInfoArr_T, icWrapper.IntransitCDRs.TransID)
			fmt.Println("Total time for 4G and WIFI for for ASBC_O :", totalTimeAsbc4GWifi_O, " and for ASBC_T : ", totalTimeAsbc4GWifi_T)
        } else {
           // Flush T/S cdr when no ASBC cdr is received with NetworkChangeInfo
           fmt.Println("Flushing T/S cdr as NetworkChangeInfo is not present")
           for _, doc := range commitDocs {
                if val, ok := doc.(map[string]interface{})["DocKey"].(string); ok {
                    err, exist , cdr := icWrapper.IntransitCDRs.GetDocumentAll(val)
                    if err != nil {
                        fmt.Println("Key:", val, "Error occured while fetching the Rated CDR")
                        continue
                    }
                    if !exist {
                        fmt.Println("Rated CDR with key ",val," doesnot exist in couchbase")
                        continue  
                    }

                    if netType,ok:=cdr["networkType"].(string);ok{
                        if netType == "4G" || netType == "WIFI"{
                            if direction,ok:=cdr["direction"].(string);ok{
                                if direction == "O" && rtpDetails.isRtpDataPresent_O {
                                    cdr["totalRtpSent"]=rtpDetails.TotalRtpSent_O
                                    cdr["totalRtpReceived"]=rtpDetails.TotalRtpReceived_O
                                    cdr["averagePacketSize"]=rtpDetails.AveragePacketSize_O
                                    cdr["totalRtpSent_precision"]=action.RtpSentPrecision
                                    cdr["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                    cdr["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                                }else if direction == "T" && rtpDetails.isRtpDataPresent_T {
                                    cdr["totalRtpSent"]=rtpDetails.TotalRtpSent_T
                                    cdr["totalRtpReceived"]=rtpDetails.TotalRtpReceived_T
                                    cdr["averagePacketSize"]=rtpDetails.AveragePacketSize_T
                                    cdr["totalRtpSent_precision"]=action.RtpSentPrecision
                                    cdr["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                    cdr["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                                }
                            }
                        }
                    }
					// icWrapper.GppInfo.Append("InvokingNewSct", "", "==>")
                    // err = gpp.ExecuteServices(cdr, icWrapper, action.SctName, icWrapper.GppInfo, icWrapper.IntransitCDRs.TransID)
					// icWrapper.GppInfo.Append("", action.SctName, "<==END|")
                    // if err != nil {
                    //     fmt.Println("failed to execute the services")
                    // }
                }
            }
           return errors.New("Flushing T/S cdr as NetworkChangeInfo is not present") 
        }
		//traverse each RatedCDR and split based on ASBC
		for _, doc := range commitDocs {
            var netType string
            var aggrDuration int = 0
            var partNumberCounter int = 1
            if val, ok := doc.(map[string]interface{})["DocKey"].(string); ok{
                err, exist , cdr := GetDocumentAll(val)
                if err != nil {
                    fmt.Println("Key:", val, "Error occured while fetching the Rated CDR")
                    continue
                }
                if !exist {
                    fmt.Println("Rated CDR with key ",val," doesnot exist in couchbase")
                    continue  
                }
                var outCdrListAsbc []map[string]interface{} = make([]map[string]interface{}, 0)
                var orgCdrStartTimeGermanInEpoch int64
                if startTime, ok := cdr["startTimeGerman"].(string); ok{
                    orgCdrStartTimeGermanInEpoch = convertTimeToEpoch(startTime, icWrapper.IntransitCDRs.TransID)
                    if orgCdrStartTimeGermanInEpoch == 0{
                        fmt.Println("Error occured while converting time to epoch")
                        continue
                    }
                    //cdr["startTimeGerman"] = val
                } else{
                    fmt.Println("StartTime is not present in the rated CDR")
                    //return errors.New("StartTime is not present in the rated CDR")
                    continue
                }
				newAccessNetworkInfoArr := delete_empty(AccessNetworkInfoArr)
                lenOfAccessNetworkInfoArr := len(newAccessNetworkInfoArr) 
                if netTypeval,ok:=cdr["networkType"].(string);ok{
                    netType=netTypeval
                }
                var i int = 0
                originalCDR := CopyableMap(cdr).DeepCopy()
                for i < lenOfAccessNetworkInfoArr {
					if newAccessNetworkInfoArr[i].AccessChangeTime <= orgCdrStartTimeGermanInEpoch {
                        fmt.Println("Current networkchange time is before or equal to startTimeGerman hence ignoring")
                        i++
                        continue
                    }
                    // copying only for the first iteration as after subsequent iteration
                    // endtime becomes start time for next splitted cdr
                    if i == 0 {
                        cdr["startTimeGerman"] = orgCdrStartTimeGermanInEpoch
                    }
					var duration int64
                    if val, ok := cdr["startTimeGerman"].(int64) ; ok{
                        if newAccessNetworkInfoArr[i].AccessChangeTime == val {
                            fmt.Println("Current networkchange time is equal to the startTimeGerman hence ignoring")
                            i++
                            continue
                        }
                        endTime := newAccessNetworkInfoArr[i].AccessChangeTime
                        duration = endTime-val
                        cdr["aggrduration"] = duration
                        aggrDuration = aggrDuration + int(duration)
                        if partialIndicatorCounter == 0 {
                            cdr["partialIndicator"] = "F"
                            cdr["partNumber"] = partNumberCounter
                            partialIndicatorCounter++
                            partNumberCounter++
                        }else{
                            cdr["partialIndicator"] = "I"
                            cdr["partNumber"] = partNumberCounter
                            partialIndicatorCounter++
                            partNumberCounter++
                        }
                        if netType == "4G" || netType == "WIFI"{
                            if dir,ok:=cdr["direction"].(string);ok{
                                if dir == "O" && rtpDetails.isRtpDataPresent_O || dir == "T" && rtpDetails.isRtpDataPresent_T{
                                    rtpSent, rtpRec, avgPackateSize :=CalculateRTPData(rtpDetails,dir,icWrapper.IntransitCDRs.TransID,float64(duration),totalTimeAsbc4GWifi_O,totalTimeAsbc4GWifi_T)
                                    cdr["totalRtpSent"]=rtpSent
                                    cdr["totalRtpReceived"]=rtpRec
                                    cdr["averagePacketSize"]=avgPackateSize
                                    cdr["totalRtpSent_precision"]=action.RtpSentPrecision
                                    cdr["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                    cdr["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                                }
                            }
                        }else{
                            delete(cdr,"totalRtpSent")
                            delete(cdr,"totalRtpReceived")
                            delete(cdr,"averagePacketSize")
                            delete(cdr,"totalRtpSent_precision")
                            delete(cdr,"totalRtpReceived_precision")
                            delete(cdr,"averagePacketSize_precision")
                        }
                        splitOccurred = true
                        cdr["asbcSplit"] = true
                        spilttedCDR := CopyableMap(cdr).DeepCopy()
                        outCdrListAsbc = append(outCdrListAsbc, spilttedCDR)
                        cdr["startTimeGerman"] = endTime
                        if strings.HasSuffix(doc.(map[string]interface{})["ServiceType"].(string), "O"){
                            if newAccessNetworkInfoArr[i].Direction == "O"{
                                cdr["networkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
                                cdr["networkProviderFlag"] = true
                                netType=newAccessNetworkInfoArr[i].NetworkType
                            }else if newAccessNetworkInfoArr[i].Direction == "T"{
                                cdr["secondPartyNetworkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
                                cdr["secondPartyNetworkProviderFlag"] = true
                            }
                        } else {
                            if newAccessNetworkInfoArr[i].Direction == "T"{
                                cdr["networkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
                                cdr["networkProviderFlag"] = true
                                netType=newAccessNetworkInfoArr[i].NetworkType
                            }else if newAccessNetworkInfoArr[i].Direction == "O"{
                                cdr["secondPartyNetworkProvider"] = newAccessNetworkInfoArr[i].AccessNetworkInformation
                                cdr["secondPartyNetworkProviderFlag"] = true
                            }
                        }
                    }else{
                        fmt.Println("StartTime is missing in the CDR")
                        // Should We exit if startTime is not present in the RATED CDR
						icWrapper.GppInfo.Append("InvokingNewSct", "", "==>")
                        err = gpp.ExecuteServices(cdr, icWrapper, action.SctName, icWrapper.GppInfo, icWrapper.IntransitCDRs.TransID)
						icWrapper.GppInfo.Append("", action.SctName, "<==END|")
                        if err != nil {
                            fmt.Println("failed to execute the services")
                            return errors.New("failed to execute the services")
                        }
                        return errors.New("StartTime is missing in the CDR")
                    }
                    i++
                }
                if partialIndicatorCounter >= 1 {
                    var duration int
                    startTime :=   cdr["startTimeGerman"].(int64)
                    if val, ok := originalCDR["aggrduration"].(float64) ; ok {
                        cdr["aggrduration"] = int(val) - aggrDuration
                        duration= int(val) - aggrDuration
                    }
                    cdr["partialIndicator"] = "L"
                    cdr["partNumber"] = partNumberCounter
                    cdr["startTimeGerman"] = startTime
                    if netType == "4G" || netType == "WIFI"{
                        if dir,ok:=cdr["direction"].(string);ok{
                            if dir == "O" && rtpDetails.isRtpDataPresent_O || dir == "T" && rtpDetails.isRtpDataPresent_T{
                                rtpSent, rtpRec, avgPackateSize :=CalculateRTPData(rtpDetails,dir,icWrapper.IntransitCDRs.TransID,float64(duration),totalTimeAsbc4GWifi_O,totalTimeAsbc4GWifi_T)
                                cdr["totalRtpSent"]=rtpSent
                                cdr["totalRtpReceived"]=rtpRec
                                cdr["averagePacketSize"]=avgPackateSize
                                cdr["totalRtpSent_precision"]=action.RtpSentPrecision
                                cdr["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                cdr["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                            }
                        }
                    }else{
                        delete(cdr,"totalRtpSent")
                        delete(cdr,"totalRtpReceived")
                        delete(cdr,"averagePacketSize")
                        delete(cdr,"totalRtpSent_precision")
                        delete(cdr,"totalRtpReceived_precision")
                        delete(cdr,"averagePacketSize_precision")
                    }
                    splitOccurred = true
                    cdr["asbcSplit"] = true
                    spilttedCDR := CopyableMap(cdr).DeepCopy()
                    outCdrListAsbc = append(outCdrListAsbc, spilttedCDR)
                }
                if splitOccurred {
                    originalCDR["partialIndicator"] = "T"
                    originalCDR["partNumber"] = partNumberCounter
                    originalCDR["asbcSplit"] = true
                    if endTime , ok := originalCDR["endTimeGerman"].(int64) ; ok && endTime != 0 {
                        if startTime, ok := originalCDR["startTimeGerman"].(int64) ; ok && startTime != 0 {
                            totalDuration := endTime - startTime
                            originalCDR["asbcDuration"] = totalDuration
                            if val, ok := originalCDR["aggrduration"].(float64) ; ok{
                                if totalDuration != int64(val) {
                                    fmt.Println("Duration in rated cdr and asbc duration are not same , RatedDuration :", originalCDR["duration"] , " asbcDuration :", originalCDR["asbcDuration"])
                                }
                            }
                        }
                    }
                }
                if netType,ok:=originalCDR["networkType"].(string);ok{
                    if netType == "4G" || netType == "WIFI"{
                        if direction,ok:=originalCDR["direction"].(string);ok{
                            if direction == "O" && rtpDetails.isRtpDataPresent_O {
                                originalCDR["totalRtpSent"]=rtpDetails.TotalRtpSent_O
                                originalCDR["totalRtpReceived"]=rtpDetails.TotalRtpReceived_O
                                originalCDR["averagePacketSize"]=rtpDetails.AveragePacketSize_O
                                originalCDR["totalRtpSent_precision"]=action.RtpSentPrecision
                                originalCDR["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                originalCDR["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                            }else if direction == "T" && rtpDetails.isRtpDataPresent_T {
                                originalCDR["totalRtpSent"]=rtpDetails.TotalRtpSent_T
                                originalCDR["totalRtpReceived"]=rtpDetails.TotalRtpReceived_T
                                originalCDR["averagePacketSize"]=rtpDetails.AveragePacketSize_T
                                originalCDR["totalRtpSent_precision"]=action.RtpSentPrecision
                                originalCDR["totalRtpReceived_precision"]=action.RtpReceivedPrecision
                                originalCDR["averagePacketSize_precision"]=action.AvgPackateSizePrecision
                            }
                        }
                    }
                }
                outCdrListAsbc = append(outCdrListAsbc, originalCDR)
                for _, cdr := range outCdrListAsbc {
                    if cdr == nil{
                        continue
                    }
					// icWrapper.GppInfo.Append("InvokingNewSct", "", "==>")
                    // err = ExecuteServices(cdr, icWrapper, action.SctName, icWrapper.GppInfo, icWrapper.IntransitCDRs.TransID)
					// icWrapper.GppInfo.Append("", action.SctName, "<==END|")
                    // if err != nil {
                    //     fmt.Println("failed to execute the services, invoking sct for next splitted cdr")
                    // }
                }
            }
            partialIndicatorCounter = 0
		}
	}else{
		fmt.Println("Rated CDR details are absent")
        return errors.New("Rated CDR details are absent")
	}
    return nil
}

func convertTimeToEpoch(inputTime, transID string) int64 {
	if strings.ContainsAny(inputTime,"-T") {
		timeFormat := "%YEAR%-%month_nz%-%date_z%T%hour_24%:%minute_z%:%second_z%-07:00"
		actualFormatedTime := GetDynamicCNotationDateTimeSubstutions(timeFormat)
		tm, err := time.Parse(actualFormatedTime, inputTime)
		if err != nil {
			fmt.Println("Time:", actualFormatedTime, "Error occured while parsing the time")
            return 0
		}
		return tm.Unix()
	} else {
		var charSign string
		var year , month , day , hour , minute , second , sign, zoneh, zonem, final string
		if len(inputTime) > 0 {
			year  = string(inputTime[1]) + string(inputTime[0])
			month = string(inputTime[3]) + string(inputTime[2])
			day   = string(inputTime[5]) + string(inputTime[4])
			hour  = string(inputTime[7]) + string(inputTime[6])
			minute = string(inputTime[9]) + string(inputTime[8])
			second = string(inputTime[11]) + string(inputTime[10])
			sign = string(inputTime[12]) + string(inputTime[13])
			zoneh = string(inputTime[15]) + string(inputTime[14])
			zonem = string(inputTime[17]) + string(inputTime[16])
		}
		if strings.EqualFold(sign,"2B") {
			charSign = "+"
		} else if strings.EqualFold(sign, "2D") {
			charSign = "-"
		} else {
			fmt.Println("Time:", inputTime, "Invalid sign")
		}
		final = "20" + year + "-" + month + "-" + day + "T" + hour + ":" + minute + ":" + second +
			charSign + zoneh + ":" + zonem
		// Parse the date string
		t, err := time.Parse(time.RFC3339Nano, final)
		if err != nil {
			fmt.Println("Time:", inputTime, "failed to parse timestamp")
            return 0
		}
		return t.Unix()
	}
}

func sortBasedOnTime(accessNetworkInfoArr []AccessNetworkInfo , transID string){ 
	fmt.Println(accessNetworkInfoArr)
	sort.SliceStable(accessNetworkInfoArr, func(i, j int) bool {
		return accessNetworkInfoArr[i].AccessChangeTime < accessNetworkInfoArr[j].AccessChangeTime
	})
	fmt.Println(accessNetworkInfoArr)
}
type AccessNetworkInfo struct{
	AccessNetworkInformation string		`json:"accessNetworkInformation"`
	AccessChangeTime int64		`json:"accessChangeTime"`
	Direction string	
    NetworkType string
}

type RtpDetails struct {
    TotalRtpSent_O float64
    TotalRtpReceived_O float64
    AveragePacketSize_O float64
    isRtpDataPresent_O bool
    TotalRtpSent_T float64
    TotalRtpReceived_T float64
    AveragePacketSize_T float64
    isRtpDataPresent_T bool
}

func CalculateRTPData (rtpDetails RtpDetails, direction, transID string, duration float64, totalTimeAsbc4GWifi_O,totalTimeAsbc4GWifi_T  int64,) (float64,float64,float64){
    var totalRtpSent float64
    var totalRtpReceived float64
    var averagePacketSize float64
    if direction == "O"{
        if totalTimeAsbc4GWifi_O != 0{
            fmt.Println("Setting rtp data", rtpDetails.TotalRtpSent_O,duration,totalTimeAsbc4GWifi_O)
            totalRtpSent,_ = decimal.NewFromFloat((rtpDetails.TotalRtpSent_O) * duration / float64(totalTimeAsbc4GWifi_O)).Round(2).Float64()
            totalRtpReceived,_ = decimal.NewFromFloat((rtpDetails.TotalRtpReceived_O) * duration / float64(totalTimeAsbc4GWifi_O)).Round(2).Float64()
        }
        averagePacketSize =rtpDetails.AveragePacketSize_O
    }else{
        if totalTimeAsbc4GWifi_T != 0{
            fmt.Println("Setting rtp data", rtpDetails.TotalRtpSent_T,duration,totalTimeAsbc4GWifi_T)
            totalRtpSent,_ = decimal.NewFromFloat((rtpDetails.TotalRtpSent_T) * duration / float64(totalTimeAsbc4GWifi_T)).Round(2).Float64()
            totalRtpReceived,_ = decimal.NewFromFloat((rtpDetails.TotalRtpReceived_T) * duration / float64(totalTimeAsbc4GWifi_T)).Round(2).Float64()
        }
        averagePacketSize =rtpDetails.AveragePacketSize_T
    }
    return totalRtpSent,totalRtpReceived,averagePacketSize
}

func deleteCorrelationDoc(correlationKey string) error{
	err := DeleteDoc(correlationKey)
	if err != nil {
		fmt.Println("Failed to delete the document", correlationKey)
		return err
	}
	return nil
}

func main(){

}