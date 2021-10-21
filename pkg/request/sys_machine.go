package request

import (
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type MachineReq struct {
	Id                uint   `json:"id" form:"id"`
	Host              string `json:"host" form:"host"`
	SshPort           int    `json:"sshPort" form:"sshPort"`
	Version           string `json:"version" form:"version"`
	Name              string `json:"name" form:"name"`
	Arch              string `json:"arch" form:"arch"`
	Cpu               string `json:"cpu" form:"cpu"`
	Memory            string `json:"memory" form:"memory"`
	Disk              string `json:"disk" form:"disk"`
	LoginName         string `json:"loginName" form:"loginName"`
	LoginPwd          string `json:"loginPwd" form:"loginPwd"`
	Status            *uint  `json:"status" form:"status"`
	Remark            string `json:"remark" form:"remark"`
	resp.Page
}

type CreateMachineReq struct {
	Host      string  `json:"host" validate:"required"`
	SshPort   req.NullUint `json:"sshPort" validate:"required"`
	Version   string  `json:"version"`
	Name      string  `json:"name"`
	Arch      string  `json:"arch"`
	Cpu       string  `json:"cpu"`
	Memory    string  `json:"memory"`
	Disk      string  `json:"disk"`
	LoginName string  `json:"loginName" validate:"required"`
	LoginPwd  string  `json:"loginPwd" validate:"required"`
	Status    req.NullUint `json:"status"`
	Remark    string  `json:"remark"`
}

type MachineShellWsReq struct {
	Host      string  `json:"host" form:"host"`
	SshPort   req.NullUint `json:"sshPort" form:"sshPort"`
	LoginName string  `json:"loginName" form:"loginName"`
	LoginPwd  string  `json:"loginPwd" form:"loginPwd"`
	InitCmd   string  `json:"initCmd" form:"initCmd"`
	Cols      req.NullUint `json:"cols" form:"cols"`
	Rows      req.NullUint `json:"rows" form:"rows"`
}

func (s CreateMachineReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Host"] = "主机地址"
	m["SshPort"] = "ssh端口"
	m["LoginName"] = "登录名"
	m["LoginPwd"] = "登录密码"
	return m
}

type UpdateMachineReq struct {
	Host      *string  `json:"host"`
	SshPort   *req.NullUint `json:"sshPort"`
	Version   *string  `json:"version"`
	Name      *string  `json:"name"`
	Arch      *string  `json:"arch"`
	Cpu       *string  `json:"cpu"`
	Memory    *string  `json:"memory"`
	Disk      *string  `json:"disk"`
	LoginName *string  `json:"loginName"`
	LoginPwd  *string  `json:"loginPwd"`
	Status    *req.NullUint `json:"status"`
	Remark    *string  `json:"remark"`
}
