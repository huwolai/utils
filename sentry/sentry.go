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

func Setup(dsn string, tags map[string]string) error  {
	var err error
	client,err = raven.NewWithTags(dsn,tags)
	if err!=nil{
		return err
	}
	return nil
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
				client.Capture(packet, flags)
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
