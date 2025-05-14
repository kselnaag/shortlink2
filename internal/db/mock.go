package db

import (
	"fmt"
	T "shortlink2/internal/types"
	"sync"
)

var _ T.IDB = (*DBmock)(nil)

type DBmock struct {
	log  T.ILog
	cfg  T.ICfg
	db   map[string]string
	rwmu sync.RWMutex
}

func NewDBmock(cfg T.ICfg, log T.ILog) *DBmock {
	mockdb := make(map[string]string, 8)
	mockdb["5clp60"] = "http://lib.ru"
	// dbmock.Store("5clp60", "http://lib.ru")
	// dbmock.Store("dhiu79", "http://google.ru")
	return &DBmock{
		log: log,
		cfg: cfg,
		db:  mockdb,
	}
}

func (m *DBmock) SaveLinkPair(hash, link string) bool {
	if (len(hash) == 0) || (len(link) == 0) {
		return false
	}
	m.rwmu.Lock()
	m.db[hash] = link
	m.rwmu.Unlock()
	return true
}

func (m *DBmock) LoadLinkPair(hash string) string {
	m.rwmu.RLock()
	link, ok := m.db[hash]
	m.rwmu.RUnlock()
	if !ok {
		return ""
	}
	return link
}

func (m *DBmock) DeleteLinkPair(hash string) bool {
	m.rwmu.Lock()
	delete(m.db, hash)
	m.rwmu.Unlock()
	return true
}

func (m *DBmock) ConnectDB() func(e error) {
	m.log.LogInfo("mock db connected")
	return func(e error) {
		if e != nil {
			// m.log.LogError(e, "DBmock.Connect(): db graceful_shutdown with error")
			m.log.LogError(fmt.Errorf("%s: %w", "DBmock.Connect(): db graceful_shutdown with error", e))
		}
		m.log.LogInfo("mock db disconnected")
	}
}
