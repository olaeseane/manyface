package user

import (
	"bytes"
	"database/sql"
	"errors"

	"github.com/aidarkhanov/nanoid"
	"manyface.net/internal/utils"
)

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (repo *UserRepo) Register(password string) ([]string, string, error) {
	salt := utils.RandStringRunes(8)
	hashPassword := utils.HashIt(password, salt)
	mnemonic, err := utils.MakeMnemonic(repo.db)
	if err != nil {
		return nil, "", err
	}
	userID := nanoid.New()
	seed := utils.MakeSeed(mnemonic, salt)
	res, err := repo.db.Exec("INSERT INTO user (user_id, seed, password) VALUES (?, ?, ?)", userID, seed, hashPassword)
	if err != nil {
		return nil, "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return nil, "", err
	}
	_, err = res.LastInsertId()
	return mnemonic, userID, err
}

func (repo *UserRepo) Login(userID string, password string) error {
	var dbHashPassword []byte
	err := repo.db.QueryRow("SELECT password FROM user WHERE user_id = ?", userID).Scan(&dbHashPassword)
	if err == sql.ErrNoRows || err != nil {
		return err
	}

	salt := string(dbHashPassword[0:8])
	inHashPassword := utils.HashIt(password, salt)
	if !bytes.Equal(inHashPassword, dbHashPassword) {
		return errors.New("password mismatched")
	}

	return nil
}

/*
func (repo *UserRepo) RegisterV1beta1(username, password string) (int64, error) {
	var userID uint
	err := repo.db.QueryRow("SELECT user_id FROM userV1beta1 WHERE username = ?", username).Scan(&userID)
	if err != sql.ErrNoRows || err == nil {
		return -1, err
	}

	salt := utils.RandStringRunes(8)
	hashPassword := utils.HashIt(password, salt)
	res, err := repo.db.Exec("INSERT INTO userV1beta1 (username, password) VALUES (?, ?)", username, hashPassword)
	if err != nil {
		return -1, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return -1, err
	}
	return res.LastInsertId()
}

func (repo *UserRepo) LoginV1beta1(username, password string) (int64, error) {
	var (
		userID         uint
		dbHashPassword []byte
	)
	err := repo.db.QueryRow("SELECT user_id, password FROM userV1beta1 WHERE username = ?", username).Scan(&userID, &dbHashPassword)
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

func (repo *UserRepo) LoginV2beta1(userID int64, password string, mnemonic []string) (int64, error) {
	var seed, dbHashPassword []byte

	err := repo.db.QueryRow("SELECT seed, password FROM userV2beta1 WHERE user_id = ?", userID).Scan(&seed, &dbHashPassword)
	if err == sql.ErrNoRows || err != nil {
		return -1, err
	}

	if password != "" {
		salt := string(dbHashPassword[0:8])
		inHashPassword := utils.HashIt(password, salt)
		if !bytes.Equal(inHashPassword, dbHashPassword) {
			return -1, errors.New("password mismatched")
		}
		return int64(userID), nil
	} else {
		salt := string(seed[0:8])
		inSeed := utils.MakeSeed(mnemonic, salt)
		if !bytes.Equal(inSeed, seed) {
			return -1, errors.New("seed mismatched")
		}
		return int64(userID), nil
	}
}

*/
