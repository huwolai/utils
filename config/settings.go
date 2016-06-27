package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

var environments = map[string]string{
	"production":    "config/prod.json",
	"preproduction": "config/pre.json",
	"tests":         "config/tests.json",
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
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)

		util.CheckErr(err)
	}
	jsonErr := json.Unmarshal(content, &settings)
	util.CheckErr(jsonErr)
}


func GetValue(key string) string {


	return settings[key]
}

func IsTestEnvironment() bool {
	return env == "tests"
}
