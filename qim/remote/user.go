package remote

import (
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
	"time"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

func UserLogin(username,password string) (map[string]interface{},error)  {
	imanagerUrl :=getImanagerPhpUrl()

	timestamp := time.Now().Unix()

	signStr := util.MD5(util.MD5(fmt.Sprintf("%s_%d",username,timestamp)))

	param :=map[string]string{
		"username": username,
		"password": password,
		"sign": signStr,
		"time": fmt.Sprintf("%d",timestamp),
	}
	response,err :=network.Post(imanagerUrl+"/Cust/Init/login",param,nil)
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


func AddUser(username string,mobile string,email string,password string,nickname string) (map[string]interface{},error)  {

	////签名
	//tm := time.Now().Unix()
	//signStr := username+"_"+fmt.Sprintf("%d",tm)
	//md5Ctx := md5.New()
	//md5Ctx.Write([]byte(signStr))
	//sign := md5Ctx.Sum(nil)
	//signS := hex.EncodeToString(sign)
	//
	//md5Ctx = md5.New()
	//md5Ctx.Write([]byte(signS))
	//sign = md5Ctx.Sum(nil)
	//signS = hex.EncodeToString(sign)

	imanagerUrl :=getImanagerPhpUrl()

	queryParam :=map[string]string{
		"username": username,
		"password": password,
		"nickname": nickname,
		"email": email,
		"mobile": mobile,
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