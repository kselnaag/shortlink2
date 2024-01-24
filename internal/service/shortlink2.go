package svc

import (
	"hash/crc32"
	T "shortlink2/internal/types"
	"strconv"
)

var _ T.ISvcShortLink2 = (*SvcShortLink2)(nil)

type SvcShortLink2 struct {
	db  T.IDB
	log T.ILog
}

func NewSvcShortLink2(db T.IDB, log T.ILog) *SvcShortLink2 {
	return &SvcShortLink2{
		db:  db,
		log: log,
	}
}

func (s *SvcShortLink2) GetLinkLongFromLinkShort(hash string) string {
	return s.db.LoadLinkPair(hash)
}

func (s *SvcShortLink2) SetLinkPair(link string) bool {
	hash := calcLinkShort(link)
	return s.db.SaveLinkPair(hash, link)
}

func (s *SvcShortLink2) DelLinkPair(hash string) bool {
	return s.db.DeleteLinkPair(hash)
}

func calcLinkShort(link string) string {
	hashlen := 6
	radixlen := 36
	hash := crc32.ChecksumIEEE([]byte(link))
	str := strconv.FormatUint(uint64(hash), radixlen)
	if len(str) > hashlen {
		idx := len(str) - hashlen
		return str[idx:]
	}
	return str
}
