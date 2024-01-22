package types

type IHTTPServer interface {
	Run() func(e error)
}
