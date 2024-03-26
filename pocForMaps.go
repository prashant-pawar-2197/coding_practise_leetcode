package main

import (
	"fmt"
)

func main()  {
	m := map[string]interface{}{"Prashant": map[string]interface{}{"state" : map[string]interface{}{"hm":"Pune"}}}
	if _, ok := m["Prashant"]; ok{
		if _, ok := m["Prashant"].(map[string]interface{})["state"]; ok{
			if _, ok := m["Prashant"].(map[string]interface{})["state"].(map[string]interface{})["city"]; ok{
				fmt.Println("Found")
			}else{
				fmt.Println("city absent")
			}
		}else{
			fmt.Println("state absent")
		}
	}
}