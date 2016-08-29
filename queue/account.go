package queue

import (
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"time"
)

const (
	//账户金额改变
	ACCOUNT_AMOUNT_EVENT_CHANGE ="ACCOUNT_AMOUNT_EVENT_CHANGE"
)


type AccountEvent struct  {
	EventKey string
	//事件名
	EventName string
	//事件版本
	Version string
	//事件数据
	Content *AccountEventContent

}

func NewAccountEvent() *AccountEvent  {

	return &AccountEvent{}
}

func NewAccountEventContent() *AccountEventContent  {

	return &AccountEventContent{}
}
//账户事件正文
type AccountEventContent struct  {
	AppId string
	//账户ID
	OpenId string
	//行为
	Action string
	//账户记录的唯一标识
	SubTradeNo string
	//变动金额
	ChangeAmount int64
}


//创建请求生产者
func createAccountQueue() *amqp.Channel {
	name :="account"
	requestChannel = GetChannel()
	//声明一个trade Exchange
	err := requestChannel.ExchangeDeclare(name+"Ex", "direct", true, false, false, false, nil)
	util.CheckErr(err)
	//声明一个声明一个trade Queue
	queue,err := requestChannel.QueueDeclare(name+"Queue",true,false,false,false,nil)
	util.CheckErr(err)
	//将队里绑定到对应的Exchange
	err = requestChannel.QueueBind(queue.Name,name,name+"Ex",false,nil)
	util.CheckErr(err)

	return  requestChannel
}

//发布订单事件
func PublishAccountEvent(event *AccountEvent) error  {
	if requestChannel==nil{
		requestChannel  =createAccountQueue()
	}

	if event.Version=="" {
		event.Version = EVENT_VERSION_V1
	}

	msgbytes,err := json.Marshal(event)
	if err!=nil{
		log.Error(err)
		return err
	}
	name :="account"
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         msgbytes,
	}
	err = requestChannel.Publish(name+"Ex", name, false, false, msg)

	return err
}

//消费订单事件
func ConsumeAccountEvent(consumer string,fn func(accountEvent *AccountEvent, dv amqp.Delivery))  {
	if requestChannel==nil{
		requestChannel  =createAccountQueue()
	}
	name :="account"
	msgs, err := requestChannel.Consume(name+"Queue", consumer, false, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var request *AccountEvent
				util.ReadJsonByByte(d.Body,&request)
				fn(request,d)
			}
		}()
	}else{
		log.Error(err)
	}
}
