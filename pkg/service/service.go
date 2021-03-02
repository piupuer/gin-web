package service

import (
	"errors"
	"fmt"
	"gin-web/pkg/global"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
)

type MysqlService struct {
	tx *gorm.DB // 事务对象实例
	db *gorm.DB // 无事务对象实例
}

// 初始化服务
func New(c *gin.Context) MysqlService {
	// 获取事务对象
	tx := global.GetTx(c)
	return MysqlService{
		tx: tx,
		db: global.Mysql,
	}
}

// 查询, model需使用指针, 否则可能无法绑定数据
func (s *MysqlService) Find(query *gorm.DB, page *response.PageInfo, model interface{}) (err error) {
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
			err = query.Limit(limit).Offset(offset).Find(model).Error
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
	return
}

// Scan查询, 适用于多表联合查询, model需使用指针, 否则可能无法绑定数据
func (s *MysqlService) Scan(query *gorm.DB, page *response.PageInfo, model interface{}) (err error) {
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
func (s *MysqlService) Create(req interface{}, model interface{}) (err error) {
	utils.Struct2StructByJson(req, model)
	// 创建数据
	err = s.tx.Create(model).Error
	return
}

// 根据编号更新
func (s *MysqlService) UpdateById(id uint, req interface{}, model interface{}) error {
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
func (s *MysqlService) DeleteByIds(ids []uint, model interface{}) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(model).Error
}
