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

// jsonq需要每一次新建不同实例, 否则可能查询条件重叠, 找不到记录
func (s RedisService) JsonQuery() *gojsonq.JSONQ {
	return gojsonq.New()
}

// 从缓存中获取model全部数据, 返回json字符串, 参数list为结构体数组, 必须传地址否则可能没数据
func (s RedisService) GetListFromCache(tableName string, list interface{}) (string, error) {
	// 缓存键由数据库名与表名组成
	cacheKey := fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, tableName)
	res, err := s.redis.Get(cacheKey).Result()
	if list != nil {
		utils.Json2Struct(res, list)
	}
	return res, err
}

// 从缓存中获取model全部数据, 返回json字符串, 参数m为结构体, 必须传地址否则可能没数据
func (s RedisService) GetItemByIdFromCache(id uint, m interface{}, tableName string) error {
	json, err := s.GetListFromCache(tableName, nil)
	if err != nil {
		return err
	}
	var list []interface{}
	// id在JSONQ中以int存在
	res := s.JsonQuery().FromString(json).Where("id", "=", int(id)).Get()
	if m != nil {
		// 转换为结构体
		utils.Struct2StructByJson(res, &list)
		if len(list) == 0 {
			return gorm.ErrRecordNotFound
		} else {
			utils.Struct2StructByJson(list[0], m)
		}
	}
	return nil
}
