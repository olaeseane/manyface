package messenger

import (
	"database/sql"
	"sync"

	"go.uber.org/zap"
	"manyface.net/internal/blobstorage"
	"manyface.net/internal/session"
	"maunium.net/go/mautrix"
)

type MsgServer struct {
	wg         *sync.WaitGroup
	db         *sql.DB
	logger     *zap.SugaredLogger
	mtrxServer string
	conns      map[string]*Conn

	UnimplementedMessengerServer
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
	BS     *blobstorage.FSStorage
}

type Face struct {
	ID       string `json:"face_id" example:"98db6968134c5a069d6a513c3993c8af"`
	Nick     string `json:"nick" example:"Bob"`
	Purpose  string `json:"purpose" example:"For work"`
	Bio      string `json:"bio" example:"Some bio"`
	Comments string `json:"comments" example:"Some comments"`
	UserID   int64  `json:"user_id" example:"10"`
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
