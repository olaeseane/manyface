package session

import (
	"context"
	"database/sql"
	"errors"

	"manyface.net/internal/utils"
)

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		DB: db,
	}
}

func (sm *SessionManager) Create(userID int64) (string, error) {
	sessID := utils.RandStringRunes(32)
	res, err := sm.DB.Exec("INSERT INTO session (sess_id, user_id) VALUES (?, ?)", sessID, userID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return "", err
	}
	return sessID, nil
}

func (sm *SessionManager) Check(sessID string) (*Session, error) {
	sess := &Session{ID: sessID}
	if err := sm.DB.QueryRow("SELECT user_id FROM session WHERE sess_id = ?", sess.ID).Scan(&sess.UserID); err != nil {
		return nil, err
	}
	return sess, nil
}

func (sm *SessionManager) DeleteCurrent(sessID string) error {
	res, err := sm.DB.Exec("DELETE FROM session WHERE sess_id = ?", sessID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) DeleteAll(userID uint) error {
	res, err := sm.DB.Exec("DELETE FROM session WHERE user_id = ?", userID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) GetFromCtx(cxt context.Context) (*Session, error) {
	sess, ok := cxt.Value(SessKey).(*Session)
	if !ok {
		return nil, errors.New("Session not found")
	}
	return sess, nil
}
