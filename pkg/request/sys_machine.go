package request

import "gin-web/pkg/response"

// 获取机器列表结构体
type MachineListRequestStruct struct {
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
	Creator           string `json:"creator" form:"creator"`
	response.PageInfo        // 分页参数
}

// 创建机器结构体
type CreateMachineRequestStruct struct {
	Host      string  `json:"host" validate:"required"`
	SshPort   ReqUint `json:"sshPort" validate:"required"`
	Version   string  `json:"version"`
	Name      string  `json:"name"`
	Arch      string  `json:"arch"`
	Cpu       string  `json:"cpu"`
	Memory    string  `json:"memory"`
	Disk      string  `json:"disk"`
	LoginName string  `json:"loginName" validate:"required"`
	LoginPwd  string  `json:"loginPwd" validate:"required"`
	Status    ReqUint `json:"status"`
	Remark    string  `json:"remark"`
	Creator   string  `json:"creator"`
}

// 机器shell ws请求结构体
type MachineShellWsRequestStruct struct {
	Host      string  `json:"host" form:"host"`
	SshPort   ReqUint `json:"sshPort" form:"sshPort"`
	LoginName string  `json:"loginName" form:"loginName"`
	LoginPwd  string  `json:"loginPwd" form:"loginPwd"`
	InitCmd   string  `json:"initCmd" form:"initCmd"`
	Cols      ReqUint `json:"cols" form:"cols"`
	Rows      ReqUint `json:"rows" form:"rows"`
}

// 翻译需要校验的字段名称
func (s CreateMachineRequestStruct) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Host"] = "主机地址"
	m["SshPort"] = "ssh端口"
	m["LoginName"] = "登录名"
	m["LoginPwd"] = "登录密码"
	return m
}

// 更新机器结构体
type UpdateMachineRequestStruct struct {
	Host      *string  `json:"host"`
	SshPort   *ReqUint `json:"sshPort"`
	Version   *string  `json:"version"`
	Name      *string  `json:"name"`
	Arch      *string  `json:"arch"`
	Cpu       *string  `json:"cpu"`
	Memory    *string  `json:"memory"`
	Disk      *string  `json:"disk"`
	LoginName *string  `json:"loginName"`
	LoginPwd  *string  `json:"loginPwd"`
	Status    *ReqUint `json:"status"`
	Remark    *string  `json:"remark"`
	Creator   *string  `json:"creator"`
}
