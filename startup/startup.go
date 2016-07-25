package startup

import (
	"io/ioutil"
	"log"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/file"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

//是否已安装
func IsInstall() bool  {

	filename :="install.txt"
	if file.CheckFileIsExist(filename) {
		content,err :=ioutil.ReadFile(filename)
		util.CheckErr(err)
		if string(content)=="1" {

			return true
		}
	}

	err := ioutil.WriteFile(filename,[]byte("1"),0666)
	util.CheckErr(err)

	return false;
}



//初始化DB数据
func InitDBData() error  {
	content, err := ioutil.ReadFile("config/init.sql")
	if err!=nil{
		log.Println(err)
		return err
	}
	_,er := db.NewSession().Exec(string(content))
	if er!=nil{
		log.Println(er)
		return er
	}
	return err
}

