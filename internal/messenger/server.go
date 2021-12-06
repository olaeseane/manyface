package messenger

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/ttacon/chalk"
	"go.uber.org/zap"
	"manyface.net/internal/utils"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

var (
	// TODO: move to configuration?
	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan  = chalk.Cyan.NewStyle().Style
	// red := chalk.Red.NewStyle().Style
	// yellow := chalk.Yellow.NewStyle().Style
)

func NewServer(db *sql.DB, logger *zap.SugaredLogger, mtrxServer string) *MsgServer {
	return &MsgServer{
		wg:         &sync.WaitGroup{},
		db:         db,
		logger:     logger,
		mtrxServer: mtrxServer,
		// clients: make(map[string]MtrxClient),
		conns: make(map[string]*Conn),
	}
}

func (srv *MsgServer) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// TODO: make auth
	c := &Conn{}
	err := srv.db.
		QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE conn_id = ?", req.ConnectionId).
		Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
	if err != nil {
		m := fmt.Sprintf("Can't select connection number %d - %s\n", req.ConnectionId, err)
		srv.logger.Error(m)
		return &SendResponse{Result: false}, errors.New(m)
	}

	if _, ok := srv.conns[c.MtrxUserID]; !ok {
		c.cli, err = joinMtrxClient(c, srv.mtrxServer)
		if err != nil {
			return &SendResponse{Result: false}, err
		}
		c.ch = make(chan EventMessage)
		srv.conns[c.MtrxUserID] = c
		go startMtrxSyncer(srv.wg, c, srv.logger)
	}

	_, err = srv.conns[c.MtrxUserID].cli.SendText(id.RoomID(c.MtrxRoomID), req.Message)
	if err != nil {
		srv.logger.Errorf("Can't send message to %v from %v", c.FacePeerID, c.FaceUserID)
		return &SendResponse{Result: false}, errors.New("send message error")
	}
	srv.logger.Infof("Message was sent to %v from %v", c.FacePeerID, c.FaceUserID)
	return &SendResponse{Result: true}, nil
}

func (srv *MsgServer) Listen(req *ListenRequest, stream Messenger_ListenServer) error {
	// TODO: make auth
	c := &Conn{}
	err := srv.db.
		QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE conn_id = ?", req.ConnectionId).
		Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
	if err != nil {
		m := fmt.Sprintf("Can't select connection number %d - %s\n", req.ConnectionId, err)
		srv.logger.Error(m)
		return errors.New(m)
	}

	c.cli, err = joinMtrxClient(c, srv.mtrxServer)
	if err != nil {
		return err
	}
	c.ch = make(chan EventMessage)
	srv.conns[c.MtrxUserID] = c
	go startMtrxSyncer(srv.wg, c, srv.logger)

	// if ctx.Err() == context.Canceled {
	// 	return status.New(codes.Canceled, "Client cancelled, abandoning.")
	// }

	for evt := range srv.conns[c.MtrxUserID].ch {
		var senderFaceID string
		srv.db.QueryRow("SELECT face_id FROM connection WHERE mtrx_user_id = ?", evt.sender).Scan(&senderFaceID)
		resp := ListenResponse{
			Content:   evt.content,
			Timestamp: evt.timestamp,
			Sender:    senderFaceID,
		}
		if err := stream.Send(&resp); err != nil {
			fmt.Println(err)
			delete(srv.conns, c.MtrxUserID) // TODO: have to close channel, mtrx syncer as well and check memory leaking
			return err
		} else {
			fmt.Printf("stream.Send - \"%s\" from %s (%s) to %s (%s)\n", green(resp.Content), cyan(senderFaceID), evt.sender, cyan(c.FaceUserID), c.MtrxUserID)
		}
	}
	return nil
}

func (srv *MsgServer) CreateFace(name, description string, userID int64) (string, error) {
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

func (srv *MsgServer) GetFaceByID(faceID string) (*Face, error) {
	f := &Face{}
	err := srv.db.QueryRow("SELECT face_id, name, description, user_id FROM face WHERE face_id = ?", faceID).Scan(&f.ID, &f.Name, &f.Description, &f.UserID)
	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	return f, nil
}

func (srv *MsgServer) GetFacesByUser(userID int64) ([]*Face, error) {
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

func (srv *MsgServer) DelFaceByID(faceID string, userID int64) error {
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

func (srv *MsgServer) CreateConn(userID int64, faceUserID, facePeerID string) ([]*Conn, error) {
	var faceUserName, facePeerName string

	// check if faces are corrrected
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ? AND user_id = ?",
		faceUserID, userID).Scan(&faceUserName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ?", facePeerID).Scan(&facePeerName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}

	// can only create one connection for pair faceID-peerID
	var cID int64
	if err := srv.db.QueryRow("SELECT conn_id FROM connection WHERE face_id = ? AND face_peer_id = ?",
		faceUserID, facePeerID).Scan(&cID); err == nil || err != sql.ErrNoRows {
		return nil, errors.New("only one connection for pair of faces")
	}

	// TODO: create once in MsgServer for performance reason?
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
		return createdConns, err
	}
	res, err = tx.Exec("INSERT INTO connection (mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		resPeer.UserID.String(), peerPassword, resPeer.AccessToken, resRoom.RoomID.String(), resUser.UserID.String(), facePeerID, faceUserID)
	if err != nil {
		return createdConns, err
	}
	rowCnt, err = res.RowsAffected()
	if rowCnt != 1 || err != nil {
		return createdConns, err
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

func (srv *MsgServer) DeleteConn(userID int64, faceUserID, facePeerID string) error {
	// TODO: add forget room or delete user or something to stop receive message for this pair connections

	var faceUserName, facePeerName string

	// check if faces are corrrected
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ? AND user_id = ?", faceUserID, userID).Scan(&faceUserName); err == sql.ErrNoRows || err != nil {
		return err
	}
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ?", facePeerID).Scan(&facePeerName); err == sql.ErrNoRows || err != nil {
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

func (srv *MsgServer) GetConnsByUser(userID int64) ([]*Conn, error) {
	rows, err := srv.db.Query("SELECT c.conn_id, c.face_id, c.face_peer_id FROM connection c LEFT JOIN face f ON c.face_id=f.face_id WHERE user_id=?", userID)
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

func joinMtrxClient(c *Conn, mtrxServer string) (*mautrix.Client, error) {
	cli, err := mautrix.NewClient(mtrxServer, "", "")
	if err != nil {
		// logger.Errorf("Can't create client for connection %v", c.ID)
		return nil, err
	}
	_, err = cli.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: c.MtrxUserID},
		Password:         c.MtrxPassword,
		StoreCredentials: true,
	})
	if err != nil {
		// logger.Errorf("Can't client login for connection %v", c.ID)
		return nil, err
	}
	return cli, nil
}

func startMtrxSyncer(wg *sync.WaitGroup, c *Conn, logger *zap.SugaredLogger) {
	defer func() {
		// defer wg.Done()
		c.cli.StopSync()
		close(c.ch)
	}()
	syncer := c.cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		// fmt.Printf("<%s> %s (%s/%s) - %s - %s\n", cyan(string(evt.Sender)), green(evt.Content.AsMessage().Body), evt.Type.String(), evt.ID, evt.RoomID, evt.ToDeviceID)
		select {
		case c.ch <- EventMessage{
			sender:    string(evt.Sender),
			content:   evt.Content.AsMessage().Body,
			roomID:    string(evt.RoomID),
			timestamp: evt.Timestamp,
		}:
		default:
		}
	})

	// TODO: maybe use channel somehow?
	c.cli.Sync()
}
