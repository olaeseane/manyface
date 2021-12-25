package messenger

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
)

func (srv *MsgServer) CreateFaceV2beta1(nick, purpose, bio, comments string, userID int64) (string, error) {
	hash := md5.Sum([]byte(nick + "|" + purpose + "|" + bio + strconv.FormatInt(userID, 10)))
	face_id := hex.EncodeToString(hash[:])
	res, err := srv.db.Exec("INSERT INTO faceV2beta1 (face_id, nick, purpose, bio, comments, user_id) VALUES (?, ?, ?, ?, ?, ?)", face_id, nick, purpose, bio, comments, userID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return "", err
	}
	return face_id, nil
}

func (srv *MsgServer) GetFaceByIDV2beta1(faceID string, userID int64) (*Face, error) {
	f := &Face{}
	err := srv.db.QueryRow("SELECT face_id, nick, purpose, bio, comments, user_id FROM faceV2beta1 WHERE face_id = ? AND user_id = ?", faceID, userID).
		Scan(&f.ID, &f.Nick, &f.Purpose, &f.Bio, &f.Comments, &f.UserID)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return f, nil
}

func (srv *MsgServer) GetFacesByUserV2beta1(userID int64) ([]*Face, error) {
	rows, err := srv.db.Query("SELECT face_id, nick, purpose, bio, comments, user_id FROM faceV2beta1 WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	faces := make([]*Face, 0, 5)
	for rows.Next() {
		f := &Face{}
		if err := rows.Scan(&f.ID, &f.Nick, &f.Purpose, &f.Bio, &f.Comments, &f.UserID); err != nil {
			return nil, err
		}
		faces = append(faces, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return faces, nil
}

func (srv *MsgServer) DelFaceByIDV2beta1(faceID string, userID int64) error {
	// _, err := srv.db.Exec("SELECT conn_id FROM connection WHERE face_id = ? OR face_peer_id=", faceID)
	// if err != sql.ErrNoRows || err == nil {
	// 	return errors.New("face deletion error - there are connections for this face")
	// }
	res, err := srv.db.Exec("DELETE FROM faceV2beta1 WHERE face_id = ? AND user_id = ?", faceID, userID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return errors.New("face deletion error")
	}
	return nil
}
