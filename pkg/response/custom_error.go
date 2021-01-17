package response

// 自定义错误码与错误信息

const (
	Ok                  = 201
	NotOk               = 405
	Unauthorized        = 401
	Forbidden           = 403
	InternalServerError = 500
)

const (
	OkMsg                      = "操作成功"
	NotOkMsg                   = "操作失败"
	UnauthorizedMsg            = "登录过期, 需要重新登录"
	LoginCheckErrorMsg         = "用户名或密码错误"
	ForbiddenMsg               = "无权访问该资源, 请联系网站管理员授权"
	InternalServerErrorMsg     = "服务器内部错误"
	IdempotenceTokenEmptyMsg   = "幂等性token为空"
	IdempotenceTokenInvalidMsg = "幂等性token失效, 重复提交"
)

var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
}
