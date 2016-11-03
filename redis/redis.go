package redis

import (
	"github.com/fzzy/radix/extra/sentinel"
	"github.com/fzzy/radix/redis"
	"log"
	"github.com/fzzy/radix/extra/pool"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"errors"
)

var client *sentinel.Client

var MASTER_NAME string ="mymaster"


type RedisConn struct {
	pl *pool.Pool;
}
var redisConn *RedisConn

func Default() *RedisConn {

	if redisConn==nil{
		redisConn = &RedisConn{}
	}

	return redisConn
}

func Setup() error {

	return Default().Connect(config.GetValue("redis_address").ToString(),10)
}

func (self *RedisConn) Connect(addr string,size int)  error {

	log.Println("init redis..")
	//var err error;
	//client,err = sentinel.NewClient("tcp","127.0.0.1:26378",10,MASTER_NAME)
	//
	//if err!=nil{
	//	log.Println("init redis error",err)
	//	os.Exit(0)
	//}

	var err error;
	self.pl,err = pool.NewPool("tcp",addr,size)

	if err!=nil{
		log.Println("redis is error=",err)
		return err
	}

	log.Println("init redis success")

	return nil

}

//func GetConn()  (*redis.Client){
//
//	if client==nil{
//
//		Init()
//	}
//
//	conn,err  :=client.GetMaster(MASTER_NAME)
//
//	if err!=nil{
//
//		log.Fatal(err);
//		return nil;
//	}
//
//	return conn;
//}

func (self *RedisConn)  getConn()  (*redis.Client){

	if self.pl==nil{

		util.CheckErr(errors.New("请先建立redis连接！"))
	}

	conn,err  :=self.pl.Get()

	if err!=nil{

		log.Fatal(err);
		return nil;
	}

	return conn;
}

func (self *RedisConn) putConn(conn *redis.Client)  {

	//client.PutMaster(MASTER_NAME,conn);

	self.pl.Put(conn)
}

func (self *RedisConn) Set(key string,value interface{})  {

	conn := self.getConn();
	defer self.putConn(conn)

	conn.Cmd("set",key,value)



}

//expire 单位 秒
func (self *RedisConn) SetAndExpire(key string,value interface{},expire float32)  {

	conn := self.getConn();
	defer self.putConn(conn)
	conn.Cmd("set",key,value)

	conn.Cmd("expire",key,expire);


}


func (self *RedisConn) GetString(key string)  (string,error){

	conn := self.getConn();
	defer self.putConn(conn)

	result,err:=conn.Cmd("get",key).Str()

	return result,err

}

// list大小
func (self *RedisConn) Llen(key string) (int64,error) {
	conn := self.getConn();
	defer self.putConn(conn)

	result,err:=conn.Cmd("LLEN",key).Int64()

	return result,err
}

func (self *RedisConn) Lrange(key string,start,stop int64) ([]string,error) {
	conn := self.getConn();
	defer self.putConn(conn)

	result,err:=conn.Cmd("LRANGE",key,start,stop).List()

	return result,err
}

//LREM key count value
//根据参数 count 的值，移除列表中与参数 value 相等的元素。
/**
count 的值可以是以下几种：
count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
count = 0 : 移除表中所有与 value 相等的值。
返回值：
	被移除元素的数量。
	因为不存在的 key 被视作空表(empty list)，所以当 key 不存在时， LREM 命令总是返回 0 。
 */
func (self *RedisConn) Lrem(key string,count int64,value string) (int64,error) {
	conn := self.getConn();
	defer self.putConn(conn)

	result,err:=conn.Cmd("LREM",key,count,value).Int64()

	return result,err
}

/**
LTRIM key start stop
对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
举个例子，执行命令 LTRIM list 0 2 ，表示只保留列表 list 的前三个元素，其余元素全部删除。
下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
当 key 不是列表类型时，返回一个错误。
 */
func (self *RedisConn) Ltrim(key string,start,stop int64) (string,error) {
	conn := self.getConn();
	defer self.putConn(conn)

	result,err:=conn.Cmd("LTRIM",key,start,stop).Str()

	return result,err
}