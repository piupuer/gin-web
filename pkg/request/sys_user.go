package request

// User login structure
type RegisterAndLoginRequestStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 修改密码结构体
type ChangePwdRequestStruct struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
