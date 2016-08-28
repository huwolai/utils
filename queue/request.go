package queue

import (
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"log"
	"encoding/json"
	"time"
)

type RequestModel struct  {

	NotifyUrl string `json:"notify_url"`
	Data interface{} `json:"data"`

}

func NewRequestModel() *RequestModel  {

	return &RequestModel{}
}

var requestChannel *amqp.Channel

//创建请求生产者
func createRequestExchange() *amqp.Channel {


	requestChannel = GetChannel()
	//声明一个trade Exchange
	err := requestChannel.ExchangeDeclare("requestDEx", "x-delayed-message", true, false, false, false, map[string]interface{}{
		"x-delayed-type":"direct",
	})
	util.CheckErr(err)
	//声明一个声明一个trade Queue
	queue,err := requestChannel.QueueDeclare("requestDQueue",true,false,false,false,nil)
	util.CheckErr(err)
	//将队里绑定到对应的Exchange
	err = requestChannel.QueueBind(queue.Name,"requestD","requestDEx",false,nil)
	util.CheckErr(err)


	return  requestChannel
}

//消费请求消息
func ConsumeRequestMsg(fn func(requestModel *RequestModel, dv amqp.Delivery)) {
	if requestChannel==nil{
		requestChannel  =createRequestExchange()
	}
	msgs, err := requestChannel.Consume("requestDQueue", "", true, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var request *RequestModel
				util.ReadJsonByByte(d.Body,&request)
				fn(request,d)
			}
		}()
	}else{
		log.Println("the Consume is error!",err)
	}

}

//发布交易消息
func PublishRequestMsg(requestModel *RequestModel) error  {

	return PublishRequestMsgOfDelay(requestModel,0)

}

//delaySec 延迟发送时间
func PublishRequestMsgOfDelay(request *RequestModel,delaySec int) error {
	if requestChannel==nil{
		requestChannel  =createRequestExchange()
	}

	msgbytes,err := json.Marshal(request)
	if err!=nil{
		log.Println("TradeMsg convert to json is  Fail!")
		return err
	}
	delay :=int64(delaySec*1000)
	msg := amqp.Publishing{
		Headers:map[string]interface{}{
			"x-delay":delay,
		},
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         msgbytes,
	}
	err = requestChannel.Publish("requestDEx", "requestD", false, false, msg)

	return err
}

