package security

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"errors"
)


type JWTAuthenticationBackend struct {
	PublicKey  *rsa.PublicKey
}

type AuthUser struct  {
	//openid
	OpenId string
	//关联ID （第三方ID）
	Rid string

}
const (
	tokenDuration = 72
	expireOffset  = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//认证用户信息
func AuthUsers(req *http.Request) (*AuthUser,error) {

	return GetAuthUser(req)
}

//认证用户信息 并且判断openId是否是当前用户的openID
func AuthUsersAndOpenId(openId string,req *http.Request) (*AuthUser,error) {
	authuser,err :=GetAuthUser(req)
	if err!=nil{
		return authuser,err
	}
	if authuser.OpenId !=openId {
		return nil,errors.New("不是当前用户！")
	}
	return authuser,nil
}

//获取认证用户信息
func GetAuthUser(req *http.Request) (*AuthUser,error)  {
	token :=GetParamInRequest("Authorization",req)
	if token=="" {
		log.Error("没有认证信息!")
		return nil,errors.New("没有认证信息!")
	}
	jwttoken,err :=InitJWTAuthenticationBackend().FetchToken(token)
	if err!=nil{
		log.Error("解析认证信息失败:",err)
		return nil,err
	}
	if !jwttoken.Valid {
		log.Error("认证信息无效!")
		return nil,errors.New("认证信息无效!")
	}
	authUser :=&AuthUser{}
	authUser.OpenId = jwttoken.Claims["sub"].(string)
	authUser.Rid = jwttoken.Claims["r_id"].(string)
	return authUser,nil
}


func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend)  FetchToken(authorization string) (token *jwt.Token,err error){
	token, err =jwt.Parse(authorization, func(token *jwt.Token)(interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return backend.PublicKey, nil
	})
	return token,err;
}




func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(config.GetValue("publickey_path").ToString())
	if err != nil {
		log.Error(err)
		panic(err)
	}
	log.Info("读取到公钥:",publicKeyFile)
	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	publicKeyFile.Close()
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		panic(err)
	}
	return rsaPub
}
