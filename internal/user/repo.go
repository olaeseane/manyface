package user

import (
	"bytes"
	"database/sql"
	"errors"

	"manyface.net/internal/utils"
)

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo *UserRepo) Register(username, password string) (int64, error) {
	var userID uint
	err := repo.db.QueryRow("SELECT user_id FROM user WHERE username = ?", username).Scan(&userID)
	if err != sql.ErrNoRows || err == nil {
		return -1, err
	}

	salt := utils.RandStringRunes(8)
	hashPassword := utils.HashIt(password, salt)
	res, err := repo.db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, hashPassword)
	if err != nil {
		return -1, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return -1, err
	}
	return res.LastInsertId()
}

func (repo *UserRepo) Login(username, password string) (int64, error) {
	var (
		userID         uint
		dbHashPassword []byte
	)
	err := repo.db.QueryRow("SELECT user_id, password FROM user WHERE username = ?", username).Scan(&userID, &dbHashPassword)
	if err == sql.ErrNoRows || err != nil {
		return -1, err
	}

	salt := string(dbHashPassword[0:8])
	inHashPassword := utils.HashIt(password, salt)
	if !bytes.Equal(inHashPassword, dbHashPassword) {
		return -1, errors.New("password mismatched")
	}

	return int64(userID), nil
}
