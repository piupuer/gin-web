package global

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

const (
	ProName = "gin-web"
	ProEnvName = "GIN_WEB"
)
