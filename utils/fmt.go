// Package utils @Author Bing
// @Date 2023/10/25 16:24:00
// @Desc
package utils

import (
	"net/http"
)

type Obj struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Data   H      `json:"data"`
}

type H map[string]interface{}

func obj(code int, msg string, data H) H {
	Log("code: %d, msg: %s", code, msg)
	return H{"status": http.StatusText(code), "code": code, "msg": msg, "data": data}
}

func Ok() H {
	return obj(http.StatusOK, "", nil)
}

func OkWithMsg(code int, msg string) H {
	return obj(code, msg, nil)
}

func OkWithDetail(code int, msg string, data H) H {
	return obj(code, msg, data)
}

func Failed() H {
	return obj(http.StatusInternalServerError, "", nil)
}

func FailedWithMsg(code int, msg string) H {
	return obj(code, msg, nil)
}
