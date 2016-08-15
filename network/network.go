package network

import (
	"github.com/sendgrid/rest"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)



func Post(url string, body []byte,headers map[string]string) (*rest.Response,err error)  {

	log.Debug("请求地址:",url)
	request :=rest.Request{
		Method:rest.Post,
		BaseURL:url,
		Body:body,
		Headers:headers,
	}
	response, err := rest.API(request)
	if err != nil {
		log.Error("请求失败:",err)
		return nil,err
	}

	log.Debug("返回结果:",response.Body)

	return response,nil
}


func GetJson(url string,queryParams map[string]string,headers map[string]string) (byts []byte,err error) {

	request :=rest.Request{
		Method:rest.Get,
		BaseURL:url,
		Headers:headers,
		QueryParams:queryParams,
	}
	response, err := rest.API(request)
	if err != nil {

		return nil,err
	}

	return []byte(response.Body),nil
}