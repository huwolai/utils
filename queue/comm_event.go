package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

const EVENT_VERSION_V1  = "v1"
var commEventChannel *amqp.Channel
type CommEvent struct  {
	//事件头
	Header CommEventHeader
	//事件正文
	Content interface{}

}

//事件头
type CommEventHeader struct {
	//事件KEY
	EventKey string
	//事件名
	EventName string
	//事件版本
	Version string
}

func NewCommEvent() CommEvent  {

	return CommEvent{}
}

//创建请求生产者
func createCommEventQueue(queueName string) *amqp.Channel {
	name :="commevent"
	commEventChannel = GetChannel()
	//声明一个 Exchange
	err := commEventChannel.ExchangeDeclare(name+"Ex", "topic", true, false, false, false, nil)
	util.CheckErr(err)
	//声明一个声明一个 Queue
	queue,err := commEventChannel.QueueDeclare(queueName,true,false,false,false,nil)
	util.CheckErr(err)
	//将队里绑定到对应的Exchange
	err = commEventChannel.QueueBind(queue.Name,name,name+"Ex",false,nil)
	util.CheckErr(err)

	return  commEventChannel
}

//发布订单事件
func PublishCommEvent(event CommEvent) error  {

	name :="commevent"
	if commEventChannel==nil{
		commEventChannel = GetChannel()
		//声明一个 Exchange
		err := commEventChannel.ExchangeDeclare(name+"Ex", "topic", true, false, false, false, nil)
		util.CheckErr(err)
	}

	if event.Header.Version=="" {
		event.Header.Version = EVENT_VERSION_V1
	}

	msgbytes,err := json.Marshal(event)
	if err!=nil{
		log.Error(err)
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         msgbytes,
	}
	err = commEventChannel.Publish(name+"Ex", name, false, false, msg)

	return err
}

//消费订单事件
func ConsumeCommEvent(queueName string,fn func(event *CommEvent, dv amqp.Delivery))  {
	if commEventChannel==nil{
		commEventChannel  =createCommEventQueue(queueName)
	}
	msgs, err := commEventChannel.Consume(queueName, "", false, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var request *CommEvent
				err = util.ReadJsonByByte(d.Body,&request)
				if err!=nil{
					log.Error(err)
					return
				}
				fn(request,d)
			}
		}()
	}else{
		log.Error(err)
	}
}