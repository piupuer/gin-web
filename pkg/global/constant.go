package global

// request id
const (
	RequestIdHeader     = "X-Request-Id"
	RequestIdContextKey = "RequestId"
	TxCtxKey = "tx"
)

// mode
const (
	Dev   = "development"
	Stage = "staging"
	Prod  = "production"
)

const (
	ProName    = "gin-web"
	ProEnvName = "GIN_WEB"
)

// fsm categories
const (
	FsmCategoryLeave uint = iota + 1
)
