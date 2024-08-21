package ctl

import (
	"micro_todoList/pkg/e"

	"github.com/gin-gonic/gin"
)

// Response 基础序列化器
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// 带数组形式的返回

type DataList struct {
	Item  interface{} `json:"item"`
	Total int64       `json:"total"`
}

// 用户登陆专属，带用户信息和token的返回resp

type TokenData struct {
	User        interface{} `json:"user"`
	AccessToken string      `json:"access_token"`
}

// 带有追踪信息的错误返回

type TrackedErrorResponse struct {
	Response
	TraceId string `json:"trace_id"`
}

// RespSuccess 带data成功返回
func RespSuccess(ctx *gin.Context, data interface{}, code ...int) *Response {
	status := e.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = "操作成功"
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    e.GetMsg(status),
	}

	return r
}

func RespError(ctx *gin.Context, err error, data string, code ...int) *Response {
	status := e.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    e.GetMsg(status),
		Error:  err.Error(),
	}

	return r
}
