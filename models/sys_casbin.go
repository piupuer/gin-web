package models

// Casbin权限访问控制表, 参见github.com/casbin/gorm-adapter/v2/adapter.go CasbinRule
// 可以根据项目实际需要动态设定, 这里用到了3个字段 角色关键字/资源名称/请求类型
type SysCasbin struct {
	PType       string `gorm:"size:100;comment:'策略类型'"`
	V0          string `gorm:"size:100;comment:'角色关键字'"`
	V1          string `gorm:"size:100;comment:'资源名称'"`
	V2          string `gorm:"size:100;comment:'请求类型'"`
	V3          string `gorm:"size:100"`
	V4          string `gorm:"size:100"`
	V5          string `gorm:"size:100"`
}

func (m SysCasbin) TableName() string {
	// 这里与Casbin官方统一表名称
	return "casbin_rule"
}
