package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"net/http"
)

type Role struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	//角色标识
	Role string `json:"role"`
	//名称
	Name string `json:"name"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

//添加角色
func RoleAdd(c *gin.Context)  {

	var roles []*Role
	err :=c.BindJSON(&roles)
	if err!=nil||(roles==nil||len(roles)<=0){
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误！")
		return
	}

	//插入角色
	err =InsertRoles(roles)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"服务插入失败！")
		return
	}

	util.ResponseSuccess(c.Writer)
}

//删除角色
func RoleDel(c *gin.Context)  {
	role :=c.Param("role")
	appId :=c.Param("app_id")

	err :=DeleteRole(role,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"删除失败！")
		return
	}
	util.ResponseSuccess(c.Writer)
	return
}

//查询角色
func RoleList(c *gin.Context)  {

	appId :=c.Query("app_id")

	roles,err :=QueryRoles(appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询失败！")
		return
	}

	if roles==nil||len(roles)<=0{
		roles = []*Role{}
	}
	c.JSON(http.StatusOK,roles)
}

func DeleteRole(role,appId string) error  {

	_,err :=db.NewSession().DeleteFrom("qyx_role").Where("role=? and app_id=?",role,appId).Exec()

	return err
}

//查询角色
func QueryRoles(appId string) ([]*Role,error)  {

	var roles []*Role
	_,err :=db.NewSession().Select("*").From("qyx_role").Where("app_id=?",appId).LoadStructs(&roles)

	return roles,err
}

//是否有角色
func HasRoles(roles []string,openId string,appId string) bool  {
	var count int64
	err :=db.NewSession().Select("count(*)").From("qyx_role_user").Where("open_id=?",openId).Where("role in ?",roles).Where("app_id=?",appId).LoadValue(&count)
	util.CheckErr(err)
	if count>0 {
		return true
	}
	return false
}

//是否有角色
func HasRole(role string,openId string,appId string) bool  {
	var count int64
	err :=db.NewSession().Select("count(*)").From("qyx_role_user").Where("open_id=?",openId).Where("role=?",role).Where("app_id=?",appId).LoadValue(&count)
	util.CheckErr(err)
	if count>0 {
		return true
	}
	return false
}

func InsertRoles(roles []*Role) error {

	tx,err :=db.NewSession().Begin()
	if err!=nil{
		return err
	}
	defer func() {
		if err := recover();err!=nil{
			tx.Rollback()
			panic(err)
		}
	}()
	for _,role :=range roles {
		err = InsertRoleTx(role,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
	}

	err=tx.Commit()
	return err

}

func InsertRoleTx(role *Role,tx *dbr.Tx) error {

	_,err :=tx.InsertInto("qyx_role").Columns("app_id","role","name","flag","json").Record(role).Exec()

	return err
}