package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"log"
	"errors"
	"net/http"
	"net"
	"time"
)

var environments = map[string]string{
	"production":    "config/prod.json",
	"preproduction": "config/pre.json",
	"tests":         "config/tests.json",
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


var settings map[string]string
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}

	var configMap map[string]string
	err := LoadSettingsByLocalEnv(env,&configMap)
	util.CheckErr(err)

	var remoteConfigMap map[string]string
	err = LoadSettingByConfigCenter(env,&remoteConfigMap)
	util.CheckErr(err)

	for k,v := range remoteConfigMap  {
		configMap[k] = v
	}

	settings = configMap

}

//通过本地环境加载配置
func LoadSettingsByLocalEnv(env string,resultMap map[string]string) (error) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)

		util.CheckErr(err)
	}
	jsonErr := json.Unmarshal(content,resultMap)

	return jsonErr
}

//从配置中心加载配置
func LoadSettingByConfigCenter(env string,resultMap map[string]string) (error)  {

	url,err :=GetConfigApiUrl()
	if err!=nil{
		return nil,err
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
		return "",errors.New("请再环境变量里配置APPID!")
	}
	env := os.Getenv("ENV")
	if env=="" {
		log.Println("warn:没有配置环境变量GO_ENV 将默认使用预生产环境(preproduction)")
		env = "preproduction"
	}
	configUrl := os.Getenv("CONFIG_URL")
	if configUrl=="" {
		return "",errors.New("请再环境变量里配置APPID!")
	}

	return configUrl+"/" +appId + "/" +env,nil
}

func GetValue(key string) string {

	if settings==nil {

	}

	return settings[key]
}

func IsTestEnvironment() bool {
	return env == "tests"
}
