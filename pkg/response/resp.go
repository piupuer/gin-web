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

// 分页封装
type PageInfo struct {
	PageNum  uint `json:"pageNum" form:"pageNum"`   // 当前页码
	PageSize uint `json:"pageSize" form:"pageSize"` // 每页显示条数
	Total    uint `json:"total"`                    // 数据总条数
}

// 带分页数据封装
type PageData struct {
	PageInfo
	List interface{} `json:"list"` // 数据列表
}

// 计算limit/offset, 如果需要用到返回的PageSize, PageNum, 务必保证Total值有效
func (s *PageInfo) GetLimit() (limit uint, offset uint) {
	// 传入参数可能不合法, 设置默认值
	// 每页显示条数不能小于1
	if s.PageSize < 1 {
		s.PageSize = 10
	}
	// 页码不能小于1
	if s.PageNum < 1 {
		s.PageNum = 1
	}

	// 如果偏移量比总条数还多
	if s.Total > 0 {
		if s.PageSize > s.Total {
			s.PageSize = s.Total
		}
		if s.PageNum > s.Total {
			s.PageNum = s.Total
		}
	}

	// 计算最大页码
	maxPageNum := s.Total/s.PageSize + 1
	if s.Total%s.PageSize == 0 {
		maxPageNum = s.Total / s.PageSize
	}
	// 页码不能小于1
	if maxPageNum < 1 {
		maxPageNum = 1
	}

	// 超出最后一页
	if s.PageNum > maxPageNum {
		s.PageNum = maxPageNum
	}

	limit = s.PageSize
	offset = limit * (s.PageNum - 1)
	return
}

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Success(c *gin.Context) {
	Result(c, Ok, CustomError[Ok], map[string]interface{}{})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	Result(c, Ok, CustomError[Ok], data)
}

func SuccessWithMsg(c *gin.Context, msg string) {
	Result(c, Ok, msg, map[string]interface{}{})
}

func Fail(c *gin.Context) {
	FailWithCode(c, NotOk)
}

func FailWithMsg(c *gin.Context, msg string) {
	Result(c, NotOk, msg, map[string]interface{}{})
}

func FailWithCode(c *gin.Context, code int) {
	// 查找给定的错误码存在对应的错误信息, 默认使用NotOk
	msg := CustomError[NotOk]
	if val, ok := CustomError[code]; ok {
		msg = val
	}
	Result(c, code, msg, map[string]interface{}{})
}
