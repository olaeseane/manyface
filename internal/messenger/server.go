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

	pb "manyface.net/grpc"
	"manyface.net/internal/common"
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

func NewServer(db *sql.DB, logger *zap.SugaredLogger) *MsgServer {
	return &MsgServer{
		wg:     &sync.WaitGroup{},
		db:     db,
		logger: logger,
		// clients: make(map[string]MtrxClient),
		conns: make(map[string]*Conn),
	}
}

func (srv *MsgServer) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {
	// TODO: make auth
	c := &Conn{}
	err := srv.db.
		QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE face_id = ? AND face_peer_id = ?", req.SenderFaceId, req.ReceiverFaceId).
		Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
	if err != nil {
		m := fmt.Sprintf("Can't select connection for faces %v and %v - %v\n", req.SenderFaceId, req.ReceiverFaceId, err)
		srv.logger.Error(m)
		return &pb.SendResponse{Result: false}, errors.New(m)
	}

	if _, ok := srv.conns[c.MtrxUserID]; !ok {
		c.cli, err = joinMtrxClient(c)
		if err != nil {
			return &pb.SendResponse{Result: false}, err
		}
		c.ch = make(chan EventMessage)
		srv.conns[c.MtrxUserID] = c
		go startMtrxSyncer(srv.wg, c, srv.logger)
	}

	_, err = srv.conns[c.MtrxUserID].cli.SendText(id.RoomID(c.MtrxRoomID), req.Message)
	if err != nil {
		srv.logger.Errorf("Can't send message to %v from %v", req.ReceiverFaceId, req.SenderFaceId)
		return &pb.SendResponse{Result: false}, errors.New("send message error")
	}
	srv.logger.Infof("Message was sent to %v from %v", req.ReceiverFaceId, req.SenderFaceId)
	return &pb.SendResponse{Result: true}, nil
}

func (srv *MsgServer) Listen(req *pb.ListenRequest, stream pb.Messenger_ListenServer) error {
	// TODO: make auth
	c := &Conn{}
	err := srv.db.
		QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE face_id = ? AND face_peer_id = ?", req.SenderFaceId, req.ReceiverFaceId).
		Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
	if err != nil {
		m := fmt.Sprintf("Can't select connection for faces %v and %v - %v\n", req.SenderFaceId, req.ReceiverFaceId, err)
		srv.logger.Error(m)
		return errors.New(m)
	}

	if _, ok := srv.conns[c.MtrxUserID]; !ok {
		c.cli, err = joinMtrxClient(c)
		if err != nil {
			return err
		}
		c.ch = make(chan EventMessage)
		srv.conns[c.MtrxUserID] = c
		go startMtrxSyncer(srv.wg, c, srv.logger)
	}

	for evt := range srv.conns[c.MtrxUserID].ch {
		resp := pb.ListenResponse{
			Content:        evt.content,
			Timestamp:      evt.timestamp,
			SenderFaceId:   evt.sender,
			ReceiverFaceId: "",
		}
		if err := stream.Send(&resp); err != nil {
			return err
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
	faces := make([]*Face, 0, 5)
	rows, err := srv.db.Query("SELECT face_id, name, description, user_id FROM face WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	// TODO: can't create new connection if there is already one for these faces?
	var faceUserName, facePeerName string

	// check if faces are corrrected
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ? AND user_id = ?", faceUserID, userID).Scan(&faceUserName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}
	if err := srv.db.QueryRow("SELECT name FROM face WHERE face_id = ?", facePeerID).Scan(&facePeerName); err == sql.ErrNoRows || err != nil {
		return nil, err
	}

	// TODO: create once in MsgServer for performance reason?
	cli, err := mautrix.NewClient("http://localhost:8008", "", "") // TODO: change to config env
	if err != nil {
		return nil, err
	}

	userPassword := common.RandStringRunes(10)
	resUser, err := cli.RegisterDummy(&mautrix.ReqRegister{
		Username: common.RandStringRunes(10),
		Password: userPassword,
	})
	if err != nil {
		return nil, err
	}
	peerPassword := common.RandStringRunes(10)
	resPeer, err := cli.RegisterDummy(&mautrix.ReqRegister{
		Username: common.RandStringRunes(10),
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
			Name:     faceUserName + "-" + facePeerName,
			Invite:   []id.UserID{resPeer.UserID},
			Preset:   "trusted_private_chat",
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

/*
func (srv *MsgServer) StartSyncers() {
	rows, err := srv.db.Query("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection")
	if err != nil {
		srv.logger.Fatal("Can't read connections - ", err)
	}
	// var conns []*Conn
	for rows.Next() {
		c := &Conn{}
		if err := rows.Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID); err != nil {
			srv.logger.Errorf("Can't scan attributes of connection %v", c.ID)
			continue
		}
		srv.conns[c.MtrxUserID] = c
		// conns = append(conns, c)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()
	for _, c := range srv.conns {
		srv.wg.Add(1)
		go startConnSync(srv.wg, c, srv.logger) // TODO: maybe use channel somehow?
	}
	srv.wg.Wait()
}
*/

func joinMtrxClient(c *Conn) (*mautrix.Client, error) {
	cli, err := mautrix.NewClient("http://localhost:8008", "", "") // TODO: change to config env
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
		fmt.Printf("<%s> %s (%s/%s) - %s - %s\n", cyan(string(evt.Sender)), green(evt.Content.AsMessage().Body), evt.Type.String(), evt.ID, evt.RoomID, evt.ToDeviceID)
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
