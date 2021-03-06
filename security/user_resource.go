package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"net/http"
)

type UserResource struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//角色
	OpenId string `json:"open_id"`
	//资源ID
	ResourceId string `json:"resource_id"`
	Action  string `json:"action"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

func UserResourceList(c *gin.Context)  {
	openId :=c.Param("open_id")
	appId :=c.Param("app_id")

	urs,err := QueryUserResource(openId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询用户资源失败！")
		return
	}

	if urs==nil || len(urs)<=0 {
		urs = []*UserResource{}
	}
	c.JSON(http.StatusOK,urs)
}

func QueryUserResource(openId,appId string) ([]*UserResource,error)  {
	var urs []*UserResource
	_,err :=db.NewSession().Select("distinct qyx_role_resource.app_id,qyx_role_resource.resource_id,qyx_role_resource.action,qyx_role_user.open_id").From("qyx_role_user").Join("qyx_role_resource","qyx_role_user.role=qyx_role_resource.role and qyx_role_user.app_id=qyx_role_resource.app_id").Where("qyx_role_user.open_id=? and qyx_role_user.app_id=?",openId,appId).LoadStructs(&urs)
	return urs,err
}