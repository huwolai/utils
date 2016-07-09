package db

import (
	"fmt"
	"time"
	"github.com/gocraft/dbr"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/config"
)
var conn *dbr.Connection


//初始化MYSQL
func InitMysql() {

	fmt.Println("init mysql...");
	loc,_ := time.LoadLocation("Local")

	mysql_host :=config.GetValue("mysql_host")
	mysql_db :=config.GetValue("mysql_db")
	mysql_user :=config.GetValue("mysql_user")
	mysql_password :=config.GetValue("mysql_password")
	connInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",mysql_user,mysql_password,mysql_host,mysql_db,loc.String())
	fmt.Println(connInfo);
	var err error;
	conn,err = dbr.Open("mysql",connInfo,nil)
	util.CheckErr(err)
	conn.SetMaxOpenConns(2000)
	conn.SetMaxIdleConns(1000)
	conn.Ping()

	fmt.Println("mysql inital is success");
}

func NewSession() *dbr.Session {
	if conn==nil {
		InitMysql()
	}

	return conn.NewSession(nil)
}
