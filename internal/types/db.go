package types

type IDB interface {
	SaveLinkPair(hash, link string) bool
	LoadLinkPair(hash string) string
	DeleteLinkPair(hash string) bool
	ConnectDB() func(e error)
}
