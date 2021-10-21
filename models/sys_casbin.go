package models

// refer to: github.com/casbin/gorm-adapter/v2/adapter.go CasbinRule
type SysCasbin struct {
	Id    uint   `gorm:"primaryKey;autoIncrement"`
	PType string `gorm:"column:ptype;size:100;index:idx_casbin_unique;comment:'enforer type'"`
	V0    string `gorm:"size:100;index:idx_casbin_unique;comment:'role keyword(SysRole.Keyword)'"`
	V1    string `gorm:"size:100;index:idx_casbin_unique;comment:'resource name'"`
	V2    string `gorm:"size:100;index:idx_casbin_unique;comment:'request method'"`
	V3    string `gorm:"size:100;index:idx_casbin_unique"`
	V4    string `gorm:"size:100;index:idx_casbin_unique"`
	V5    string `gorm:"size:100;index:idx_casbin_unique"`
}

// role and casbin
type SysRoleCasbin struct {
	Keyword string `json:"keyword"` // role keyword
	Method  string `json:"method"`  // api method
	Path    string `json:"path"`    // api path
}
