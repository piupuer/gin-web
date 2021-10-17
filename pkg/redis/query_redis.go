package redis

import (
	"errors"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/thedevsaddam/gojsonq/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"sync"
)

// 以redis client为基础, 加入一些类似gorm的查询方法(目前语法为gorm v2.0)
// 核心是将Preload方法同步过来, 使查询更方便

// 自定义redis查询结构
type QueryRedis struct {
	ctx *gin.Context
	// 错误信息
	Error error
	// redis对象
	redis redis.UniversalClient
	// 是否需要克隆
	clone int
	// 查询声明, 类似于之前版本的search
	Statement      *Statement
	cacheStore     *sync.Map
	NamingStrategy schema.Namer
}

// 初始化服务
func New(c *gin.Context) *QueryRedis {
	nc:= gin.Context{}
	if c != nil {
		nc = *c
	}
	return &QueryRedis{
		ctx:   &nc,
		redis: global.Redis,
		// 初始化时应该为1, 链式操作时才能clone
		clone: 1,
	}
}

// AddError add error to db
func (s *QueryRedis) AddError(err error) error {
	if s.Error == nil {
		s.Error = err
	} else if err != nil {
		s.Error = fmt.Errorf("%v; %w", s.Error, err)
	}
	return s.Error
}

func (s QueryRedis) getInstance() *QueryRedis {
	if s.clone > 0 {
		tx := &QueryRedis{
			ctx:   s.ctx,
			redis: s.redis,
			// 表命名相关(主要针对表前缀)
			NamingStrategy: global.Mysql.NamingStrategy,
		}

		if s.clone == 1 {
			// clone with new statement
			tx.Statement = &Statement{
				DB: tx,
			}
		} else {
			// gorm clone>1用于共享条件, 参考https://gorm.io/zh_CN/docs/session.html#WithConditions
			// 为降低复杂度这一步省略
		}

		return tx
	}

	return &s
}

// 校验表名是否正常, 有可能没有执行Table方法
func (s QueryRedis) check() bool {
	ins := s.getInstance()
	// 未指定json字符串时考虑表名称
	if !ins.Statement.json && strings.TrimSpace(ins.Statement.Table) == "" {
		s.Error = fmt.Errorf("invalid table name: '%s'", ins.Statement.Table)
		return false
	}
	return true
}

// 执行查询前的一些初始化操作
func (s QueryRedis) beforeQuery(db *QueryRedis) *QueryRedis {
	stmt := db.Statement
	if stmt.Model == nil {
		stmt.Model = stmt.Dest
	} else if stmt.Dest == nil {
		stmt.Dest = stmt.Model
	}

	if stmt.Model != nil && !stmt.count {
		if err := stmt.Parse(stmt.Model); err != nil && (!errors.Is(err, schema.ErrUnsupportedDataType) || (stmt.Table == "")) {
			db.AddError(err)
		}
	}

	if stmt.Dest != nil {
		stmt.ReflectValue = reflect.ValueOf(stmt.Dest)
		for stmt.ReflectValue.Kind() == reflect.Ptr {
			stmt.ReflectValue = stmt.ReflectValue.Elem()
		}
		if !stmt.ReflectValue.IsValid() {
			db.AddError(fmt.Errorf("invalid value"))
		}
	}
	return db
}

// 从缓存中获取model全部数据, 返回json字符串
func (s *QueryRedis) get(tableName string) *gojsonq.JSONQ {
	jsonStr := ""
	if !s.Statement.json {
		// 缓存键由数据库名与表名组成
		cacheKey := fmt.Sprintf("%s_%s", global.Conf.Mysql.Database, tableName)
		var err error
		str, err := s.redis.Get(s.ctx, cacheKey).Result()
		global.Log.Debug(s.ctx, "[QueryRedis.get]读取redis缓存: %s", tableName)
		if err != nil {
			global.Log.Debug(s.ctx, "[QueryRedis.get]读取redis缓存异常: %v", err)
		} else {
			// 解压缩字符串
			jsonStr = utils.DeCompressStrByZlib(str)
		}
		if jsonStr == "" {
			// 如果是空字符串, 将其设置为空数组, 否则list会被转为nil
			jsonStr = "[]"
		}
	} else {
		jsonStr = s.Statement.jsonStr
	}
	query := s.jsonQuery(jsonStr)
	var nullList interface{}
	list := query.Get()
	if s.Statement.first {
		// 取第一条数据
		switch list.(type) {
		case []interface{}:
			v, _ := list.([]interface{})
			if len(v) > 0 {
				list = v[0]
			} else {
				// 设置为空元素
				list = nullList
			}
		}
	}
	// 类型为int64表示查询数据条数, 直接跳过结构体查询以及预加载
	if _, ok := s.Statement.Dest.(*int64); !ok {
		// 获取json结果并转换为结构体
		utils.Struct2StructByJson(list, s.Statement.Dest)
		if list != nil {
			// 处理预加载
			s.processPreload()
		} else {
			s.AddError(gorm.ErrRecordNotFound)
		}
	}

	return query
}

// jsonq需要每一次新建不同实例, 否则可能查询条件重叠, 找不到记录
func (s QueryRedis) jsonQuery(str string) *gojsonq.JSONQ {
	// 使用jsonq
	query := gojsonq.New().FromString(str)
	// 添加where条件
	for _, condition := range s.Statement.whereConditions {
		query = query.Where(condition.key, condition.cond, condition.val)
	}
	// 添加order条件
	for _, condition := range s.Statement.orderConditions {
		if condition.asc {
			query = query.SortBy(condition.property)
		} else {
			query = query.SortBy(condition.property, "desc")
		}
	}
	// 添加limit/offset
	query.Limit(s.Statement.limit)
	query.Offset(s.Statement.offset)
	return query
}
