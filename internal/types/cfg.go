package types

type ICfg interface {
	GetVal(string) string
	Parse() ICfg
}

const (
	SL_APP_NAME    = "SL_APP_NAME"
	SL_APP_PROTOCS = "SL_APP_PROTOCS"
	SL_LOG_LEVEL   = "SL_LOG_LEVEL"
	SL_HTTP_IP     = "SL_HTTP_IP"
	SL_HTTP_PORT   = "SL_HTTP_PORT"
	SL_GRPC_PORT   = "SL_GRPC_PORT"
)
