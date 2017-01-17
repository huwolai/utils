package queue

import (
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type RequestModel struct  {

	NotifyUrl string `json:"notify_url"`
	Data interface{} `json:"data"`

}

func NewRequestModel() *RequestModel  {

	return &RequestModel{}
}


//创建请求生产者
func createRequestExchange(queueName string) *amqp.Channel {


	requestChannel := GetChannel()
	//声明一个trade Exchange
	err := requestChannel.ExchangeDeclare("requestEx", "x-delayed-message", true, false, false, false, map[string]interface{}{
		"x-delayed-type":"topic",
	})
	util.CheckErr(err)
	//声明一个声明一个trade Queue
	queue,err := requestChannel.QueueDeclare(queueName,true,false,false,false,nil)
	util.CheckErr(err)
	//将队里绑定到对应的Exchange
	err = requestChannel.QueueBind(queue.Name,"http.request","requestEx",false,nil)
	util.CheckErr(err)


	return  requestChannel
}

//消费请求消息
func ConsumeRequestMsg(queueName string,fn func(requestModel *RequestModel, dv amqp.Delivery)) {
	requestChannel  :=createRequestExchange(queueName)
	msgs, err := requestChannel.Consume(queueName, "", false, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var request *RequestModel
				err = util.ReadJsonByByte(d.Body,&request)
				if err!=nil{
					log.Error(err)
					return
				}
				fn(request,d)
			}
		}()
	}else{
		log.Error("the Consume is error!",err)
	}

}

//发布交易消息
func PublishRequestMsg(requestModel *RequestModel) error  {

	return PublishRequestMsgOfDelay(requestModel,0)

}

//delaySec 延迟发送时间
func PublishRequestMsgOfDelay(request *RequestModel,delaySec int) error {
	requestChannel := GetChannel()
	//声明一个 Exchange
	err := requestChannel.ExchangeDeclare("requestEx", "x-delayed-message", true, false, false, false, map[string]interface{}{
		"x-delayed-type":"topic",
	})
	util.CheckErr(err)

	msgbytes,err := json.Marshal(request)
	if err!=nil{
		log.Error("TradeMsg convert to json is  Fail!")
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
	err = requestChannel.Publish("requestEx", "http.request", false, false, msg)

	return err
}

