package db

import (
	"database/sql"
	T "shortlink2/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

var _ T.IDB = (*DBsqlite)(nil)

type DBsqlite struct {
	log    T.ILog
	cfg    T.ICfg
	dbpath string
	db     *sql.DB
}

func NewDBsqlite(cfg T.ICfg, log T.ILog, dbpath string) *DBsqlite {
	return &DBsqlite{
		log:    log,
		cfg:    cfg,
		dbpath: dbpath,
	}
}

func (s *DBsqlite) SaveLinkPair(hash, link string) bool {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(err, "DBsqlite.SaveLinkPair(): unable to ping db")
		return false
	}
	_, err1 := s.db.Exec("INSERT INTO shortlink VALUES (?, ?)", hash, link)
	if err1 != nil {
		s.log.LogError(err1, "DBsqlite.SaveLinkPair(): unable to INSERT values")
		return false
	}
	return true
}

func (s *DBsqlite) LoadLinkPair(hash string) string {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(err, "DBsqlite.LoadLinkPair(): unable to ping db")
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
		s.log.LogError(err, "DBsqlite.DeleteLinkPair(): unable to ping db")
		return false
	}
	_, err1 := s.db.Exec("DELETE FROM shortlink WHERE hash = ?", hash)
	if err1 != nil {
		s.log.LogError(err1, "DBsqlite.DeleteLinkPair(): unable to DELETE values")
		return false
	}
	return true
}

func (s *DBsqlite) InitDB() {
	if err := s.db.Ping(); err != nil {
		s.log.LogError(err, "DBsqlite.InitDB(): unable to ping db")
		return
	}
	_, err1 := s.db.Exec("CREATE TABLE IF NOT EXISTS shortlink (hash TEXT PRIMARY KEY, link TEXT NOT NULL, CHECK (link <> ''))")
	if err1 != nil {
		s.log.LogError(err1, "DBsqlite.InitDB(): unable to CREATE TABLE")
		return
	}
	_, err2 := s.db.Exec("INSERT INTO shortlink VALUES ('5clp60', 'http://lib.ru'); INSERT INTO shortlink VALUES ('dhiu79', 'http://google.ru');")
	if err2 != nil {
		s.log.LogError(err2, "DBsqlite.InitDB(): unable to INSERT values")
		return
	}
}

func (s *DBsqlite) ConnectDB() func(e error) {
	db, err := sql.Open("sqlite3", s.dbpath)
	if err != nil {
		s.log.LogError(err, "DBsqlite.ConnectDB(): unable to connect")
		return func(e error) {}
	}
	s.db = db
	s.InitDB()
	s.log.LogInfo("DBsqlite connected")
	return func(e error) {
		if err := s.db.Close(); err != nil {
			s.log.LogError(err, "DBsqlite.ConnectDB(): db graceful_shutdown error")
		}
		if e != nil {
			s.log.LogError(e, "DBsqlite.ConnectDB(): db graceful_shutdown with error")
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
