package user

import (
	"database/sql"

	"go.uber.org/zap"
	"manyface.net/internal/session"
)

type User struct {
	ID       int64    `json:"user_id,omitempty"`
	Username string   `json:"username"`
	Seed     string   `json:"seed"`
	Password string   `json:"password"`
	Mnemonic []string `json:"mnemonic"`
}

type UserHandler struct {
	Logger *zap.SugaredLogger
	Repo   *UserRepo
	SM     *session.SessionManager
}

type UserRepo struct {
	db *sql.DB
}
