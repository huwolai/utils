package network

import (
	"io/ioutil"
	"github.com/sendgrid/rest"
)



func Post(url string, body []byte,headers map[string]string) (byts []byte,err error)  {

	request :=rest.Request{
		Method:rest.Post,
		BaseURL:url,
		Body:body,
		Headers:headers,
	}
	response, err := rest.API(request)
	if err != nil {

		return nil,err
	}

	bys,er := ioutil.ReadAll(response.Body)

	return bys,er
}


func GetJson(url string,queryParams map[string]string,headers map[string]string)  {

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

	bys,er := ioutil.ReadAll(response.Body)

	return bys,er
}