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

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func ResultPage(c *gin.Context, code int, msg string, data interface{}, total int) {
	c.JSON(http.StatusOK, PageResp{
		Code:  code,
		Data:  data,
		Msg:   msg,
		Total: total,
	})
}

func Success(c *gin.Context) {
	Result(c, SUCCESS, errorMsg[SUCCESS], map[string]interface{}{})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, errorMsg[SUCCESS], data)
}

func SuccessWithMsg(c *gin.Context, msg string) {
	Result(c, SUCCESS, msg, map[string]interface{}{})
}

func SuccessPageWithData(c *gin.Context, data interface{}, total int) {
	ResultPage(c, SUCCESS, errorMsg[SUCCESS], data, total)
}

func Fail(c *gin.Context) {
	Result(c, FAIL, errorMsg[FAIL], map[string]interface{}{})
}

func FailWithMsg(c *gin.Context, msg string) {
	Result(c, FAIL, msg, map[string]interface{}{})
}

func FailWithCode(c *gin.Context, code int) {
	Result(c, code, errorMsg[code], map[string]interface{}{})
}
