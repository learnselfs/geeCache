// Package utils @Author Bing
// @Date 2023/10/25 11:47:00
// @Desc
package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "【utils】", log.Ldate|log.Ltime|log.LstdFlags|log.Llongfile)

}

func Log(fmt string, v ...interface{}) {
	Logger.Printf(fmt, v...)
}
