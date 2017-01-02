package remote

import (
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
	"fmt"
	"crypto/md5"
	"time"
)

const (
	//手机注册
	REG_TYPE_MOBILE = 1
	//邮箱注册
	REG_TYPE_EMAIL = 2
)

//添加用户 regtype 1.手机注册 2.邮箱注册
func AddUser(username string,password string,nickname string,regtype int) (map[string]interface{},error)  {

	//签名
	tm := time.Now().Unix()
	signStr := username+"_"+fmt.Sprintf("%d",tm)
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStr))
	sign := md5Ctx.Sum(nil)
	md5Ctx = md5.New()
	md5Ctx.Write(sign)
	sign = md5Ctx.Sum(nil)

	imanagerUrl :=getImanagerPhpUrl()

	queryParam :=map[string]string{
		"username": username,
		"password": password,
		"nickname": nickname,
		"time": fmt.Sprintf("%d",tm),
		"sign": string(sign),
		"regtype": fmt.Sprintf("%d",regtype),
	}
	response,err :=network.Get(imanagerUrl+"/Cust/Init/addUser",queryParam,nil)
	if err!=nil{
		log.Error(err)
		log.Error("请求失败！")
		return nil,err
	}

	if response.StatusCode != http.StatusOK {
		log.Error("状态错误：",response.StatusCode)
		return nil,errors.New("不是有效的HTTP状态！")
	}
	log.Info(response.Body)
	return GetResultMap(response.Body)

}