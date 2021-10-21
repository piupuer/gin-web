package models

import "github.com/piupuer/go-helper/models"

const (
	SysMachineStatusUnhealthy    uint   = 0
	SysMachineStatusHealthy      uint   = 1
	SysMachineStatusUnhealthyStr string = "unhealthy"
	SysMachineStatusHealthyStr   string = "healthy"
)

var SysMachineStatusConst = map[uint]string{
	SysMachineStatusUnhealthy: SysMachineStatusUnhealthyStr,
	SysMachineStatusHealthy:   SysMachineStatusHealthyStr,
}

type SysMachine struct {
	models.M
	Host      string `gorm:"comment:'host(IP/Domain)'" json:"host"`
	SshPort   int    `gorm:"comment:'ssh port'" json:"sshPort"`
	Version   string `gorm:"comment:'os version'" json:"version"`
	Name      string `gorm:"comment:'os name'" json:"name"`
	Arch      string `gorm:"comment:'os architecture'" json:"arch"`
	Cpu       string `gorm:"comment:'CPU model'" json:"cpu"`
	Memory    string `gorm:"comment:'memory size'" json:"memory"`
	Disk      string `gorm:"comment:'disk size'" json:"disk"`
	LoginName string `gorm:"comment:'login name'" json:"loginName"`
	LoginPwd  string `gorm:"comment:'login password'" json:"loginPwd"`
	Status    *uint  `gorm:"default:0;comment:'status(0:unhealthy 1:healthy)'" json:"status"`
	Remark    string `gorm:"comment:'remark'" json:"remark"`
}
