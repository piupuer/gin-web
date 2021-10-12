package models

// Casbin权限访问控制表, 参见github.com/casbin/gorm-adapter/v2/adapter.go CasbinRule
// 可以根据项目实际需要动态设定, 这里用到了3个字段 角色关键字/资源名称/请求类型
type SysCasbin struct {
	Id    uint   `gorm:"primaryKey;autoIncrement"`
	PType string `gorm:"column:ptype;size:100;index:idx_casbin_unique;comment:'策略类型'"`
	V0    string `gorm:"size:100;index:idx_casbin_unique;comment:'角色关键字'"`
	V1    string `gorm:"size:100;index:idx_casbin_unique;comment:'资源名称'"`
	V2    string `gorm:"size:100;index:idx_casbin_unique;comment:'请求类型'"`
	V3    string `gorm:"size:100;index:idx_casbin_unique"`
	V4    string `gorm:"size:100;index:idx_casbin_unique"`
	V5    string `gorm:"size:100;index:idx_casbin_unique"`
}

// 角色权限规则
type SysRoleCasbin struct {
	Keyword string `json:"keyword"` // 角色关键字
	Method  string `json:"method"`  // 请求方式
	Path    string `json:"path"`    // 访问路径
}
