package main

import (
	"encoding/json"
//	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type CustomFunctionCodec struct {
	CustomSetName    string `json:"customSetName"`
	DestinationField string `json:"destinationField"`
	PatternMediaName string `json:"patternMediaName"` //for RTP/AVP
	PatternMediaDesc string `json:"patternMediaDesc"` // for rtpmap:
	IsbcFlag bool `json:"isbcFlag"`
	MediaNameField	string `json:"mediaNameField"`
	DescriptionField	string `json:"descriptionField"`
	MediaNameField_ISBC	string `json:"mediaNameField_ISBC"`
	DescriptionField_ISBC	string `json:"descriptionField_ISBC"`
}

// CustomFunction to Decode Ims Codec for Voice CDR
func actionCustomFunctionCodec(action CustomFunctionCodec, cdr map[string]interface{}, transID string) error {
	//var customSetName string = action.CustomSetName
	var destCdrField string = action.DestinationField
	var patternMediaName string = action.PatternMediaName
	var patternMediaDesc string = action.PatternMediaDesc

	counter := 0
	counterA := 0
	// Regex format sample --> "(m=audio) (.*) (RTP/AVP) (\\d+) (.*)"
	re, _ := regexp.Compile(patternMediaName)

	if action.IsbcFlag {
		for {
			CdrFieldName := strings.Replace(action.MediaNameField_ISBC, "$x", strconv.Itoa(counter), -1)
			CdrFieldName = strings.Replace(CdrFieldName, "$y", strconv.Itoa(counterA), -1)
			counter++
			val, CdrFieldAvailble := cdr[CdrFieldName]
				if !CdrFieldAvailble {
					fmt.Println("No sdpmedianame present that contains m=audio and CdrFieldName: ",CdrFieldName)
					break
				}else{
					mediaName, ok := val.(string)
					if !ok {
						continue
					}
				if re.MatchString(mediaName) {
					codecValInStr := re.ReplaceAllString(mediaName, "$4")
					codecVal, _ := strconv.Atoi(codecValInStr)
					if codecVal < 96 {
						fmt.Println("CODEC derived from customset")
						// err, found, doc := customSet.GetMatch(customSetName, codecValInStr)
						// if err != nil {
						// 	fmt.Println("CustomSet not configured properly")
						// 	return errors.New("CustomSet not configured properly")
						// }
						// if !found {
						// 	fmt.Println("No matching value present in the look up, moving to next iteration")	
						// 	continue
						// }
						// cdr[destCdrField] = doc.Value
						// continue
					}else{
						//Here we are replacing imsInfo_sDPMediaComponent_0_sDPMediaName with imsInfo_sDPMediaComponent_0_sDPMediaDescription_$i
						sDPMediaDescriptionArr := strings.Replace(CdrFieldName, "sDP-Media-Name", action.DescriptionField_ISBC, -1)
						patternForMediaDesc := strings.Replace(patternMediaDesc, "$i", codecValInStr, -1)
						// Regex format sample --> "(a=rtpmap:)" + codecValInStr + "(.*)"
						regEx := regexp.MustCompile(patternForMediaDesc)
						secondCounter := 0
						for {
							sDPMediaDescriptionEle := strings.Replace(sDPMediaDescriptionArr, "$i", strconv.Itoa(secondCounter), -1)
							val, sDPMediaDescriptionEleAvailable := cdr[sDPMediaDescriptionEle]
		
							if !sDPMediaDescriptionEleAvailable {
								fmt.Println("No sdpmediadescription present which contains rtpmap, checking in next instance")
								break
							} else {
								if regEx.MatchString(val.(string)) {
									valueFromsDPMediaDescriptionEle := regEx.ReplaceAllString(val.(string), "$2")
									cdr[destCdrField] = valueFromsDPMediaDescriptionEle
									break
								}
							}
							secondCounter++
						}
					}
				}
			}
		}
		return nil
	}else{
		//Iterating through imsInfo_sDPMediaComponent_$i_sDPMediaName
		for {
			CdrFieldName := strings.Replace(action.MediaNameField, "$i", strconv.Itoa(counter), -1)
			counter++
			val, CdrFieldAvailble := cdr[CdrFieldName]
			if !CdrFieldAvailble {
				fmt.Println("No sdpmedianame present that contains m=audio and CdrFieldName: ",CdrFieldName)
				break
			} else {
				//if val is not a string then
				mediaName, ok := val.(string)
				if !ok {
					continue
				}
				if re.MatchString(mediaName) {
					codecValInStr := re.ReplaceAllString(mediaName, "$4")
					codecVal, _ := strconv.Atoi(codecValInStr)
					if codecVal < 96 {
						fmt.Println("CODEC derived from customset")
						// err, found, doc := customSet.GetMatch(customSetName, codecValInStr)
						// if err != nil {
						// 	fmt.Println("CustomSet not configured properly")
						// 	return errors.New("CustomSet not configured properly")
						// }
						// if !found {
						// 	fmt.Println("No matching value present in the look up")
						// 	return errors.New("No matching value present in the look up")
						// }
						// cdr[destCdrField] = doc.Value
						break

					} else {
						//Here we are replacing imsInfo_sDPMediaComponent_0_sDPMediaName with imsInfo_sDPMediaComponent_0_sDPMediaDescription_$i
						sDPMediaDescriptionArr := strings.Replace(CdrFieldName, "sDPMediaName", action.DescriptionField, -1)
						patternForMediaDesc := strings.Replace(patternMediaDesc, "$i", codecValInStr, -1)
						// Regex format sample --> "(a=rtpmap:)" + codecValInStr + "(.*)"
						regEx := regexp.MustCompile(patternForMediaDesc)
						secondCounter := 0
						for {
							sDPMediaDescriptionEle := strings.Replace(sDPMediaDescriptionArr, "$i", strconv.Itoa(secondCounter), -1)
							val, sDPMediaDescriptionEleAvailable := cdr[sDPMediaDescriptionEle]

							if !sDPMediaDescriptionEleAvailable {
								fmt.Println("No sdpmediadescription present which contains rtpmap:")
								break
							} else {
								if regEx.MatchString(val.(string)) {
									valueFromsDPMediaDescriptionEle := regEx.ReplaceAllString(val.(string), "$2")
									cdr[destCdrField] = valueFromsDPMediaDescriptionEle
									return nil
								}
							}
							secondCounter++
						}
					}
				}
			}
		}	
	}
	return nil
}

