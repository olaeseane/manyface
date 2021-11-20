package messenger

import (
	"database/sql"
	"sync"

	"go.uber.org/zap"
	"manyface.net/internal/session"
	"maunium.net/go/mautrix"

	pb "manyface.net/grpc"
)

type MsgServer struct {
	wg     *sync.WaitGroup
	db     *sql.DB
	logger *zap.SugaredLogger
	conns  map[string]*Conn

	pb.UnimplementedMessengerServer
}

type EventMessage struct {
	content   string
	sender    string
	roomID    string
	timestamp int64
}

type MessengerHandler struct {
	Logger *zap.SugaredLogger
	Srv    *MsgServer
	SM     *session.SessionManager
}

type Face struct {
	ID          string `json:"face_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int64  `json:"user_id"`
}

type Conn struct {
	ID              int64  `json:"conn_id,omitempty"`
	MtrxUserID      string `json:"mtrx_user_id,omitempty"`
	MtrxPassword    string `json:"mtrx_password,omitempty"`
	MtrxAccessToken string `json:"mtrx_access_token,omitempty"`
	MtrxRoomID      string `json:"mtrx_room_id,omitempty"`
	MtrxPeerID      string `json:"mtrx_peer_id,omitempty"`
	FaceUserID      string `json:"face_user_id,omitempty"`
	FacePeerID      string `json:"face_peer_id,omitempty"`
	cli             *mautrix.Client
	ch              chan EventMessage
}
