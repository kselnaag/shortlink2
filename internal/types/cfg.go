package types

type CfgEnv struct {
	SL_APP_NAME    string `env:"SL_APP_NAME"`
	SL_APP_PROTOCS string `env:"SL_APP_PROTOS"`
	SL_LOG_MODE    string `env:"SL_LOG_MODE"`
	SL_LOG_LEVEL   string `env:"SL_LOG_LEVEL"`
	SL_HTTP_MODE   string `env:"SL_HTTP_MODE"`
	SL_HTTP_IP     string `env:"SL_HTTP_IP"`
	SL_HTTP_PORT   string `env:"SL_HTTP_PORT"`
	SL_DB_MODE     string `env:"SL_DB_MODE"`
	SL_DB_IP       string `env:"SL_DB_IP"`
	SL_DB_PORT     string `env:"SL_DB_PORT"`
	SL_DB_LOGIN    string `env:"SL_DB_LOGIN"`
	SL_DB_PASS     string `env:"SL_DB_PASS"`
	SL_DB_DBNAME   string `env:"SL_DB_DBNAME"`
	SL_GRPC_PORT   string `env:"SL_GRPC_PORT"`
}
