package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"errors"
	"net/http"
	"net"
	"time"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

var environments = map[string]string{
	"production":    "config/prod.json",
	"preproduction": "config/pre.json",
	"tests":         "config/tests.json",
}

type ConfigValue struct  {

	Value interface{}
}

func (self*ConfigValue) ToString() string  {

	switch v := self.Value.(type){
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	}

	return fmt.Sprintf("%s",self.Value)
}

func (self*ConfigValue) ToInt() int {
	switch v := self.Value.(type){
	case int:
		return v
	case string:
		k,_ := strconv.Atoi(v)
		return k
	case float32:

		return int(v)
	case int64:
		return int(v)
	default:
		fmt.Println(v)
		util.CheckErr(errors.New("不能转换为int类型111"))

	}

	return 0
}

func (self*ConfigValue) ToFloat() float32 {
	switch v := self.Value.(type){
	case float32:
		return v
	case int:

		return float32(v)
	case string:
		f,_ := strconv.ParseFloat(v,20)

		return float32(f)
	}

	util.CheckErr(errors.New("不能转换为float类型"))

	return 0
}

var c *http.Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*10)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 20,
	},
}


var settings map[string]interface{}
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	fmt.Println("环境["+env+"]")
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}

	var configMap map[string]interface{}
	err := LoadSettingsByLocalEnv(env,&configMap)
	util.CheckErr(err)

	var remoteConfigMap map[string]interface{}
	err = LoadSettingByConfigCenter(env,&remoteConfigMap)
	util.CheckErr(err)

	for k,v := range remoteConfigMap  {
		configMap[k] = v
	}

	settings = configMap

}

//通过本地环境加载配置
func LoadSettingsByLocalEnv(env string,resultMap *map[string]interface{}) (error) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)

		util.CheckErr(err)
	}
	jsonErr := json.Unmarshal(content,resultMap)

	return jsonErr
}

//从配置中心加载配置
func LoadSettingByConfigCenter(env string,resultMap *map[string]interface{}) (error)  {

	url,err :=GetConfigApiUrl()
	if err!=nil{
		return err
	}

	response,err := c.Get(url+"/config")
	if err!=nil {
		return err
	}

	defer response.Body.Close()

	err = util.ReadJson(response.Body,resultMap)

	return err

}

func GetConfigApiUrl() (string,error) {

	appId := os.Getenv("APPID")
	if appId=="" {
		return "",errors.New("请在环境变量里配置APPID!")
	}
	env := os.Getenv("GO_ENV")
	if env=="" {
		log.Println("warn:没有配置环境变量GO_ENV 将默认使用预生产环境(preproduction)")
		env = "preproduction"
	}
	configUrl := os.Getenv("CONFIG_URL")
	if configUrl=="" {
		return "",errors.New("请在环境变量里配置CONFIG_URL!")
	}

	return configUrl+"/" +appId + "/" +env,nil
}

func GetValue(key string) *ConfigValue {

	if settings==nil {
		Init()
	}

	value :=&ConfigValue{settings[key]}

	return value
}

func IsTestEnvironment() bool {
	return env == "tests"
}
