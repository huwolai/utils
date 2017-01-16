package remote

import (
	"encoding/json"
	"bytes"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

func GetResultMaps(body string) ([]map[string]interface{},error) {

	var resultMap map[string]interface{}
	mdz:=json.NewDecoder(bytes.NewBuffer([]byte(body)))
	err := mdz.Decode(&resultMap)
	if err!=nil{
		return nil,err
	}
	status :=resultMap["status"].(float64)
	if status==0 {
		dataintes := resultMap["data"].([]interface{})
		dataMaps := []map[string]interface{}{}
		for _,inte :=range dataintes {
			dataMaps = append(dataMaps,inte.(map[string]interface{}))
		}
		return dataMaps,nil
	}else{
		log.Error(resultMap["msg"])
		return nil,errors.New(resultMap["msg"].(string))
	}
}

func GetResultMap(body string) (map[string]interface{},error) {

	var resultMap map[string]interface{}
	mdz:=json.NewDecoder(bytes.NewBuffer([]byte(body)))
	err := mdz.Decode(&resultMap)
	if err!=nil{
		return nil,err
	}
	status :=resultMap["status"].(float64)
	if status==0 {
		dataMap := resultMap["data"].(map[string]interface{})
		return dataMap,nil
	}else{
		log.Error(resultMap["msg"])
		return nil,errors.New(resultMap["msg"].(string))
	}
}