package types

type IHTTPServer interface {
	Run() func(e error)
}

type HTTPMess struct {
	Method string `json:"M"`
	Hash   string `json:"H"`
	Link   string `json:"L"`
}
