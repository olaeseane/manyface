package messenger

import (
	"context"
	"errors"
	"fmt"

	"github.com/ttacon/chalk"

	"maunium.net/go/mautrix/id"
)

var (
	// TODO: move to configuration?
	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan  = chalk.Cyan.NewStyle().Style
	// red := chalk.Red.NewStyle().Style
	// yellow := chalk.Yellow.NewStyle().Style
)

func (srv *Proxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
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

func (srv *Proxy) Listen(req *ListenRequest, stream Messenger_ListenServer) error {
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

	// TODO: change key MtrxUserID to connID
	if _, ok := srv.conns[c.MtrxUserID]; !ok {
		c.cli, err = joinMtrxClient(c, srv.mtrxServer)
		if err != nil {
			return err
		}
		c.ch = make(chan EventMessage)
		srv.conns[c.MtrxUserID] = c
		go startMtrxSyncer(srv.wg, c, srv.logger)
	}

	// if ctx.Err() == context.Canceled {
	// 	return status.New(codes.Canceled, "Client cancelled, abandoning.")
	// }

	for evt := range srv.conns[c.MtrxUserID].ch {
		var senderFaceID, receiveFaceID string
		srv.db.QueryRow("SELECT face_id, face_peer_id FROM connection WHERE mtrx_user_id = ?", evt.sender).Scan(&senderFaceID, &receiveFaceID)
		resp := ListenResponse{
			Content:        evt.content,
			Timestamp:      evt.timestamp,
			Sender:         senderFaceID,
			ReceiverFaceId: receiveFaceID,
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
