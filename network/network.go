package network

import (
	"fmt"
	"net/http"
	"net"
	"time"
	"io"
	"bytes"
	"io/ioutil"
)

var c *http.Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*10)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 20,
	},
}

func Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)  {

	return c.Post(url,bodyType,body)
}

func PostJson(url string,jsonData []byte) (byts []byte,err error)  {

	resp,err := Post(url,"application/json;utf-8",bytes.NewReader(jsonData))
	if err!=nil{

		return nil,err
	}

	bys,er := ioutil.ReadAll(resp.Body)

	return bys,er
}

func GetJson(url string)  {


}