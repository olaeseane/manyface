package messenger

import (
	"database/sql"
	"errors"

	"manyface.net/internal/utils"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

func (srv *Proxy) CreateConn(userID string, faceUserID, facePeerID string) ([]*Conn, error) {
	var faceUserName, facePeerName string

	// check if faces are corrrected
	if err := srv.db.QueryRow("SELECT nick FROM face WHERE face_id = ? AND user_id = ?",
		faceUserID, userID).Scan(&faceUserName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	if err := srv.db.QueryRow("SELECT nick FROM face WHERE face_id = ?", facePeerID).Scan(&facePeerName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}

	// can only create one connection for pair faceID-peerID
	var cID int64
	if err := srv.db.QueryRow("SELECT conn_id FROM connection WHERE face_id = ? AND face_peer_id = ?",
		faceUserID, facePeerID).Scan(&cID); err == nil || err != sql.ErrNoRows {
		return nil, errors.New("only one connection for pair of faces")
	}

	// TODO: create once in Proxy for performance reason?
	cli, err := mautrix.NewClient(srv.mtrxServer, "", "")
	if err != nil {
		return nil, err
	}

	userPassword := utils.RandStringRunes(10)
	resUser, err := cli.RegisterDummy(&mautrix.ReqRegister{
		Username: utils.RandStringRunes(10),
		Password: userPassword,
	})
	if err != nil {
		return nil, err
	}
	peerPassword := utils.RandStringRunes(10)
	resPeer, err := cli.RegisterDummy(&mautrix.ReqRegister{
		Username: utils.RandStringRunes(10),
		Password: peerPassword,
	})
	if err != nil {
		return nil, err
	}

	cli.SetCredentials(resUser.UserID, resUser.AccessToken)
	resRoom, err := cli.CreateRoom(
		&mautrix.ReqCreateRoom{
			Visibility: "private",
			// RoomAliasName: faceUserName + "-" + facePeerName,
			Name:   faceUserName + "<->" + facePeerName,
			Invite: []id.UserID{resPeer.UserID},
			// Preset:   "trusted_private_chat",
			Preset:   "private_chat",
			IsDirect: true,
		})
	if err != nil {
		return nil, err
	}
	tx, _ := srv.db.Begin()
	defer tx.Rollback()
	res, err := tx.Exec("INSERT INTO connection (mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		resUser.UserID.String(), userPassword, resUser.AccessToken, resRoom.RoomID.String(), resPeer.UserID.String(), faceUserID, facePeerID)
	if err != nil {
		return nil, err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return nil, err
	}
	connID, _ := res.LastInsertId()
	createdConns := make([]*Conn, 2)
	createdConns[0] = &Conn{
		ID:              connID,
		MtrxUserID:      resUser.UserID.String(),
		MtrxPassword:    userPassword,
		MtrxAccessToken: resUser.AccessToken,
		MtrxRoomID:      resRoom.RoomID.String(),
		MtrxPeerID:      resPeer.UserID.String(),
		FaceUserID:      faceUserID,
	}
	// srv.wg.Add(1)
	// srv.conns[createdConns[0].MtrxUserID] = createdConns[0]
	// go startConnSync(srv.wg, createdConns[0], srv.logger) // TODO: maybe use channel somehow?
	cli.Logout()
	cli.ClearCredentials()

	cli.SetCredentials(resPeer.UserID, resPeer.AccessToken)
	if _, err := cli.JoinRoomByID(resRoom.RoomID); err != nil {
		return nil, err
	}
	res, err = tx.Exec("INSERT INTO connection (mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		resPeer.UserID.String(), peerPassword, resPeer.AccessToken, resRoom.RoomID.String(), resUser.UserID.String(), facePeerID, faceUserID)
	if err != nil {
		return nil, err
	}
	rowCnt, err = res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return nil, err
	}
	connID, _ = res.LastInsertId()
	createdConns[1] = &Conn{
		ID:              connID,
		MtrxUserID:      resPeer.UserID.String(),
		MtrxPassword:    peerPassword,
		MtrxAccessToken: resPeer.AccessToken,
		MtrxRoomID:      resRoom.RoomID.String(),
		MtrxPeerID:      resUser.UserID.String(),
		FaceUserID:      facePeerID,
	}
	// srv.wg.Add(1)
	// srv.conns[createdConns[1].MtrxUserID] = createdConns[1]
	// go startConnSync(srv.wg, createdConns[1], srv.logger) // TODO: maybe use channel somehow?
	cli.Logout()
	cli.ClearCredentials()
	tx.Commit()

	return createdConns, nil
}

func (srv *Proxy) DeleteConn(userID string, faceUserID, facePeerID string) error {
	// TODO: add forget room or delete user or something to stop receive message for this pair connections

	var faceUserName, facePeerName string

	// check if faces are corrrected
	if err := srv.db.QueryRow("SELECT nick FROM face WHERE face_id = ? AND user_id = ?", faceUserID, userID).Scan(&faceUserName); err == sql.ErrNoRows || err != nil {
		return err
	}
	if err := srv.db.QueryRow("SELECT nick FROM face WHERE face_id = ?", facePeerID).Scan(&facePeerName); err == sql.ErrNoRows || err != nil {
		return err
	}

	tx, _ := srv.db.Begin()
	defer tx.Rollback()
	res, err := tx.Exec("DELETE FROM connection WHERE face_id = ? AND face_peer_id = ?", faceUserID, facePeerID)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return err
	}
	res, err = tx.Exec("DELETE FROM connection WHERE face_id = ? AND face_peer_id = ?", facePeerID, faceUserID)
	if err != nil {
		return err
	}
	rowCnt, err = res.RowsAffected()
	if err != nil || rowCnt != 1 {
		return err
	}
	tx.Commit()
	return nil
}

func (srv *Proxy) GetConnsByUser(userID string) ([]*Conn, error) {
	rows, err := srv.db.Query("SELECT c.conn_id, c.face_id, c.face_peer_id FROM connection c LEFT JOIN face f ON c.face_id=f.face_id WHERE f.user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	conns := make([]*Conn, 0, 5)
	for rows.Next() {
		c := &Conn{}
		if err := rows.Scan(&c.ID, &c.FaceUserID, &c.FacePeerID); err != nil {
			return nil, err
		}
		conns = append(conns, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return conns, nil
}
