package security

import (
	"github.com/gin-gonic/gin"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"github.com/gocraft/dbr"
	"net/http"
)

type RoleUser struct {
	Id int64 `json:"id"`
	AppId string `json:"app_id"`
	Role string `json:"role"`
	OpenId string `json:"open_id"`
	Flag string `json:"flag"`
	Json string `json:"json"`
}

//添加用户角色
func RoleUserAdd(c *gin.Context)  {

	appId :=c.Param("app_id")
	openId :=c.Param("open_id")

	var roleUsers []*RoleUser
	err :=c.BindJSON(&roleUsers)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"数据格式有误！")
		return
	}
	if roleUsers!=nil&&len(roleUsers)>0{
		for _,ru :=range roleUsers{
			ru.AppId = appId
			ru.OpenId = openId
		}
	}

	err = InsertRoleUsers(roleUsers)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"插入用户角色失败！")
		return
	}

	util.ResponseSuccess(c.Writer)
}

//查询用户角色
func RoleUserList(c *gin.Context)  {
	appId :=c.Param("app_id")
	openId := c.Param("open_id")

	rus,err := QueryRoleUsers(openId,appId)
	if err!=nil{
		log.Error(err)
		util.ResponseError400(c.Writer,"查询用户角色失败！")
		return
	}

	if rus==nil||len(rus)<=0 {
		rus = []*RoleUser{}
	}

	c.JSON(http.StatusOK,rus)
}




//-------------------------------------  private ------------------------------



//-------------------------------------  db ------------------------------

//删除用户的角色
func DeleteRoleUser(role,openId,appId string,tx *dbr.Tx) error {
	_,err := tx.DeleteFrom("qyx_role_user").Where("role=? and open_id=? and app_id=?",role,openId,appId).Exec()

	return err
}

func QueryRoleUsers(openId,appId string) ([]*RoleUser,error)  {
	var rus  []*RoleUser
	_,err :=db.NewSession().Select("*").From("qyx_role_user").Where("open_id=? and app_id=?",openId,appId).LoadStructs(&rus)
	return rus,err
}

func InsertRoleUsers(roleusers []*RoleUser) error {

	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err !=nil{
			tx.Rollback()
			panic(err)
		}
	}()

	for _,ru :=range roleusers {
		err := InsertRoleUser(ru,tx)
		if err!=nil{
			log.Error(err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func InsertRoleUser(roleuser *RoleUser,tx *dbr.Tx) (error)  {

	_,err :=tx.InsertInto("qyx_role_user").Columns("app_id","role","open_id","flag","json").Record(roleuser).Exec()

	return err
}