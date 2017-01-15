package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"fmt"
	"errors"
	"net/http"
)

var client *raven.Client

func Setup(dsn string) error  {
	return raven.SetDSN(dsn)
}

//重要错误
func CaptureMajorErr(errStr string,flag string)  {
	packet := raven.NewPacket(errStr, raven.NewException(errors.New(errStr), raven.NewStacktrace(2, 3, nil)))
	raven.Capture(packet, map[string]string{
		"type": flag,
	})

}

func CaptureErr(err error)  {

	raven.CaptureError(err,nil)
}

func Recovery(onlyCrashes bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {


			flags := map[string]string{
				"endpoint": c.Request.RequestURI,
			}
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rvalStr := fmt.Sprint(rval)
				packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
				raven.Capture(packet, flags)
				c.Writer.WriteHeader(http.StatusInternalServerError)
			}
			if !onlyCrashes {
				for _, item := range c.Errors {
					packet := raven.NewPacket(item.Error(), &raven.Message{item.Error(), []interface{}{item.Meta}})
					client.Capture(packet, flags)
				}
			}
		}()
		c.Next()
	}
}
