package security

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/app"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

const (
	//APP认证
	SECURITY_LEVEL_APP  = "app_authed"
	//用户认证
	SECURITY_LEVEL_USER = "user_authed"
	//资源认证
	SECURITY_LEVEL_RESOURCE = "resource_authed"
)

type AppSecurity struct {
	AppId string
	Sign string
}

type UserSecurity struct {
	OpenId string
	Token string
	Rid string
}

type Security struct {
	// 应用安全
	AppSecurity *AppSecurity
	//用户安全
	UserSecurity *UserSecurity
	//安全等级
	Level string
}

//认证
func Auth(req *http.Request) (*Security,error) {

	securityLevel :=GetSecurityLevel(req)
	if securityLevel==""{
		log.Warn("没有认证信息！")
		return nil,errors.New("没有认证信息！")
	}
	

	if securityLevel == SECURITY_LEVEL_APP{//app级别的权限
		log.Info("APP 认证方式..")
		appSign,err :=app.Auth(req)
		if err!=nil{
			return nil,err
		}
		appSecurity :=&AppSecurity{}
		appSecurity.AppId = appSign.App.AppId
		appSecurity.Sign = appSign.Sign

		return &Security{Level:securityLevel,AppSecurity:appSecurity},nil
	}


	if securityLevel == SECURITY_LEVEL_USER {//用户级别的权限
		log.Info("用户Authorization认证方式..")
		authU,err :=AuthUsers(req)
		if err!=nil{
			return nil,err
		}
		userSecurity :=&UserSecurity{}
		userSecurity.OpenId = authU.OpenId
		userSecurity.Rid = authU.Rid
		userSecurity.Token = app.GetParamInRequest("Authorization",req)

		return &Security{Level:securityLevel,UserSecurity:userSecurity},nil
	}
	log.Warn("没有认证方式！")
	return nil,errors.New("没有此认证方式！")
}

//权限认证和open_id认证
func AuthAndOpenId(openId string,req *http.Request) (*Security,error)  {
	sec,err := Auth(req)
	if err!=nil{
		return nil,err
	}
	if !OpenIdIsOk(openId,sec) {//用户不被允许
		return nil,errors.New("用户不被允许操作！")
	}

	return  sec,nil
}

//open_id是否被允许
func OpenIdIsOk(openId string,sec *Security) bool  {

	if sec.Level==SECURITY_LEVEL_USER {
		if openId == sec.UserSecurity.OpenId {

			return true
		}
	}

	if sec.Level==SECURITY_LEVEL_APP {

		return true
	}

	return false
}

func GetSecurityLevel(req *http.Request) string  {
	var securityLevel string
	token :=app.GetParamInRequest("Authorization",req)
	if token!=""{
		securityLevel = SECURITY_LEVEL_USER
	}

	sign :=app.GetParamInRequest("sign",req)
	if sign!=""{
		securityLevel = SECURITY_LEVEL_APP
	}

	return securityLevel
}