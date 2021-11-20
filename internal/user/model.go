package user

import (
	"database/sql"

	"go.uber.org/zap"
	"manyface.net/internal/session"
)

type User struct {
	ID       int64  `json:"user_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	Logger *zap.SugaredLogger
	Repo   *UserRepo
	SM     *session.SessionManager
}

type UserRepo struct {
	db *sql.DB
}
