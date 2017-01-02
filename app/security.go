package app

import (
	"net/http"
	"errors"
	"strconv"
	"fmt"
	"crypto/md5"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"time"
)

const APP_ID_KEY  = "app_id"

type AppSign struct {
	App *App
	Sign string
}

var cacheAppMap map[string]*App

// 设置缓存中的APP
func CacheAppSet(app *App)  {
	if cacheAppMap==nil{
		cacheAppMap = make(map[string]*App,0)
	}

	cacheAppMap[app.AppId] = app
}

//获取缓存中的APP
func CacheAppWithAppId(appId string) *App  {
	if cacheAppMap==nil{
		cacheAppMap = make(map[string]*App,0)
	}

	return cacheAppMap[appId]
}

// 认证APP信息
func Auth(req *http.Request) (*AppSign,error)  {
	appId := GetParamInRequest("app_id",req)
	if appId==""{
		return nil,errors.New("app_id不能为空")
	}
	//从缓存中获取APP信息
	app :=CacheAppWithAppId(appId)
	if app==nil{
		 dbApp,err := QueryAppWithId(appId)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
		if dbApp!=nil{
			//缓存APP信息
			CacheAppSet(dbApp)
		}
		app = dbApp
	}

	if app==nil{
		return nil,errors.New("应用信息未找到!请检查APPID是否正确");
	}
	sign :=GetParamInRequest("sign",req)
	if sign =="" {

		return nil,errors.New("签名信息(sign)不能为空!");
	}
	gotSign := sign

	noncestr :=GetParamInRequest("noncestr",req)
	timestamp :=GetParamInRequest("timestamp",req)

	if noncestr=="" {
		return nil,errors.New("随机码不能为空!");
	}

	if timestamp=="" {
		return nil,errors.New("时间戳不能为空!");
	}


	timestam64,_ := strconv.ParseInt(timestamp,10,64)
	timeBtw := time.Now().Unix()-int64(timestam64)
	if timeBtw > 5*60 {

		return nil,errors.New("签名已失效!");
	}

	signStr:= fmt.Sprintf("%s%s%s",app.AppKey,noncestr,timestamp)
	log.Info("signStr=",signStr)
	wantSign :=fmt.Sprintf("%X",md5.Sum([]byte(signStr)))

	if gotSign!=wantSign {
		fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
		return nil,errors.New("请求不合法!");
	}

	appSign := &AppSign{App:app,Sign:wantSign}

	return appSign,nil;
}

func QueryAppWithId(id string) (*App,error)  {
	var dbApp *App
	_,err := db.NewSession().Select("id","app_id","app_key","app_name","app_desc","status").From("qyx_app").Where("app_id=? and status=?",id,"1").LoadStructs(&dbApp)
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	return dbApp,nil
}

//在请求中获取AppId
func GetParamInRequest(key string,req *http.Request) string  {
	var value string = req.Header.Get(key)
	if value=="" {
		if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
			value = values[0]
		}
	}



	return value

}

//获取APPID
func GetAppIdInRequest(req *http.Request) string {

	return GetParamInRequest(APP_ID_KEY,req)
}