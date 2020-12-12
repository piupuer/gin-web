package service

import (
	"errors"
	"gin-web/pkg/global"
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

// 创建, model需使用指针, 否则可能无法插入数据
func (s *MysqlService) Create(req interface{}, model interface{}) (err error) {
	utils.Struct2StructByJson(req, model)
	// 创建数据
	err = s.tx.Create(model).Error
	return
}

// 根据编号更新
func (s *MysqlService) UpdateById(id uint, req interface{}) error {
	// 通过反射获取请求类型
	reqType := reflect.TypeOf(req)
	oldModel := reflect.New(reqType)
	oldModelIns := oldModel.Interface()
	query := s.tx.Model(oldModelIns).Where("id = ?", id).First(oldModelIns)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在, 更新失败")
	}

	// 比对增量字段
	newModel := reflect.New(reqType)
	newModelIns := newModel.Interface()
	utils.CompareDifferenceStructByJson(oldModelIns, req, newModelIns)

	// 更新指定列
	return query.Updates(newModelIns).Error
}

// 批量删除, model需使用指针, 否则可能无法插入数据
func (s *MysqlService) DeleteByIds(ids []uint, model interface{}) (err error) {
	return s.tx.Where("id IN (?)", ids).Delete(model).Error
}
