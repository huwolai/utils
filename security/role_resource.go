package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"net/http"
)

type RoleResource struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//角色
	Role string `json:"role"`
	//资源ID
	ResourceId string `json:"resource_id"`
	Action  string `json:"action"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

//添加角色资源
func RoleResourceAdd(c *gin.Context)  {

	var roleresources []*RoleResource
	err :=c.BindJSON(&roleresources)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误！")
		return
	}

	role :=c.Param("role")
	appId :=c.Param("app_id")

	if roleresources!=nil&&len(roleresources)>0{
		for _,rrs :=range roleresources {
			rrs.AppId = appId
			rrs.Role = role
		}
	}

	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	err =InsertRoleResourceTx(roleresources,tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		util.ResponseError400(c.Writer,"添加失败！")
		return
	}

	err = tx.Commit()
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"提交失败！")
		return
	}

	util.ResponseSuccess(c.Writer)
}

//查询角色资源
func RoleResourceList(c *gin.Context)  {

	appId := c.Param("app_id")
	role :=c.Param("role")
	if appId==""{
		util.ResponseError400(c.Writer,"app_id不能为空！")
		return
	}
	if role==""{
		util.ResponseError400(c.Writer,"role不能为空！")
		return
	}

	roleresources,err := QueryRoleResources(role,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询角色资源失败！")
		return
	}
	c.JSON(http.StatusOK,roleresources)
}

func QueryRoleResources(role,appId string) ([]*RoleResource,error)  {
	var roleresources []*RoleResource
	_,err := db.NewSession().Select("*").From("role_resource").Where("app_id=? and role=?",appId,role).LoadStructs(&roleresources)
	return roleresources,err
}

func InsertRoleResources(rrss []*RoleResource) error {
	tx,err :=db.NewSession().Begin()
	if err!=nil{
		return err
	}
	defer func() {
		if err = recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()

	for _,rrs :=range rrss {
		err = InsertRoleResourceTx(rrs,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func InsertRoleResourceTx(rrs *RoleResource,tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("role_resource").Columns("app_id","role_id","role","resource_id","action","flag","json").Record(rrs).Exec()

	return err
}