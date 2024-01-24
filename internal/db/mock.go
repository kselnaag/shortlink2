package db

import (
	T "shortlink2/internal/types"
	"sync"
)

var _ T.IDB = (*DBMock)(nil)

type DBMock struct {
	log  T.ILog
	cfg  *T.CfgEnv
	db   map[string]string
	rwmu sync.RWMutex
}

// dbmock.Store("5clp60", "http://lib.ru")
// dbmock.Store("dhiu79", "http://google.ru")
func NewDBMock(cfg *T.CfgEnv, log T.ILog) *DBMock {
	return &DBMock{
		log: log,
		cfg: cfg,
		db:  make(map[string]string, 8),
	}
}

func (m *DBMock) SaveLinkPair(hash, link string) bool {
	if (len(hash) == 0) || (len(link) == 0) {
		return false
	}
	m.rwmu.Lock()
	m.db[hash] = link
	m.rwmu.Unlock()
	return true
}

func (m *DBMock) LoadLinkPair(hash string) string {
	m.rwmu.RLock()
	link, ok := m.db[hash]
	m.rwmu.RUnlock()
	if !ok {
		return ""
	}
	return link
}

func (m *DBMock) DeleteLinkPair(hash string) bool {
	m.rwmu.Lock()
	delete(m.db, hash)
	m.rwmu.Unlock()
	return true
}

func (m *DBMock) ConnectDB() func(e error) {
	m.log.LogInfo("mock db connected")
	return func(e error) {
		if e != nil {
			m.log.LogError(e, "DBMock.Connect(): db graceful_shutdown error")
		}
		m.log.LogInfo("mock db disconnected")
	}
}
