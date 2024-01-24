package types

type ISvcShortLink2 interface {
	GetLinkLongFromLinkShort(hash string) string
	SetLinkPair(link string) bool
	DelLinkPair(hash string) bool
}
