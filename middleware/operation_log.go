package middleware

import (
	v1 "gin-web/api/v1"
	"gin-web/models"
	"gin-web/pkg/cache_service"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 操作日志
func OperationLog(c *gin.Context) {
	// 开始时间
	startTime := time.Now()
	// 避免服务器出现异常, 这里用defer保证一定可以执行
	defer func() {
		// 下列请求比较频繁无需写入日志
		if c.Request.Method == http.MethodGet || 
			c.Request.Method == http.MethodOptions {
			return
		}
		// 结束时间
		endTime := time.Now()

		log := models.SysOperationLog{
			// Ip地址
			Ip: c.ClientIP(),
			// 请求方式
			Method: c.Request.Method,
			// 请求路径(去除url前缀)
			Path: strings.TrimPrefix(c.Request.URL.Path, "/"+global.Conf.System.UrlPathPrefix),
			// 请求耗时
			Latency: endTime.Sub(startTime),
			// 浏览器标识
			UserAgent: c.Request.UserAgent(),
		}

		// 获取当前登录用户
		user := v1.GetCurrentUser(c)

		// 用户名
		if user.Id > 0 {
			log.Username = user.Username
			log.RoleName = user.Role.Name
		} else {
			log.Username = "未登录"
			log.RoleName = "未登录"
		}

		// 获取当前接口
		cache := cache_service.New(c)
		apis, err := cache.GetApis(&request.ApiListRequestStruct{
			Method: log.Method,
			Path:   log.Path,
		})
		if err == nil && len(apis) == 1 {
			log.ApiDesc = apis[0].Desc
		} else {
			log.ApiDesc = "无"
		}
		
		// 获取Ip所在地
		log.IpLocation = "未知地址"

		// 响应状态码
		log.Status = c.Writer.Status()
		// 响应数据
		resp, exists := c.Get(global.Conf.System.OperationLogKey)
		var data string
		if exists {
			data = utils.Struct2Json(resp)
		} else {
			data = "无"
		}
		// 通过gzip压缩
		log.Data = utils.Str2BytesByGzip(data)
		// 写入数据库
		global.Mysql.Create(&log)
	}()
	c.Next()
}
