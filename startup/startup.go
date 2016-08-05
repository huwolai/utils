package startup

import (
	"io/ioutil"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/file"
	"gitlab.qiyunxin.com/tangtao/utils/util"
        "github.com/rubenv/sql-migrate"
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

	migrations := &migrate.FileMigrationSource{
		Dir: "config/sql",
	}

	_, err := migrate.Exec(db.NewSession().DB, "mysql", migrations, migrate.Up)

	return err
}

