package remote

import (
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"net/http"
)

//创建经费群
func CreateFundGroup(_custId,_token string,groupName string,custIds string,amount float32) (map[string]interface{},error)  {

	imanagerUrl :=getImanagerPhpUrl()
	queryParam :=map[string]string{
		"_custid": _custId,
		"_token": _token,
		"custids": custIds,
		"chatname": groupName,
		"amount": fmt.Sprintf("%g",amount),
	}
	response,err :=network.Get(imanagerUrl+"/Chat/SpecialChat/createFundsChat",queryParam,nil)
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

//创建普通群组
func CreateGroup(_custId,_token string,groupName string,custIds string) (map[string]interface{},error)  {

	imanagerUrl :=getImanagerPhpUrl()
	queryParam :=map[string]string{
		"_custid": _custId,
		"_token": _token,
		"custids": custIds,
		"chatname": groupName,
	}
	response,err :=network.Get(imanagerUrl+"/Chat/Chat/createChat",queryParam,nil)
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


//获取群详情
func GetGroupDetail(_custId,_token string,groupNo string) (map[string]interface{},error) {
	imanagerUrl :=getImanagerPhpUrl()
	queryParam :=map[string]string{
		"_custid": _custId,
		"_token": _token,
		"chatid": groupNo,
	}
	response,err :=network.Get(imanagerUrl+"/Chat/Chat/getChatDetail",queryParam,nil)
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

//custids成员custids 逗号隔开
func AddGroupMember(_custId,_token string,groupNo string,custids string) error  {
	imanagerUrl :=getImanagerPhpUrl()
	queryParam :=map[string]string{
		"_custid": _custId,
		"_token": _token,
		"chatid": groupNo,
		"custids": custids,
	}
	response,err :=network.Get(imanagerUrl+"/Chat/Chat/addChatCusts",queryParam,nil)
	if err!=nil{
		log.Error(err)
		log.Error("请求失败！")
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Error("状态错误：",response.StatusCode)
		return errors.New("不是有效的HTTP状态！")
	}

	log.Info(response.Body)

	_,err = GetResultMaps(response.Body)

	return err
}

func getImanagerPhpUrl() string {

	return config.GetValue("imanager_api_url").ToString()
}

