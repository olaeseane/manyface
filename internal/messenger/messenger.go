package messenger

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	sync "sync"

	"github.com/ttacon/chalk"
	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// TODO: move to configuration?
	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan  = chalk.Cyan.NewStyle().Style
	red   = chalk.Red.NewStyle().Style
	// yellow := chalk.Yellow.NewStyle().Style
)

func NewProxy(db *sql.DB, logger *zap.SugaredLogger, mtrxServer string) *Proxy {
	return &Proxy{
		db:         db,
		logger:     logger,
		mtrxServer: mtrxServer,
		mu:         &sync.RWMutex{},
	}
}

func (proxy *Proxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// TODO: make auth
	var c *Conn
	if v, ok := proxy.conns.Load(req.ConnectionId); !ok {
		err := proxy.db.
			QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE conn_id = ?", req.ConnectionId).
			Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
		if err != nil {
			m := fmt.Sprintf("Can't select connection number %d - %s\n", req.ConnectionId, err)
			proxy.logger.Error(m)
			return nil, status.Errorf(codes.InvalidArgument, m)
		}
		c.cli, err = proxy.joinMtrxClient(c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		c.quitCh = make(chan struct{})
		c.dataCh = proxy.startMtrxSyncer(c)
		proxy.conns.Store(req.ConnectionId, c)
	} else {
		c, ok = v.(*Conn)
		if !ok {
			return nil, status.Errorf(codes.Internal, errors.New("error cast").Error())
		}
	}

	_, err := c.cli.SendText(id.RoomID(c.MtrxRoomID), req.Message)
	if err != nil {
		m := fmt.Sprintf("Can't send message to %v from %v", c.FacePeerID, c.FaceUserID)
		proxy.logger.Errorf(m)
		return nil, status.Errorf(codes.Internal, m)
	}
	proxy.logger.Infof("Message was sent to %v from %v", c.FacePeerID, c.FaceUserID)
	return &SendResponse{ConnectionId: c.ID}, nil
}

func (proxy *Proxy) Listen(req *ListenRequest, stream Messenger_ListenServer) error {
	c := &Conn{}
	if v, ok := proxy.conns.Load(req.ConnectionId); !ok {
		err := proxy.db.
			QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE conn_id = ?", req.ConnectionId).
			Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
		if err != nil {
			m := fmt.Sprintf("Can't select connection number %d - %s\n", req.ConnectionId, err)
			proxy.logger.Error(m)
			return status.Errorf(codes.InvalidArgument, m)
		}
		c.cli, err = proxy.joinMtrxClient(c)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		c.quitCh = make(chan struct{})
		c.dataCh = proxy.startMtrxSyncer(c)
		proxy.conns.Store(req.ConnectionId, c)
	} else {
		c, ok = v.(*Conn)
		if !ok {
			return status.Errorf(codes.Internal, errors.New("error cast").Error())
		}
	}

	log.Printf("Starting stream for connection ID: %d", req.ConnectionId)

	ctx := stream.Context()
	for {
		select {
		case evt := <-c.dataCh:
			var senderFaceID, receiveFaceID string
			proxy.db.QueryRow("SELECT face_id, face_peer_id FROM connection WHERE mtrx_user_id = ?", evt.sender).Scan(&senderFaceID, &receiveFaceID)
			resp := ListenResponse{
				Message:        evt.content,
				Timestamp:      evt.timestamp,
				Sender:         evt.sender,
				SenderFaceId:   senderFaceID,
				ReceiverFaceId: receiveFaceID,
			}
			if err := stream.Send(&resp); err != nil {
				fmt.Println(red(err.Error()))
				return status.Errorf(codes.Internal, err.Error())
			} else {
				fmt.Printf("stream.Send - \"%s\" from %s (%s) to %s (%s)\n", green(resp.Message), cyan(senderFaceID), evt.sender, cyan(c.FaceUserID), c.MtrxUserID)
			}
		case <-ctx.Done():
			log.Printf("Closing stream for connection ID: %d", req.ConnectionId)
			c.quitCh <- struct{}{}
			c.cli.LogoutAll()
			proxy.conns.Delete(req.ConnectionId)
			return nil
		}
	}

}

/*
func (proxy *Proxy) Listen(req *ListenRequest, stream Messenger_ListenServer) error {
	// TODO: make auth
	c := &Conn{}
	err := proxy.db.
		QueryRow("SELECT conn_id, mtrx_user_id, mtrx_password, mtrx_access_token, mtrx_room_id, mtrx_peer_id, face_id, face_peer_id FROM connection WHERE conn_id = ?", req.ConnectionId).
		Scan(&c.ID, &c.MtrxUserID, &c.MtrxPassword, &c.MtrxAccessToken, &c.MtrxRoomID, &c.MtrxPeerID, &c.FaceUserID, &c.FacePeerID)
	if err != nil {
		m := fmt.Sprintf("Can't select connection number %d - %s\n", req.ConnectionId, err)
		proxy.logger.Error(m)
		return status.Errorf(codes.InvalidArgument, m)
	}

	// TODO: change key MtrxUserID to connID
	if _, ok := proxy.conns[c.MtrxUserID]; !ok {
		c.cli, err = proxy.joinMtrxClient(c)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		c.quitCh = make(chan struct{})
		c.dataCh = proxy.startMtrxSyncer(c)
		proxy.mu.Lock()
		proxy.conns[c.MtrxUserID] = c
		proxy.mu.Unlock()
	}

	ctx := stream.Context()
	// stream.Con
	// if stream.ctx.Err() == context.Canceled {
	// 	return status.New(codes.Canceled, "Client cancelled, abandoning.")
	// }

	for {
		proxy.mu.RLock()
		select {
		case evt := <-proxy.conns[c.MtrxUserID].dataCh:
			proxy.mu.RUnlock()
			var senderFaceID, receiveFaceID string
			proxy.db.QueryRow("SELECT face_id, face_peer_id FROM connection WHERE mtrx_user_id = ?", evt.sender).Scan(&senderFaceID, &receiveFaceID)
			resp := ListenResponse{
				Message:        evt.content,
				Timestamp:      evt.timestamp,
				Sender:         evt.sender,
				SenderFaceId:   senderFaceID,
				ReceiverFaceId: receiveFaceID,
			}
			if err := stream.Send(&resp); err != nil {
				fmt.Println(red(err.Error()))
				// proxy.conns[c.MtrxUserID].ch <- EventMessage{}
				// delete(proxy.conns, c.MtrxUserID) // TODO: have to close channel, mtrx syncer as well and check memory leaking
				// return status.Errorf(codes.Internal, err.Error())
			} else {
				fmt.Printf("stream.Send - \"%s\" from %s (%s) to %s (%s)\n", green(resp.Message), cyan(senderFaceID), evt.sender, cyan(c.FaceUserID), c.MtrxUserID)
			}
		default:
			if ctx.Err() == context.Canceled {
				fmt.Println("Context().Done()")
				proxy.conns[c.MtrxUserID].quitCh <- struct{}{}
				proxy.conns[c.MtrxUserID].cli.LogoutAll()
				delete(proxy.conns, c.MtrxUserID) // TODO: have to close channel, mtrx syncer as well and check memory leaking
				return nil
			}
		}
	}
}
*/

func (proxy *Proxy) joinMtrxClient(c *Conn) (*mautrix.Client, error) {
	cli, err := mautrix.NewClient(proxy.mtrxServer, "", "")
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

func (proxy *Proxy) startMtrxSyncer(c *Conn) chan EventMessage {
	ch := make(chan EventMessage)
	syncer := c.cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		// fmt.Printf("<%s> %s (%s/%s) - %s - %s\n", cyan(string(evt.Sender)), green(evt.Content.AsMessage().Body), evt.Type.String(), evt.ID, evt.RoomID, evt.ToDeviceID)
		ch <- EventMessage{
			sender:    string(evt.Sender),
			content:   evt.Content.AsMessage().Body,
			roomID:    string(evt.RoomID),
			mtype:     evt.Type.String(),
			timestamp: evt.Timestamp,
		}
	})

	go c.cli.Sync()
	go func() {
		<-c.quitCh
		c.cli.StopSync()
		close(ch)
	}()

	return ch
}
