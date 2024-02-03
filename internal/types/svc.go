package types

type ISvcShortLink2 interface {
	GetLinkPair(hash string) string
	SetLinkPair(link string) string
	DelLinkPair(hash string) bool
}
