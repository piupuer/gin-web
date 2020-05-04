package models

// 系统接口表
type SysApi struct {
	Model
	Method   string `gorm:"comment:'请求方式'" json:"method"`
	Path     string `gorm:"comment:'访问路径'" json:"path"`
	Category string `gorm:"comment:'所属类别'" json:"category"`
	Desc     string `gorm:"comment:'说明'" json:"desc"`
	Creator  string `gorm:"comment:'创建人'" json:"creator"`
}

func (m SysApi) TableName() string {
	return m.Model.TableName("sys_api")
}
