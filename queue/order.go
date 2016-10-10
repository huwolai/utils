package queue

import (
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"time"
)

const (
	//订单创建
	ORDER_EVENT_CREATED ="ORDER_EVENT_CREATED"

	//订单正在支付
	ORDER_EVENT_PAYING = "ORDER_EVENT_PAYING"

	//订单已付款
	ORDER_EVENT_PAID = "ORDER_EVENT_PAID"
)

type OrderEvent struct  {
	//事件KEY
	EventKey string
	//事件名
	EventName string
	//事件版本
	Version string
	//事件数据
	//事件数据
	Content *OrderEventContent

}

type OrderEventContent struct {
	//订单号
	OrderNo string
	//订单类型
	OrderType int
	//创建时间
	CreateTime string
	//订单标题
	Title string
	//订单金额
	Amount float64
	//下单用户
	OpenId string
	Json string
	Flag string
	//订单项
	Items []*OrderEventItem
	//扩展数据 (mobile)
	ExtData map[string]interface{}
}

type OrderEventItem struct {
	OrderNo string
	Num int
	Title string
	Price float64
	TotalPrice float64
	Flag string
	Json string
}

func NewOrderEventContent() *OrderEventContent   {

	return &OrderEventContent{}
}

func NewOrderEvent() *OrderEvent  {

	return &OrderEvent{}
}

func NewOrderEventItem() *OrderEventItem {

	return &OrderEventItem{}
}

//创建请求生产者
func createOrderQueue() *amqp.Channel {
	name :="order"
	requestChannel = GetChannel()
	//声明一个trade Exchange
	err := requestChannel.ExchangeDeclare(name+"Ex", "topic", true, false, false, false, nil)
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
func PublishOrderEvent(event *OrderEvent) error  {
	if requestChannel==nil{
		requestChannel  =createOrderQueue()
	}

	if event.Version=="" {
		event.Version = EVENT_VERSION_V1
	}

	msgbytes,err := json.Marshal(event)
	if err!=nil{
		log.Error(err)
		return err
	}
	name :="order"
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
func ConsumeOrderEvent(fn func(event *OrderEvent, dv amqp.Delivery))  {
	if requestChannel==nil{
		requestChannel  =createOrderQueue()
	}
	name :="order"
	msgs, err := requestChannel.Consume(name+"Queue", "", false, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var request *OrderEvent
				util.ReadJsonByByte(d.Body,&request)
				fn(request,d)
			}
		}()
	}else{
		log.Error(err)
	}
}
