package qim

import (
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"net/http"
	"errors"
	"fmt"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"strings"
)

type BaseMsgContent struct {
	//消息唯一编号
	Msgno string `json:"msgno"`
	//发送者昵称
	FromCustName string `json:"from_cust_name,omitempty"`
	//发送者ID
	FromCustId string `json:"from_cust_id,omitempty"`
	//正文类型
	Type int `json:"type,omitempty"`
	//当前会话ID
	SessionId string `json:"session_id,omitempty"`
	//当前会话类型
	SessionType int `json:"session_type,omitempty"`
	//业务类型
	OperationType int `json:"operationType"`
	//正文
	Content string `json:"content,omitempty"`
}

func SendMsgForUser(toCustId string, content interface{}) error {

	return SendMsgForUsers([]string{toCustId},content)
}

//发送点对点的消息
func SendMsgForUsers(toCustIds []string, content interface{}) error {

	log.Info("发送IM消息给:",strings.Join(toCustIds,","))
	if content==nil{
		log.Error("消息正文不能为空！")
		return  errors.New("消息正文不能为空！")
	}

	imapiUrl := config.GetValue("im_api_url").ToString()

	//生成消息唯一标识
	msgNo := CreateMsgNo("0",strings.Join(toCustIds,","),"0")
	//下面所有参数注解都是本人猜测 不一定正确
	queryParam := map[string]string{
		"type": "2", //2. 发消息
		"client": "0", // 0.不是客户端
		"msg_no": msgNo, //消息唯一编号
		"sys_type": "3", //0、所有在线人，1、好友friend（custs里面的friends）、2、聊天群（chatids里面的群组）3、指定人、4、kick人
		"cust_ids": "["+strings.Join(toCustIds,",")+"]",
		"chat_ids": "",
		"content": util.ToJson2(content),
	}
	//请求IM接口
	response, err := network.Get(imapiUrl, queryParam, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	if response.StatusCode != http.StatusOK {
		log.Error("返回状态错误:", response.StatusCode)
		return errors.New("返回状态错误")
	}
	log.Info("发送IM消息发送成功:",strings.Join(toCustIds,","))
	log.Info("response=", response.Body)
	return nil
}

func CreateToCustIdMsgNo(toCustId string) string  {

	return CreateMsgNo("0",toCustId,"0")
}

func CreateMsgNo(fromCustId string,toCustId string,chatId string) string {
	msgNoStr := fmt.Sprintf("%s:%s:%s:%s",fromCustId,toCustId,chatId,time.Now().UnixNano())

	return GetSign64(msgNoStr)
}