package types

type IDB interface {
	LoadLinkPair(hash string) (string, string)
	SaveLinkPair(hash, link string) bool
	DeleteLinkPair() bool
	ConnectDB() func(e error)
}
