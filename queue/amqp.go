package queue

import (
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"time"
)

var connection *amqp.Connection


func SetupAMQP(url string,restartCallbackFuc func ())  {

	if url=="" {
		url ="amqp://guest:guest@localhost:5672/"
	}
	go restConn(url,restartCallbackFuc)
	//var err error
	//connection, err = amqp.Dial(url)
	//util.CheckErr(err)
	//errChanel := make( chan *amqp.Error)
	//connection.NotifyClose(errChanel)
	//select {
	//case   amqerr :=<-errChanel:
	//	log.Error(amqerr)
	//connection, err = amqp.Dial(url)
	//if err!=nil{
	//	log.Error(err)
	//	time.Sleep(2*time.Second)
	//	continue
	//}
	//log.Info("connection=",connection)
	//connection.NotifyClose(errChanel)
	////log.Error("rabbitmq is close!",amqerr.Error())
	//restartCallbackFuc("")

	//}
}

func restConn(url string,restartCallbackFuc func ())  {
	var err error
	connection, err = amqp.Dial(url)
	if err!=nil{
		log.Error(err)

	}else{
		log.Info("amq连接成功！")
		restartCallbackFuc()
		errChanel := make( chan *amqp.Error)
		connection.NotifyClose(errChanel)
		select {
		case  <-errChanel:
			restConn(url,restartCallbackFuc)
		}
		return
	}
	//util.CheckErr(err)

	if err!=nil{
		for {
			time.Sleep(2*time.Second)
			restConn(url,restartCallbackFuc)

		}
	}


}

func GetChannel() *amqp.Channel {
	channel, err := connection.Channel()
	util.CheckErr(err)
	return channel
}