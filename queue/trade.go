package queue

import (
	"github.com/streadway/amqp"
	"time"
	"encoding/json"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)
type TradeMsg struct {
	//签名
	Sign string `json:"sign"`
	//交易号
	TradeNo string `json:"trade_no"`
	//第三方系统中的交易号
	OutTradeNo string `json:"out_trade_no"`
	//第三方交易类型
	OutTradeType string `json:"out_trade_type"`
	// 备用数据
	Memo string  `json:"memo"`
	//预付款代号
	ImprestCode string `json:"imprest_code"`
	//应用ID
	AppId string  `json:"app_id"`
	//用户openID
	OpenId string `json:"open_id"`
	//交易时间
	TradeTime int64 `json:"trade_time"`
	//交易金额
	Amount int64  `json:"amount"`
	//交易标题
	Title string `json:"title"`
	//交易备注
	Remark string `json:"remark"`
	//交易通知地址
	NotifyUrl string `json:"notify_url"`

}

func NewTradeMsg() *TradeMsg  {

	return &TradeMsg{}
}

//创建交易生产者
func createTradeExchange() *amqp.Channel {


	tradeChannel := GetChannel()
	//声明一个trade Exchange
	err := tradeChannel.ExchangeDeclare("tradeDEx", "x-delayed-message", true, false, false, false, map[string]interface{}{
		"x-delayed-type":"direct",
	})
	util.CheckErr(err)
	//声明一个声明一个trade Queue
	queue,err := tradeChannel.QueueDeclare("tradeDQueue",true,false,false,false,nil)
	util.CheckErr(err)
	//将队里绑定到对应的Exchange
	err = tradeChannel.QueueBind(queue.Name,"tradeD","tradeDEx",false,nil)
	util.CheckErr(err)


	return  tradeChannel
}

func PublishTradeMsgOfDelay(tradeMsg *TradeMsg,delaySec int) error {
	tradeChannel  :=createTradeExchange()

	msgbytes,err := json.Marshal(tradeMsg)
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
	err = tradeChannel.Publish("tradeDEx", "tradeD", false, false, msg)

	return err
}

//发布交易消息
func PublishTradeMsg(tradeMsg *TradeMsg) error  {

	return PublishTradeMsgOfDelay(tradeMsg,0)

}

//消费交易消息
func ConsumeTradeMsg(fn func(tradeMsg *TradeMsg, dv amqp.Delivery)) {
	tradeChannel  :=createTradeExchange()
	msgs, err := tradeChannel.Consume("tradeDQueue", "", false, false, false, false, nil)

	if err==nil{
		go func() {

			for d := range msgs {
				var tradMsg *TradeMsg
				err := util.ReadJsonByByte(d.Body,&tradMsg)
				if err!=nil{
					log.Info("交易数据格式有误！")
					log.Error(err)
					continue
				}
				fn(tradMsg,d)
			}
		}()
	}else{
		log.Error("the Consume is error!",err)
	}

}