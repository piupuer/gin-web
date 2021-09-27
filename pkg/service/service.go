package service

import (
	"context"
	"errors"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

type MysqlService struct {
	ctx *gin.Context // 上下文
	tx  *gorm.DB     // 事务对象实例
	db  *gorm.DB     // 无事务对象实例
}

// 初始化服务
func New(c *gin.Context) MysqlService {
	nc := gin.Context{}
	if c != nil {
		nc = *c
	}
	s := MysqlService{
		ctx: &nc,
	}
	tx := global.GetTx(&nc)
	rc := s.RequestIdContext("")
	s.tx = tx.WithContext(rc)
	s.db = global.Mysql.WithContext(rc)
	return s
}

var (
	findByIdsCache = cache.New(24*time.Hour, 48*time.Hour)
	findCountCache = cache.New(5*time.Minute, 48*time.Hour)
)

// 获取携带request id的上下文
func (s MysqlService) RequestIdContext(requestId string) context.Context {
	if s.ctx != nil && requestId == "" {
		requestId = s.ctx.GetString(global.RequestIdContextKey)
	}
	return global.RequestIdContext(requestId)
}

// 查询指定id, model需使用指针, 否则可能无法绑定数据
func (s MysqlService) FindById(id uint, model interface{}, setCache bool) (err error) {
	return s.FindByKeys("id", id, model, setCache)
}

// 查询指定id列表, model需使用指针, 否则可能无法绑定数据
func (s MysqlService) FindByIds(ids []uint, model interface{}, setCache bool) (err error) {
	return s.FindByKeys("id", ids, model, setCache)
}

// 查询指定key列表, model需使用指针, 否则可能无法绑定数据(如不使用cache可设置为false)
func (s MysqlService) FindByKeys(key string, ids interface{}, model interface{}, setCache bool) (err error) {
	return s.FindByKeysWithPreload(key, nil, ids, model, setCache)
}

// 查询指定key列表, 并且preload其他表, model需使用指针, 否则可能无法绑定数据(如不使用cache可设置为false)
func (s MysqlService) FindByKeysWithPreload(key string, preloads []string, ids interface{}, model interface{}, setCache bool) (err error) {
	var newIds interface{}
	var firstId interface{}
	// 判断ids是否数组
	idsRv := reflect.ValueOf(ids)
	idsRt := reflect.TypeOf(ids)
	newIdsRv := reflect.ValueOf(newIds)
	newIdsIsArr := false
	if idsRv.Kind() == reflect.Ptr {
		return fmt.Errorf("ids不能为指针类型")
	}
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("model必须是非空指针类型")
	}
	if key == "" {
		key = "id"
	}
	// 处理参数是否与值类似匹配
	if idsRv.Kind() == reflect.Slice {
		if idsRv.Len() == 0 {
			return
		}
		// 获取第一个元素
		firstId = idsRv.Index(0).Convert(idsRt.Elem()).Interface()
		if idsRv.Len() > 1 {
			// 参数是数组, 值不是数组
			if reflect.ValueOf(model).Elem().Kind() != reflect.Slice {
				newIds = firstId
			} else {
				// 创建新数组, 避免数组引用传递在外层被修改
				newArr := reflect.MakeSlice(reflect.TypeOf(ids), idsRv.Len(), idsRv.Len())
				// 通过golang提供的拷贝方法
				reflect.Copy(newArr, idsRv)
				newIds = newArr.Interface()
				newIdsIsArr = true
			}
		} else {
			// len=0, 将值重写
			newIds = firstId
		}
	} else {
		firstId = ids
		newIds = ids
	}
	// 刷新反射值
	newIdsRv = reflect.ValueOf(newIds)

	// 可能一个条件有多个查询结果
	if key != "id" && !newIdsIsArr && newIdsRv.Kind() != reflect.Slice && rv.Elem().Kind() == reflect.Slice {
		newIdsIsArr = true
	}
	// ids是数组, 但model却不是数组
	if newIdsIsArr && rv.Elem().Kind() != reflect.Slice {
		// ids取第一个值
		newIds = firstId
	}
	cacheKey := ""
	// 需要设置缓存
	if setCache {
		structName := ""
		if rv.Elem().Kind() == reflect.Slice {
			structName = strings.ToLower(rv.Elem().Type().Elem().String())
		} else {
			structName = strings.ToLower(rv.Elem().Type().String())
		}
		preload := "preload_nothing"
		if len(preloads) > 0 {
			preload = "preload_" + strings.ToLower(strings.Join(preloads, "_"))
		}
		// 缓存key组成: table+preloads+key+ids+modelIsArr
		cacheKey = fmt.Sprintf("%s_%s_%s_%s_find", structName, preload, key, utils.Struct2Json(newIds))
		if rv.Elem().Kind() != reflect.Slice {
			cacheKey = fmt.Sprintf("%s_%s_%s_%s_first", structName, preload, key, utils.Struct2Json(newIds))
		}
		oldCache, ok := findByIdsCache.Get(cacheKey)
		if ok {
			// 通过反射回写数据, 而不是直接赋值
			// model = oldCache
			crv := reflect.ValueOf(oldCache)
			if rv.Elem().Kind() == reflect.Struct && crv.Kind() == reflect.Slice {
				rv.Elem().Set(crv.Index(0))
			} else if rv.Elem().Kind() == reflect.Slice && crv.Kind() == reflect.Struct {
				// 结构体写入数组第一个元素
				newArr1 := reflect.MakeSlice(rv.Elem().Type(), 1, 1)
				v := newArr1.Index(0)
				v.Set(crv)
				// 创建新数组, 避免数组引用传递在外层被修改
				newArr2 := reflect.MakeSlice(rv.Elem().Type(), 1, 1)
				reflect.Copy(newArr2, newArr1)
				rv.Elem().Set(newArr2)
			} else if rv.Elem().Kind() == reflect.Slice && crv.Kind() == reflect.Slice {
				// 创建新数组, 避免数组引用传递在外层被修改
				newArr := reflect.MakeSlice(rv.Elem().Type(), crv.Len(), crv.Len())
				// 通过golang提供的拷贝方法
				reflect.Copy(newArr, crv)
				rv.Elem().Set(newArr)
			} else {
				rv.Elem().Set(crv)
			}
			return
		}
	}
	query := s.tx
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	if !newIdsIsArr {
		err = query.
			Where(fmt.Sprintf("`%s` = ?", key), newIds).
			First(model).Error
	} else {
		if newIdsIsArr && newIdsRv.Kind() != reflect.Slice {
			// 可能一个条件有多个查询结果
			err = query.
				Where(fmt.Sprintf("`%s` = ?", key), firstId).
				Find(model).Error
		} else {
			err = query.
				Where(fmt.Sprintf("`%s` IN (?)", key), newIds).
				Find(model).Error
		}
	}
	if setCache {
		if rv.Elem().Kind() == reflect.Slice {
			// 如果model是数组, 需创建新数组, 避免数组引用传递在外层被修改
			newArr := reflect.MakeSlice(rv.Elem().Type(), rv.Elem().Len(), rv.Elem().Len())
			reflect.Copy(newArr, rv.Elem())
			findByIdsCache.Set(cacheKey, newArr.Interface(), cache.DefaultExpiration)
		} else {
			// 写入缓存
			findByIdsCache.Set(cacheKey, rv.Elem().Interface(), cache.DefaultExpiration)
		}
	}
	return
}

