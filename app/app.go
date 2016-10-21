package app

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/rubenv/sql-migrate"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"time"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type App struct {
	AppId string `json:"app_id"`
	AppKey string `json:"app_key"`
	AppName string `json:"app_name"`
	AppDesc string `json:"app_desc"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
	Status int `json:"status"`
	Json string `json:"json"`
	Flag string `json:"flag"`
}

func Setup()  {

	go func() {
		err := InitDB()
		if err!=nil {
			log.Error(err)
			log.Info("初始化APP管理的DB失败！")
			return
		}
		router :=gin.Default()
		router.GET("/v1/apps",AppsAdd)

		log.Info("init app manager success!")

		router.Run(":8081")
	}()
}

// 添加APP
func AppsAdd(c *gin.Context)  {

	var app *App
	err :=c.BindJSON(&app)
	if err!=nil {
		log.Error(err)
		return
	}
	if app.AppId=="" {
		util.ResponseError400(c.Writer,"app_id不能为空！")
		return
	}
	if app.AppName=="" {
		util.ResponseError400(c.Writer,"app_name不能为空！")
		return
	}
	//生成APPKEY
	app.AppKey=util.GenerUUId()
	app.Status = 1 // APP状态为正常

	var existApp *App
	_,err =db.NewSession().Select("*").From("qyx_app").Where("app_id=?",app.AppId).LoadStructs(&existApp)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询APP失败！")
		return
	}

	if existApp!=nil {
		util.ResponseError400(c.Writer,"APP【"+existApp.AppId+"】已存在！")
		return
	}

	//插入APP
	_,err =db.NewSession().InsertInto("qyx_app").Columns("app_id","app_key","app_name","app_desc","status","json","flag").Record(&app).Exec()
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加APP失败！")
		return
	}

	c.JSON(http.StatusOK,app)
}

func AppsWithPage(c *gin.Context)  {


}

//初始化DB数据
func InitDB() error  {

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id:   "app_init",
				Up:   []string{"CREATE TABLE IF NOT EXISTS qyx_app(id BIGINT PRIMARY KEY AUTO_INCREMENT," +
					"app_id VARCHAR(100) UNIQUE COMMENT '应用ID'," +
					"app_key VARCHAR(255) COMMENT '应用KEY'," +
					"app_name VARCHAR(255) COMMENT '应用名称'," +
					"app_desc VARCHAR(1000) COMMENT '应用描述'," +
					"create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'," +
					"update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'," +
					"status int COMMENT '应用状态 0.待审核 1.已审核'," +
					"json VARCHAR(255) COMMENT '附加数据'," +
					"flag VARCHAR(255) COMMENT '标记') CHARACTER SET utf8"},
			},
		},
	}

	_, err := migrate.Exec(db.NewSession().DB, "mysql", migrations, migrate.Up)

	return err
}