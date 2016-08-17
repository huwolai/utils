package security

import (
	"errors"
	"net/http"
)

//认证校验
func CheckAppAuth(req *http.Request) (string,error)  {

	appId := GetParamInRequest("app_id",req)

	if appId==""{

		return appId,errors.New("app_id不能为空")
	}

	return appId,nil
}

//用户认证
func CheckUserAuth(req *http.Request) (string,error)  {
	openId := GetParamInRequest("open_id",req)

	if openId==""{

		return openId,errors.New("open_id不能为空")
	}

	return openId,nil
}

func GetOpenId(req *http.Request) string {
	openId := GetParamInRequest("open_id",req)

	return openId
}

//获取APPID
func GetAppId(req *http.Request) (string,error)  {

	appId :=GetParamInRequest("app_id",req)
	if appId=="" {

		return "",errors.New("app_id不能为空")
	}
	return appId,nil
}

func GetAppId2(req *http.Request) string {
	appId :=GetParamInRequest("app_id",req)

	return appId
}

//在请求中获取AppId
func GetParamInRequest(key string,req *http.Request) string  {

	var value string
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		value = values[0]
	}
	if value=="" {
		value = req.Header.Get(key)
	}

	return value

}
