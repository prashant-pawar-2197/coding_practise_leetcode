package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"

	"strconv"
)


var TemplateMap map[string]TemplateList
var TemplateConfig Template

type Template struct {
	AggregateTemplate struct {
		Key 		string `json:"key"`
		Data		[]struct {
			Key   string                  `json:"key"`
			Value string `json:"value"`
		} `json:"data,omitempty"`
		DecodedData map[string]string
		Templates	[]TemplateList `json:"templates"`
	} `json:"aggregateTemplate"`
}
type TemplateList struct {
	TemplateName string `json:"TemplateName"`
	Template     []struct {
		FieldName string `json:"FieldName"`
		FieldType string `json:"FieldType"`
		Value     string `json:"Value"`
	} `json:"Template"`
	DecodedData  map[string]interface{}
}
func (s *Template) UnmarshalJSON(data []byte) error {
    type Template2 Template
    if err := json.Unmarshal(data, (*Template2)(s)); err != nil {
        return err
    }
	s.AggregateTemplate.DecodedData = make(map[string]string)
	for _, v := range s.AggregateTemplate.Data {
		s.AggregateTemplate.DecodedData[v.Key] = v.Value
	}
	return nil
}

// Custom unmarshal function to decode value for literal Location type
func (s *TemplateList) UnmarshalJSON(data []byte) error {
    type TemplateList2 TemplateList
    if err := json.Unmarshal(data, (*TemplateList2)(s)); err != nil {
        return err
    }

	s.DecodedData = make(map[string]interface{})

    for _, template := range s.Template {
       //ml.MavLog(ml.INFO, "", "Template item", i+1, template.FieldName, template.Value)
        if template.FieldType == "Integer" {
            val, err1 := strconv.Atoi(template.Value)
            if err1 != nil {
                //ml.MavLog(ml.ERROR, "", "error while converting string to int, using default as 0", err1)
                val = 0
            }
			s.DecodedData[template.FieldName] = val
        } else if template.FieldType == "String" {
			s.DecodedData[template.FieldName] = template.Value
        } else if template.FieldType == "Boolean" {
            if template.Value == "true" {
				s.DecodedData[template.FieldName] = true
            } else if template.Value == "false" {
				s.DecodedData[template.FieldName] = false
            } else {
                //ml.MavLog(ml.ERROR, "", "invalid value:", template.Value, "configured for type boolean, using default of false.")
				s.DecodedData[template.FieldName] = false
            }
        } else {
            //ml.MavLog(ml.ERROR, "", "Invalid FieldType:", template.FieldType, " provided for FieldName:",
                      //template.FieldName, " allowed only Integer and String.")
            //ml.MavLog(ml.ERROR, "", "by default setting FieldType to String.")
			s.DecodedData[template.FieldName] = template.Value
        }
    }
	return nil
}

func initTemplate(tmplFile string) error {
	//ml.MavLog(ml.INFO, dummy, "Loading templates from:", tmplFile)
	jsonFile, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		//ml.MavLog(ml.ERROR, dummy, err)
		return err
	}
	err = json.Unmarshal([]byte(jsonFile), &TemplateConfig)
	if err != nil {
		//ml.MavLog(ml.ERROR, dummy, err)
		return err
	}
	ConfigTempleteListCfg()

	return nil
}

func ConfigTempleteListCfg() {
	TemplateMap = make(map[string]TemplateList)
	for _, template := range TemplateConfig.AggregateTemplate.Templates {
		TemplateMap[template.TemplateName] = template
	}
}
func main() {
	initTemplate("C:/Users/pawarpr/OneDrive - Mavenir Systems, Inc/Documents/GoPractise and notes/sampleAggregateTemplate.json")
	fmt.Println(TemplateConfig.AggregateTemplate.DecodedData[TemplateConfig.AggregateTemplate.Key])
	encodedData,_ := json.Marshal(TemplateConfig)
	ioutil.WriteFile("actualAggregateData",encodedData, 0644)
}