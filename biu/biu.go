package biu

import "github.com/gin-gonic/gin"

type ResponseBody struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func MustInit() {}

func Failed(ctx *gin.Context, statusCode int, msg string) {
	ctx.JSON(statusCode, ResponseBody{
		Code: statusCode,
		Msg:  msg,
	})
}

func ResponseMsg(ctx *gin.Context, statusCode int, code int, msg string) {
	ctx.JSON(statusCode, ResponseBody{
		Code: code,
		Msg:  msg,
	})
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, ResponseBody{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}
