package mongodb

import (
	"gopkg.in/mgo.v2"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)
var (
	db *mgo.Database
)
func Setup()  {
	session, err := mgo.Dial(config.GetValue("mongodb_host").ToString())
	util.CheckErr(err)

	db =session.DB(config.GetValue("mongodb_db").ToString())

	err = db.Login(config.GetValue("mongodb_user").ToString(),config.GetValue("mongodb_password").ToString())
	util.CheckErr(err)
}

func GetDB() *mgo.Database {

	return db
}

func GetCollection(name string) *mgo.Collection  {
	return GetDB().C(name)
}