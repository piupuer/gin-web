package cache_service

import (
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/thedevsaddam/gojsonq/v2"
)

// 所有的查询可以走redis, 但数据的更新还是走mysql
type RedisService struct {
	mysql service.MysqlService // 保留mysql, 如果没开启redis可以走mysql
	redis *redis.Client        // redis对象实例
}

// 初始化服务
func New(c *gin.Context) RedisService {
	return RedisService{
		mysql: service.New(c),
		redis: global.Redis,
	}
}

// gojsonq.JSONQ不支持FindOne, 由于比较常用, 这里自行实现
func JsonQueryFindOne(s *gojsonq.JSONQ) (interface{}, error) {
	// 获取查询列表
	res := s.Get()
	switch res.(type) {
	case []interface{}:
		v, _ := res.([]interface{})
		if len(v) > 0 {
			// 只取第一条数据
			return v[0], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// jsonq需要每一次新建不同实例, 否则可能查询条件重叠, 找不到记录
func (s RedisService) JsonQuery() *gojsonq.JSONQ {
	return gojsonq.New()
}

// 从缓存中获取model全部数据, 返回json字符串, 参数list为结构体数组, 必须传地址否则可能没数据
func (s RedisService) GetListFromCache(list interface{}, tableName string) string {
	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, tableName)
	res, err := s.redis.Get(cacheKey).Result()
	if err != nil {
		global.Log.Debug(fmt.Sprintf("[GetListFromCache]读取redis缓存异常: %v", err))
	}
	if res == "" {
		// 如果是空字符串, 将其设置为空数组, 否则list会被转为nil
		res = "[]"
	}
	if list != nil {
		utils.Json2Struct(res, list)
	}
	return res
}

// 从缓存中获取model全部数据, 返回json字符串, 参数m为结构体, 必须传地址否则可能没数据
func (s RedisService) GetItemByIdFromCache(id uint, m interface{}, tableName string) error {
	json := s.GetListFromCache(nil, tableName)
	// id在JSONQ中以int存在
	res, err := JsonQueryFindOne(s.JsonQuery().FromString(json).Where("id", "=", int(id)))
	if err != nil {
		return err
	}
	// 转换为结构体
	utils.Struct2StructByJson(res, m)
	return nil
}
