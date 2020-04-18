package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// http请求响应封装
type Resp struct {
	Code int         `json:"code"` // 错误代码代码
	Data interface{} `json:"data"` // 数据内容
	Msg  string      `json:"msg"`  // 消息提示
}

// 带分页信息的响应
type PageResp struct {
	Code  int         `json:"code"`  // 错误代码代码
	Data  interface{} `json:"data"`  // 数据内容
	Msg   string      `json:"msg"`   // 消息提示
	Total int         `json:"total"` // 数据条数
}

const (
	SUCCESS   = 201
	FAIL      = 405
	EXCEPTION = 500
)

var errorMsg = map[int]string{
	SUCCESS:   "操作成功",
	FAIL:      "操作失败",
	EXCEPTION: "系统异常",
}

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func ResultPage(code int, msg string, data interface{}, total int, c *gin.Context) {
	c.JSON(http.StatusOK, PageResp{
		Code:  code,
		Data:  data,
		Msg:   msg,
		Total: total,
	})
}

func Success(c *gin.Context) {
	Result(SUCCESS, errorMsg[SUCCESS], map[string]interface{}{}, c)
}

func SuccessWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, errorMsg[SUCCESS], data, c)
}

func SuccessWithMsg(msg string, c *gin.Context) {
	Result(SUCCESS, msg, map[string]interface{}{}, c)
}

func SuccessPageWithData(data interface{}, total int, c *gin.Context) {
	ResultPage(SUCCESS, errorMsg[SUCCESS], data, total, c)
}

func Fail(c *gin.Context) {
	Result(FAIL, errorMsg[FAIL], map[string]interface{}{}, c)
}

func FailWithMsg(msg string, c *gin.Context) {
	Result(FAIL, msg, map[string]interface{}{}, c)
}

func FailWithCode(code int, c *gin.Context) {
	Result(code, errorMsg[code], map[string]interface{}{}, c)
}
