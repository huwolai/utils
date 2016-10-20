package app

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type App struct {
	AppId string
	AppKey string
	AppName string
	AppDesc string
}

func init()  {
	http.HandleFunc("/v1/apps",Apps)

}

func Apps(w http.ResponseWriter, r *http.Request)  {

	log.Info("测试")

	w.Write([]byte("999999999999"))
}