func main()  {
	cdr := make(map[string]interface{})
	var action CustomFunctionCodec
	readFileData, err := ioutil.ReadFile("codecSampleCDR.json")
	if err != nil {
		fmt.Println("Error occured")
	}
	readActionData, err := ioutil.ReadFile("decodingCodecAction.json")
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
	actionCustomFunctionCodec(action,cdr,"aaa")
	encodedData,_ := json.Marshal(cdr)
 	ioutil.WriteFile("CodecOutput",encodedData, 0644)


	// if cdr["iBCFRecord_list-Of-SDP-Media-Components_0_sIP-Response-Timestamp-Fraction"].(float64) == 715{
	// 	fmt.Println("Worked without typecasting")
	// }else {
	// 	fmt.Println("Didn't work without typecasting")
	// }
	// // "(m=audio) (.*) (RTP/AVP) (\\d+) (.*)"
	// // "(a=rtpmap:$i) (.*)"
	// //isbcFlag := true
	// counter := 0
	// counterA := 0
	// patternMediaDesc := "(a=rtpmap:$i) (.*)"
	// destCdrField := "codec"
	// re, _ := regexp.Compile("(m=audio) (.*) (RTP/AVP) (\\d+) (.*)")
	// if _ , ok := cdr["iBCFRecord_list-Of-SDP-Media-Components_0_sDP-Media-Components_0_sDP-Media-Name"] ; ok {
	// 	for {
	// 		CdrFieldName := strings.Replace("iBCFRecord_list-Of-SDP-Media-Components_$x_sDP-Media-Components_$y_sDP-Media-Name", "$x", strconv.Itoa(counter), -1)
	// 		CdrFieldName = strings.Replace(CdrFieldName, "$y", strconv.Itoa(counterA), -1)
	// 			counter++
	// 			val, CdrFieldAvailble := cdr[CdrFieldName]
	// 				if !CdrFieldAvailble{
	// 					fmt.Println("No sdpmedianame present that contains m=audio")
	// 					break
	// 				}else{
	// 					mediaName, ok := val.(string)
	// 					if !ok {
	// 						continue
	// 					}
	// 				if re.MatchString(mediaName) {
	// 					codecValInStr := re.ReplaceAllString(mediaName, "$4")
	// 					codecVal, _ := strconv.Atoi(codecValInStr)
	// 					if codecVal < 96 {
	// 						fmt.Println("Doing a lookup")
	// 					}else{
	// 						//Here we are replacing imsInfo_sDPMediaComponent_0_sDPMediaName with imsInfo_sDPMediaComponent_0_sDPMediaDescription_$i
	// 						sDPMediaDescriptionArr := strings.Replace(CdrFieldName, "sDP-Media-Name", "sDP-Media-Descriptions_$i", -1)
	// 						patternForMediaDesc := strings.Replace(patternMediaDesc, "$i", codecValInStr, -1)
	// 						// Regex format sample --> "(a=rtpmap:)" + codecValInStr + "(.*)"
	// 						regEx := regexp.MustCompile(patternForMediaDesc)
	// 						secondCounter := 0
	// 						for {
	// 							sDPMediaDescriptionEle := strings.Replace(sDPMediaDescriptionArr, "$i", strconv.Itoa(secondCounter), -1)
	// 							val, sDPMediaDescriptionEleAvailable := cdr[sDPMediaDescriptionEle]
	// 							if !sDPMediaDescriptionEleAvailable {
	// 								break
	// 								//fmt.Println("No sdpmediadescription present which contains rtpmap:")
	// 								//return errors.New("No sdpmediadescription present which contains rtpmap")
	// 							} else {
	// 								if regEx.MatchString(val.(string)) {
	// 									valueFromsDPMediaDescriptionEle := regEx.ReplaceAllString(val.(string), "$2")
	// 									cdr[destCdrField] = valueFromsDPMediaDescriptionEle
	// 									break
	// 								}
	// 							}
	// 							secondCounter++
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// 	encodedData,_ := json.Marshal(cdr)
	// 	ioutil.WriteFile("isbcCodecOutput",encodedData, 0644)
	}


