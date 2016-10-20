package app

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"github.com/gin-gonic/gin"
)

type App struct {
	AppId string
	AppKey string
	AppName string
	AppDesc string
}

func Setup(router gin.IRouter) error {
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