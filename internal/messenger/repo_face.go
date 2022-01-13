package messenger

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

func (srv *Proxy) CreateFace(nick, purpose, bio, comments string, server string, userID string) (string, error) {
	hash := md5.Sum([]byte(nick + "|" + purpose + "|" + bio + "|" + userID))
	face_id := hex.EncodeToString(hash[:])
	res, err := srv.db.Exec("INSERT INTO face (face_id, nick, purpose, bio, comments, server, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)", face_id, nick, purpose, bio, comments, server, userID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return "", err
	}
	return face_id, nil
}

func (srv *Proxy) GetFaceByID(faceID string, userID string) (*Face, error) {
	f := &Face{}
	err := srv.db.QueryRow("SELECT face_id, nick, purpose, bio, comments, server, user_id FROM face WHERE face_id = ? AND user_id = ?", faceID, userID).
		Scan(&f.ID, &f.Nick, &f.Purpose, &f.Bio, &f.Comments, &f.Server, &f.UserID)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return f, nil
}

func (srv *Proxy) GetFacesByUser(userID string) ([]*Face, error) {
	rows, err := srv.db.Query("SELECT face_id, nick, purpose, bio, comments, server, user_id FROM face WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	faces := make([]*Face, 0, 5)
	for rows.Next() {
		f := &Face{}
		if err := rows.Scan(&f.ID, &f.Nick, &f.Purpose, &f.Bio, &f.Comments, &f.Server, &f.UserID); err != nil {
			return nil, err
		}
		faces = append(faces, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return faces, nil
}

func (srv *Proxy) DelFaceByID(faceID string, userID string) error {
	// _, err := srv.db.Exec("SELECT conn_id FROM connection WHERE face_id = ? OR face_peer_id=", faceID)
	// if err != sql.ErrNoRows || err == nil {
	// 	return errors.New("face deletion error - there are connections for this face")
	// }
	res, err := srv.db.Exec("DELETE FROM face WHERE face_id = ? AND user_id = ?", faceID, userID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return errors.New("face deletion error")
	}
	return nil
}

func (srv *Proxy) UpdFace(face *Face) error {
	res, err := srv.db.Exec("UPDATE face SET nick =?, purpose = ?, bio = ?, comments = ?, server = ? WHERE face_id = ? AND user_id = ?",
		&face.Nick, &face.Purpose, &face.Bio, &face.Comments, &face.Server, &face.ID, &face.UserID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return errors.New("face update error")
	}
	return nil
}

/*
func (srv *Proxy) CreateFace(name, description string, userID int64) (string, error) {
	hash := md5.Sum([]byte(name + "|" + description))
	face_id := hex.EncodeToString(hash[:])
	res, err := srv.db.Exec("INSERT INTO face (face_id, name, description, user_id) VALUES (?, ?, ?, ?)", face_id, name, description, userID)
	if err != nil {
		return "", err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return "", err
	}
	return face_id, nil
}

func (srv *Proxy) GetFaceByID(faceID string) (*Face, error) {
	f := &Face{}
	err := srv.db.QueryRow("SELECT face_id, name, description, user_id FROM face WHERE face_id = ?", faceID).Scan(&f.ID, &f.Name, &f.Description, &f.UserID)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return f, nil
}

func (srv *Proxy) GetFacesByUser(userID int64) ([]*Face, error) {
	rows, err := srv.db.Query("SELECT face_id, name, description, user_id FROM face WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	faces := make([]*Face, 0, 5)
	for rows.Next() {
		f := &Face{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.UserID); err != nil {
			return nil, err
		}
		faces = append(faces, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return faces, nil
}

func (srv *Proxy) DelFaceByID(faceID string, userID int64) error {
	_, err := srv.db.Exec("SELECT conn_id FROM connection WHERE face_id = ? OR face_peer_id=", faceID)
	if err != sql.ErrNoRows || err == nil {
		return errors.New("face deletion error - there are connections for this face")
	}
	res, err := srv.db.Exec("DELETE FROM face WHERE face_id = ? AND user_id = ?", faceID, userID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return errors.New("face deletion error")
	}
	return nil
}
*/