// 查询, model需使用指针, 否则可能无法绑定数据
func (s MysqlService) Find(query *gorm.DB, page *response.PageInfo, model interface{}) (err error) {
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Slice) {
		return fmt.Errorf("model必须是非空指针数组类型")
	}

	countCache := false
	if page.CountCache != nil {
		countCache = *page.CountCache
	}
	if !page.NoPagination {
		if !page.SkipCount {
			// 查询条数
			fromCache := false
			// 以sql语句作为缓存键
			stmt := query.Session(&gorm.Session{DryRun: true}).Count(&page.Total).Statement
			cacheKey := s.tx.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
			if countCache {
				countCache, ok := findCountCache.Get(cacheKey)
				if ok {
					total, _ := countCache.(int64)
					page.Total = total
					fromCache = true
				}
			}
			if !fromCache {
				err = query.Count(&page.Total).Error
				if err == nil {
					findCountCache.Set(cacheKey, page.Total, cache.DefaultExpiration)
				}
			} else {
				global.Log.Debug(s.ctx, "条数缓存命中: %s, total: %d", cacheKey, page.Total)
			}
		}
		if page.Total > 0 || page.SkipCount {
			// 获取分页参数
			limit, offset := page.GetLimit()
			if page.LimitPrimary == "" {
				err = query.Limit(limit).Offset(offset).Find(model).Error
			} else {
				// 解析model
				if query.Statement.Model != nil {
					err = query.Statement.Parse(query.Statement.Model)
					if err != nil {
						return
					}
				}
				err = query.Joins(
					// 通过索引先分页再获取join其他字段, 以提高查询效率
					fmt.Sprintf(
						"JOIN (?) AS `OFFSET_T` ON `%s`.`id` = `OFFSET_T`.`%s`",
						query.Statement.Table,
						page.LimitPrimary,
					),
					query.
						Session(&gorm.Session{}).
						Select(
							fmt.Sprintf("`%s`.`%s`", query.Statement.Table, page.LimitPrimary),
						).
						Limit(limit).
						Offset(offset),
				).Find(model).Error
			}
		}
	} else {
		// 不使用分页
		err = query.Find(model).Error
		if err == nil {
			page.Total = int64(rv.Elem().Len())
			// 获取分页参数
			page.GetLimit()
		}
	}
	page.CountCache = &countCache
	return
}

// Scan查询, 适用于多表联合查询, model需使用指针, 否则可能无法绑定数据
func (s MysqlService) Scan(query *gorm.DB, page *response.PageInfo, model interface{}) (err error) {
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Slice) {
		return fmt.Errorf("model必须是非空指针数组类型")
	}

	if !page.NoPagination {
		// 查询条数
		err = query.Count(&page.Total).Error
		if err == nil && page.Total > 0 {
			// 获取分页参数
			limit, offset := page.GetLimit()
			err = query.Limit(limit).Offset(offset).Scan(model).Error
		}
	} else {
		// 不使用分页
		err = query.Scan(model).Error
		if err == nil {
			page.Total = int64(rv.Elem().Len())
			// 获取分页参数
			page.GetLimit()
		}
	}
	return
}

// 创建, model需使用指针, 否则可能无法插入数据
func (s MysqlService) Create(req interface{}, model interface{}) (err error) {
	utils.Struct2StructByJson(req, model)
	// 创建数据
	err = s.tx.Create(model).Error
	return
}

// 根据编号更新
func (s MysqlService) UpdateById(id uint, req interface{}, model interface{}) error {
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Struct) {
		return fmt.Errorf("model必须是非空指针结构体类型")
	}
	query := s.tx.Model(rv.Interface()).Where("id = ?", id).First(rv.Interface())
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在, 更新失败")
	}

	// 比对增量字段
	m := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(rv.Elem().Interface(), req, &m)

	// 更新指定列
	return query.Updates(&m).Error
}

// 批量删除, model需使用指针, 否则可能无法插入数据
func (s MysqlService) DeleteByIds(ids []uint, model interface{}) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(model).Error
}
