package models

import (
	"gin-web/pkg/utils"
)

const (
	// 机器状态
	SysMachineStatusUnConnected    uint   = 0      // 无法连接
	SysMachineStatusNormal         uint   = 1      // 正常
	SysMachineStatusUnConnectedStr string = "无法连接" // 无法连接
	SysMachineStatusNormalStr      string = "正常"   // 正常
)

// 定义map方便取值
var SysMachineStatusConst = map[uint]string{
	SysMachineStatusUnConnected: SysMachineStatusUnConnectedStr,
	SysMachineStatusNormal:      SysMachineStatusNormalStr,
}

// 机器配置
type SysMachine struct {
	Model
	Host      string `gorm:"comment:'主机地址(可以是IP或域名)'" json:"host"`
	SshPort   int    `gorm:"comment:'ssh端口号'" json:"sshPort"`
	Version   string `gorm:"comment:'操作系统版本'" json:"version"`
	Name      string `gorm:"comment:'系统名字'" json:"name"`
	Arch      string `gorm:"comment:'系统架构'" json:"arch"`
	Cpu       string `gorm:"comment:'CPU型号'" json:"cpu"`
	Memory    string `gorm:"comment:'内存'" json:"memory"`
	Disk      string `gorm:"comment:'硬盘'" json:"disk"`
	LoginName string `gorm:"comment:'登陆名'" json:"loginName"`
	LoginPwd  string `gorm:"comment:'登陆密码'" json:"loginPwd"`
	Status    *uint  `gorm:"default:0;comment:'状态(0:无法连接 1:正常)'" json:"status"`
	Remark    string `gorm:"comment:'备注'" json:"remark"`
	Creator   string `gorm:"comment:'创建人'" json:"creator"`
}

func (m *SysMachine) TableName() string {
	return m.Model.TableName("sys_machine")
}

// 获取ssh配置项
func (m *SysMachine) GetSshConfig(timeout int) utils.SshConfig {
	return utils.SshConfig{
		LoginName: m.LoginName,
		LoginPwd:  m.LoginPwd,
		Port:      m.SshPort,
		Host:      m.Host,
		Timeout:   timeout,
	}
}
