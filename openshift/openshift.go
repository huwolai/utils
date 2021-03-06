package openshift

import (
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

const (
	DEFUALT_NAMESPACE = "qiyunxin"
	PROTOCOL  =  "http"
	FIX_DOMAIN  = "svc.cluster.local"
	//服务
	SERVICE_PORT = "8080"
	//应用管理
	APPMANAGER_PORT = "8081"
	//权限管理
	SECURITYMANAGER_PORT = "8082"
)

type UserResource struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//角色
	OpenId string `json:"open_id"`
	//资源ID
	ResourceId string `json:"resource_id"`
	Action  string `json:"action"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}


//获取用户资源(权限资源)
func GetUserSources(serviceId,appId,openId string) ([]*UserResource,error)  {
	serviceUrl :=GetServiceSecurityUrl(serviceId)
	resp,err :=network.Get(serviceUrl+ "/v1/_useresources/"+openId+"/apps/"+appId,nil,nil)
	if err!=nil{
		return nil,err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil,errors.New("服务没提供资源服务！")
	} else if (resp.StatusCode == http.StatusOK) {
		var results []*UserResource
		err :=util.ReadJsonByByte([]byte(resp.Body),&results)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
		return results,err
	}else{
		return nil,errors.New("服务请求失败！")
	}
}
//获取权限管理服务地址
func GetServiceSecurityUrl(serviceId string) string  {

	return GetServiceSecurityUrlWithNamespace(DEFUALT_NAMESPACE,serviceId)

}

//获取APP服务地址
func GetServiceAppUrl(serviceId string)  string {

	return GetServiceAppUrlWithNamespace(DEFUALT_NAMESPACE,serviceId)
}
//获取API服务地址
func GetServiceApiUrl(serviceId string) string {

	return GetServiceApiUrlWithNamespace(DEFUALT_NAMESPACE,serviceId)
}

//获取API服务地址
func GetServiceApiUrlWithNamespace(namespace,serviceId string) string {

	return GetServiceUrl(namespace,serviceId,SERVICE_PORT)
}

//获取权限管理服务地址
func GetServiceSecurityUrlWithNamespace(namespace,serviceId string) string  {

	return GetServiceUrl(namespace,serviceId,SECURITYMANAGER_PORT)
}

//通过空间获取APP服务地址
func GetServiceAppUrlWithNamespace(namespace,serviceId string) string  {

	return GetServiceUrl(namespace,serviceId,APPMANAGER_PORT)
}
//获取服务的URL
func GetServiceUrl(namespace,serviceId ,port string) string {

	return PROTOCOL +"://" + serviceId+"." +namespace +"." + FIX_DOMAIN + ":" + port
}
