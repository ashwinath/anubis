package logger

import (
	"fmt"
	"log"
	"os"
)

const logFormat = log.Ldate | log.Ltime

var infoLogger = log.New(os.Stdout, "INFO: ", logFormat)
var warningLogger = log.New(os.Stdout, "WARN: ", logFormat)
var errorLogger = log.New(os.Stderr, "ERROR: ", logFormat)

func Infof(msg string, args ...interface{}) {
	infoLogger.Println(fmt.Sprintf(msg, args...))
}

func Warningf(msg string, args ...interface{}) {
	warningLogger.Println(fmt.Sprintf(msg, args...))
}

func Errorf(msg string, args ...interface{}) {
	errorLogger.Println(fmt.Sprintf(msg, args...))
}
