package global

// 本地时间格式
const (
	MsecLocalTimeFormat  = "2006-01-02 15:04:05.000"
	SecLocalTimeFormat   = "2006-01-02 15:04:05"
	DateLocalTimeFormat  = "2006-01-02"
	MonthLocalTimeFormat = "2006-01"
)

// request id
const (
	RequestIdHeader     = "X-Request-Id"
	RequestIdContextKey = "RequestId"
)

// mode
const (
	Dev   = "development"
	Stage = "staging"
	Prod  = "production"
)
