package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"github.com/rubenv/sql-migrate"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"net/http"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type Source struct {
	//资源ID
	Id string `json:"id"`
	//资源名称
	Name string `json:"name"`
	Description string `json:"description"`
	//资源
	Resource string `json:"resource"`
	//
	Permissions string `json:"permissions"`
}

type UserSource struct {
	//资源ID
	Id int64 `json:"id"`
	//应用ID
	AppId string `json:"app_id"`
	//用户ID
	OpenId string `json:"open_id"`
	//资源ID
	SourceId string `json:"source_id"`
	//资源行为
	Action string `json:"action"`
}

type UserSourceWrap struct {
	OpenId string `json:"open_id"`
	AppId string `json:"app_id"`
	Sources []*UserSource `json:"sources"`
}


var srsAll []Source

func InitSources(sources []Source)  {

	srsAll = sources
}

// 安装
func Setup()  {
	if srsAll ==nil{
		panic(errors.New("请先调用InitSources初始化资源!"))
		return
	}
	go func() {
		err := InitDB()
		if err!=nil {
			log.Error(err)
			log.Info("初始化安全管理的DB失败！")
			return
		}
		router :=gin.Default()
		router.POST("/v1/_usersources",UserSourcesAdd)
		router.GET("/v1/_sources",SourcesAll)

		log.Info("init security manager on 8082!")

		router.Run(":8082")
	}()
}

//添加用户资源
func UserSourcesAdd(c *gin.Context)  {
	var uswrap *UserSourceWrap
	err :=c.BindJSON(&uswrap)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误！")
		return
	}

	if uswrap==nil||uswrap.Sources==nil{
		util.ResponseError400(c.Writer,"传入数据不能为空！")
		return
	}

	if uswrap.AppId==""{
		util.ResponseError400(c.Writer,"app_id不能为空！")
		return
	}

	if uswrap.OpenId == "" {
		util.ResponseError400(c.Writer,"open_id不能为空！")
		return
	}

	//删除用户旧资源
	_,err =db.NewSession().DeleteFrom("qyx_usersource").Where("app_id=? and open_id=?",uswrap.AppId,uswrap.OpenId).Exec()
	if err!=nil{
		util.ResponseError400(c.Writer,"删除用户历史资源失败！")
		return
	}
	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	for _,usersource :=range uswrap.Sources  {
		usersource.AppId = uswrap.AppId
		usersource.OpenId = uswrap.OpenId
		_,err :=tx.InsertInto("qyx_usersource").Columns("app_id","open_id","source_id","action").Record(usersource).Exec()
		if err!=nil{
			tx.Rollback()
			log.Error(err)
			util.ResponseError400(c.Writer,"添加失败！")
			return
		}
	}
	err =tx.Commit()
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加失败！")
		return
	}
	c.JSON(http.StatusOK,uswrap)
}


//查询所有资源
func SourcesAll(c *gin.Context)  {

	c.JSON(http.StatusOK,srsAll)
}

//初始化DB数据
func InitDB() error  {

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id:   "app_init_4",
				Up:   []string{"CREATE TABLE IF NOT EXISTS qyx_appppp(id BIGINT PRIMARY KEY AUTO_INCREMENT," +
					"app_id VARCHAR(50) UNIQUE COMMENT '应用ID'," +
					"open_id VARCHAR(50) DEFAULT '' COMMENT '用户ID'," +
					"source_id VARCHAR(50) DEFAULT '' COMMENT '资源ID'," +
					"`action` VARCHAR(50) DEFAULT '' NOT NULL COMMENT '行为'," +
					"create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'," +
					"update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'" +
					") CHARACTER SET utf8"},
			},
		},
	}

	_, err := migrate.Exec(db.NewSession().DB, "mysql", migrations, migrate.Up)

	return err
}