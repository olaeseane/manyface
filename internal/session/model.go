package session

import (
	"database/sql"
)

// type ISessionManager interface {
// 	Check(*http.Request) (*Session, error)
// 	Create(http.ResponseWriter, *User) error
// 	DestroyCurrent(http.ResponseWriter, *http.Request) error
// 	DestroyAll(http.ResponseWriter, *User) error
// }

type SessionCtxType string

const SessKey SessionCtxType = "sess"

type SessionManager struct {
	DB *sql.DB
}

type Session struct {
	ID     string
	UserID int64
}
