package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	T "shortlink2/internal/types"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var _ T.IDB = (*DBsqlite)(nil)

type DBsqlite struct {
	log    T.ILog
	cfg    T.ICfg
	dbpath string
	db     *sql.DB
}

func NewDBsqlite(cfg T.ICfg, log T.ILog, dir string) *DBsqlite {
	return &DBsqlite{
		log:    log,
		cfg:    cfg,
		dbpath: filepath.Join(dir, "db/sqlite.db"),
	}
}

func (s *DBsqlite) SaveLinkPair(hash, link string) bool {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.SaveLinkPair(): unable to ping db", err))
		return false
	}
	_, err1 := s.db.Exec("INSERT INTO shortlink VALUES (?, ?)", hash, link)
	if err1 != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.SaveLinkPair(): unable to INSERT values", err1))
		return false
	}
	return true
}

func (s *DBsqlite) LoadLinkPair(hash string) string {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.LoadLinkPair(): unable to ping db", err))
		return ""
	}
	var pair T.DBMess
	err1 := s.db.QueryRow("SELECT hash, link FROM shortlink WHERE hash = ?", hash).Scan(&(pair.Hash), &(pair.Link))
	if err1 != nil {
		s.log.LogDebug("DBsqlite.LoadLinkPair(): unable to SELECT values")
		return ""
	}
	return pair.Link
}

func (s *DBsqlite) DeleteLinkPair(hash string) bool {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.DeleteLinkPair(): unable to ping db", err))
		return false
	}
	_, err1 := s.db.Exec("DELETE FROM shortlink WHERE hash = ?", hash)
	if err1 != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.DeleteLinkPair(): unable to DELETE values", err1))
		return false
	}
	return true
}

func (s *DBsqlite) InitDB() {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.InitDB(): unable to ping db", err))
		return
	}
	_, err1 := s.db.Exec("CREATE TABLE IF NOT EXISTS shortlink (hash TEXT PRIMARY KEY, link TEXT NOT NULL, CHECK (link <> ''))")
	if err1 != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.InitDB(): unable to CREATE TABLE", err1))
		return
	}
	_, err2 := s.db.Exec("INSERT INTO shortlink VALUES ('5clp60', 'http://lib.ru'); INSERT INTO shortlink VALUES ('dhiu79', 'http://google.ru');")
	if err2 != nil && !strings.Contains(err2.Error(), "UNIQUE constraint failed") {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.InitDB(): unable to INSERT values", err2))
		return
	}
}

func (s *DBsqlite) ConnectDB() func(e error) {
	db, err := sql.Open("sqlite3", s.dbpath)
	if err != nil {
		s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.ConnectDB(): unable to connect", err))
		return func(e error) {}
	}
	s.db = db
	s.InitDB()
	s.log.LogInfo("DBsqlite connected")
	return func(e error) {
		if err := s.db.Close(); err != nil {
			s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.ConnectDB(): db graceful_shutdown error", err))
		}
		if e != nil {
			s.log.LogError(fmt.Errorf("%s: %w", "DBsqlite.ConnectDB(): db graceful_shutdown with error", e))
		}
		s.db = nil
		s.log.LogInfo("DBsqlite disconnected")
	}
}

/* rows, err1 := s.db.Query("SELECT hash, link FROM shortlink WHERE hash = ?", hash)
for rows.Next() {
	var pair T.DBMess
	err2 := rows.Scan(&(pair.Hash), &(pair.Link))
}
err3 := rows.Err()
} */
