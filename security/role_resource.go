package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"net/http"
	"strings"
	"errors"
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

	err = InsertRoleResources(roleresources)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"添加失败！")
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


//-----------------------------------  db ----------------------------------------------------

func QueryRoleResources(role,appId string) ([]*RoleResource,error)  {
	var roleresources []*RoleResource
	_,err := db.NewSession().Select("*").From("qyx_role_resource").Where("app_id=? and role=?",appId,role).LoadStructs(&roleresources)
	return roleresources,err
}

func InsertRoleResources(rrss []*RoleResource) error {
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

	for _,rrs :=range rrss {
		err = InsertRoleResourceTx(rrs,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

//查询用户角色资源
func QueryUserRoleResource(openId,appId string) ([]*RoleResource,error)  {
	var roleResources []*RoleResource
	builder:=db.NewSession().Select("qyx_role_resource.*").From("qyx_role").Join("qyx_role_resource","qyx_role_user.role = qyx_role_resource.role").Join("qyx_role_user","qyx_role_user.role = qyx_role_resource.role")
	_,err := builder.Where("qyx_role_user.open_id=? and qyx_role_user.app_id=?",openId,appId).LoadStructs(&roleResources)
	return roleResources,err
}

var userResourceCache = map[string][]*RoleResource{}


func HasResourceWithOpenId(resource string,openId,appId string) bool {

	log.Info("访问应用:",appId)
	log.Info("访问用户:",openId)
	log.Info("访问资源:",resource)
	reosurceActions :=strings.Split(resource,":")
	if len(reosurceActions)!=2 {
		util.CheckErr(errors.New("资源输入有误！"))
		return false
	}
	roleResources := userResourceCache[openId+"-"+appId]
	if roleResources==nil{
		roleResources,err := QueryRoleResources(openId,appId)
		util.CheckErr(err)
		if roleResources!=nil&&len(roleResources)>0{
			userResourceCache[openId+"-"+appId] = roleResources
		}
	}

	if roleResources==nil||len(roleResources)<=0{
		log.Warn("用户没有角色资源！")
		return false
	}

	for _,roleResource :=range roleResources{
		log.Info("roleResource.ResourceId=",roleResource.ResourceId)
		log.Info("roleResource.Action=",roleResource.Action)
		if roleResource.ResourceId == reosurceActions[0] {
			if roleResource.Action == reosurceActions[1]{
				return true
			}
		}
	}

	return false
}


//添加角色资源
func InsertRoleResourceTx(rrs *RoleResource,tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("qyx_role_resource").Columns("app_id","role_id","role","resource_id","action","flag","json").Record(rrs).Exec()

	return err
}