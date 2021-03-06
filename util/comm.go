package util

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"bytes"
	"github.com/sumory/idgen"
	"hash"
	"sort"
	"encoding/hex"
	"crypto/md5"
	"bufio"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"time"
	"math/rand"
	"errors"
)

const (
	Error_Code_OK =0
)



func CheckErr(err error)  {
	if err != nil {
		panic(err)
	}
}

func ResponseError400AndForward(w http.ResponseWriter,msg string,forward string){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := ResultErrorAndForward{http.StatusBadRequest, msg,forward}
	log.Error(msg)
	w.WriteHeader(http.StatusBadRequest)
	WriteJson(w,err)
}

func ResponseError400(w http.ResponseWriter,msg string){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ResponseError(w,http.StatusBadRequest,msg)
}
func ResponseError401(w http.ResponseWriter,msg string){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ResponseError(w,http.StatusUnauthorized,msg)
}

func ResponseError401Msg(w http.ResponseWriter){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ResponseError(w,http.StatusUnauthorized,"认证失败")
}

func ResponseError(w http.ResponseWriter, statusCode int,msg string)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := ResultError{statusCode, msg}
	log.Error(msg)
	w.WriteHeader(statusCode)
	WriteJson(w,err)
}

func ResponseErrorS(w http.ResponseWriter, statusCode int,msg string)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := ResultError{statusCode, msg}

	w.WriteHeader(http.StatusBadRequest)
	WriteJson(w,err)
}

func ResponseErrorSS(w http.ResponseWriter,httpStatus, statusCode int,msg string)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := ResultError{statusCode, msg}

	w.WriteHeader(httpStatus)
	WriteJson(w,err)
}

func ResponseSuccess(w http.ResponseWriter)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := NewResultError(0,"OK")
	WriteJson(w,err)

}

//将对象转换为JSON
func ToJson(obj interface{})  (string,error){
	jsonData,err:= json.Marshal(obj);

	if err!=nil {
		return "",err
	}

	return string(jsonData),nil
}

func ToJson2(obj interface{})  (string){
	json,err := ToJson(obj)
	CheckErr(err)
	return json
}
func WriteJsonStr(w io.Writer,json string) {
	if json=="" {
		io.WriteString(w,"{}")
		return
	}
	io.WriteString(w,json)
}

func WriteJson(w io.Writer,obj interface{})  {

	if obj==nil {
		io.WriteString(w,"{}")
		return
	}

	if objStr,ok := obj.(string);ok {
		if objStr=="" {
			io.WriteString(w,"{}")
			return
		}
	}
	jsonData,_:= json.Marshal(obj);

	io.WriteString(w,string(jsonData))
}

func ReadJsonByByte(body []byte,obj interface{}) error {
	mdz:=json.NewDecoder(bytes.NewBuffer(body))

	mdz.UseNumber()
	err := mdz.Decode(obj)

	if  err != nil {
		return err;
	}
	return nil;
}

func ReadJson( r io.ReadCloser,obj interface{})  error {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	if err := r.Close(); err != nil {
		panic(err)
	}


	return ReadJsonByByte(body,obj);


}

func GenerUUId()  string{

	out, _ := exec.Command("uuidgen").Output()


	return strings.Replace(strings.TrimSpace(string(out)),"-","",-1)
}

//生成APPID
func GenerAppId() int64  {
	err, idWorker := idgen.NewIdWorker(1)
	CheckErr(err)
	err,appid := idWorker.NextId()
	CheckErr(err)
	return appid;
}

type ResultError struct {

	ErrCode int `json:"err_code"`
	ErrMsg string `json:"err_msg"`

}

type ResultErrorAndForward struct {

	ErrCode int `json:"err_code"`
	ErrMsg string `json:"err_msg"`
	//跳转
	Forward string `json:"forward,omitempty"`

}

func (self *ResultError) Success() bool {

	return self.ErrCode==0 || self.ErrCode==200
}


func NewResultError(errCode int,errMsg string) *ResultError  {

	resultError := &ResultError{}
	resultError.ErrCode=errCode;
	resultError.ErrMsg=errMsg

	return  resultError
}

func SignWithBaseSign(params map[string]interface{}, apiKey string,basesign string, fn func() hash.Hash) string {
	if fn == nil {
		fn = md5.New
	}
	h := fn()
	bufw := bufio.NewWriterSize(h, 128)

	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)


	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		vs := ObjToStr(v)
		bufw.WriteString(k)
		bufw.WriteByte('=')
		bufw.WriteString(vs)
		bufw.WriteByte('&')
	}
	bufw.WriteString("key=")
	bufw.WriteString(apiKey)

	if basesign!=""{
		bufw.WriteString("&")
		bufw.WriteString("basesign=")
		bufw.WriteString(basesign)
	}

	bufw.Flush()
	signature := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// Sign 支付签名.
//  params: 待签名的参数集合
//  apiKey: api密钥
//   basesign 基础sign
//  fn:     func() hash.Hash, 如果为 nil 则默认用 md5.New
func Sign(params map[string]string, apiKey string, fn func() hash.Hash) string {

	objparams :=make(map[string]interface{})

	for key,_ :=range params {

		objparams[key] = params[key]
	}

	return SignWithBaseSign(objparams,apiKey,"",fn)
}

// 基础签名
func SignBase(appKey string) (timestamp,noncestr,basesign string ) {
	noncestr =GetRandomSalt()
	timestamp =fmt.Sprintf("%d",time.Now().Unix())
	signStr := appKey+noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign =fmt.Sprintf("%X",bytes)
	return
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(num int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < num; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ObjToStr(v interface{}) string {
	var strV string
	switch v.(type) {

	case int:
		strV= fmt.Sprintf("%d",v)
		break
	case uint:
		strV= fmt.Sprintf("%d",v)
		break
	case int64:
		strV= fmt.Sprintf("%d",v)
		break
	case uint64:
		strV= fmt.Sprintf("%d",v)
		break
	case int8:
		strV= fmt.Sprintf("%d",v)
		break
	case uint8:
		strV= fmt.Sprintf("%d",v)
		break
	case int16:
		strV= fmt.Sprintf("%d",v)
		break
	case uint16:
		strV= fmt.Sprintf("%d",v)
		break
	case int32:
		strV= fmt.Sprintf("%d",v)
		break
	case uint32:
		strV= fmt.Sprintf("%s",v)
		break
	case string:
		strV= fmt.Sprintf("%s",v)
		break
	case float32:
		strV= fmt.Sprintf("%s",v)
		break
	case float64:
		strV= fmt.Sprintf("%s",v)
		break
	default:
		strV= fmt.Sprintf("%s",v)

	}
	return strV
}


//获取API请求的返回错误日志
func GetApiResponseErr(status int,payload []byte,er error)  error {
	if status==http.StatusBadRequest { //400 表示服务器有话要说
		var errMap map[string]interface{}
		err := ReadJsonByByte(payload,&errMap)
		if err!=nil{
			log.Error(err)
			return errors.New("服务返回数据有误！")
		}
		errMsg := errMap["err_msg"]
		if errMsg!=nil {
			return errors.New(errMsg.(string))
		}
	}
	if status == http.StatusOK {
		return nil
	}
	if er!=nil{
		log.Error(er)
		return errors.New("请求API出错！")
	}
	return nil

}