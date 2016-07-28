package log

import (
	 "github.com/Sirupsen/logrus"
	"os"
)

func init()  {
	env := os.Getenv("GO_ENV")
	if env=="production" {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(logrus.InfoLevel)
	}else{
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.DebugLevel)
	}
	
}

func Info(args ...interface{})  {
	
	logrus.Info(args)
}

func Debug(args ...interface{}) {

	logrus.Debug(args)
}

func Warn(args ...interface{}) {

	logrus.Warn(args)
}

func Error(args ...interface{}) {

	logrus.Error(args)
}