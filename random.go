// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"math"
)

func main() {
	/*
	fmt.Println(reflect.TypeOf(6000/100000000))
	f := 1.234567399
	ratio := math.Pow(10, float64(8))
	fmt.Println(math.Ceil(f*ratio)/ratio)
	*/
	numArr1 := []interface{}{12345678916.0, 23456789168.0, 34567891682.0}
	var numArr2[3] interface{} 
	//divisionFactor := 100000000
	/*
	for i := 0; i < len(numArr1); i++ {
		field := numArr1[i].(float64)
		
		field = field/float64(divisionFactor)
		fmt.Println(field)
		ratio := math.Pow(10, float64(2))
		numArr2[i] = math.Round(field*ratio)/ratio
	}
	
	for i := 0; i < len(numArr1); i++ {
		field , ok := numArr1[i].(float64)
		if !ok {
			//ml.MavLog(ml.INFO, transID, action.CURRENCYCONVERSION.SourceFields[i] , " is not present in cdr or dataType is invalid")
			fmt.Println("some error")
			continue               
		}
		field = field/float64(divisionFactor)
		ratio := math.Pow(10, float64(2))
		switch "ROUND" {
		case "ROUND":
				numArr2[i] = math.Round(field*ratio)/ratio
		case "CEIL":
				numArr2[i] = math.Ceil(field*ratio)/ratio
		case "FLOOR":
				numArr2[i] = math.Floor(field*ratio)/ratio
		}
		
	}
	
	
	        for i := 0; i < len(numArr1); i++ {
            field , ok := numArr1[i].(float64)
            if !ok {
                //ml.MavLog(ml.INFO, transID, action.CURRENCYCONVERSION.SourceFields[i] , " is not present in cdr or dataType is invalid")
                continue               
            }
            field = field/float64(100000000)
			if flag == "Prepaid"{
            	ratio = math.Pow(10, float64(2))
			}else if flag == "Postpaid"{
				ratio = math.Pow(10, float64(8))
			}
            switch "CEIL" {
            case "ROUND":
                    numArr2[i] = math.Round(field*ratio)/ratio
            case "CEIL":
					numArr2[i] = math.Ceil(field*ratio)/ratio
            case "FLOOR":
					numArr2[i] = math.Floor(field*ratio)/ratio
            }      
        }

	*/
	flag := "Prepaid"
	var ratio float64
	if flag == "Prepaid"{
		ratio = math.Pow(10, float64(2))
	}else if flag == "Postpaid" {
		ratio = math.Pow(10, float64(8))
	}
	for i := 0; i < len(numArr1); i++ {
		field , ok := numArr1[i].(float64)
		if !ok {
			//ml.MavLog(ml.INFO, transID, action.CURRENCYCONVERSION.SourceFields[i] , " is not present in cdr or dataType is invalid")
			continue               
		}
		field = field/float64(100000000)
		switch "CEIL" {
			case "ROUND":
					numArr2[i] = math.Round(field*ratio)/ratio
			case "CEIL":
					numArr2[i] = math.Ceil(field*ratio)/ratio
			case "FLOOR":
					numArr2[i] = math.Floor(field*ratio)/ratio
		}  
	}  	
	fmt.Println(numArr2)
}
