package app

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/rubenv/sql-migrate"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type App struct {
	AppId string
	AppKey string
	AppName string
	AppDesc string
}

func Setup(router gin.IRouter) error {
	err := InitDB()
	if err!=nil {
		log.Error(err)
		return err
	}
	log.Info("init......")
	router.GET("/v1/apps",Apps)

	return nil

}

func Apps(c *gin.Context)  {
	log.Info("测试")
	c.JSON(http.StatusOK,map[string]string{
		"test": "122",
	})
}

//初始化DB数据
func InitDB() error  {

	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id:   "app_init",
				Up:   []string{"CREATE TABLE IF NOT EXISTS app(id BIGINT PRIMARY KEY AUTO_INCREMENT," +
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