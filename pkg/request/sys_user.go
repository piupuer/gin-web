package request

// User login structure
type RegisterAndLoginStruct struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// 修改密码结构体
type ChangePwdStruct struct {
	OldPassword  string `json:"old_password"`
	NewPassword  string `json:"new_password"`
}

