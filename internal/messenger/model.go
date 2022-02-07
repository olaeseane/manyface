package messenger

import (
	"database/sql"
	"sync"

	"go.uber.org/zap"
	"manyface.net/internal/blobstorage"
	"manyface.net/internal/session"
	"maunium.net/go/mautrix"
)

type Proxy struct {
	db         *sql.DB
	logger     *zap.SugaredLogger
	mtrxServer string
	conns      sync.Map
	mu         *sync.RWMutex

	UnimplementedMessengerServer
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
	dataCh          chan EventMessage
	quitCh          chan struct{}
}

type EventMessage struct {
	content   string
	sender    string
	roomID    string
	mtype     string
	timestamp int64
}

type MessengerHandler struct {
	Logger *zap.SugaredLogger
	Srv    *Proxy
	SM     *session.SessionManager
	BS     *blobstorage.FSStorage
}

type Face struct {
	ID       string `json:"face_id,omitempty" example:"98db6968134c5a069d6a513c3993c8af"`
	Nick     string `json:"nick,omitempty" example:"Bob"`
	Purpose  string `json:"purpose,omitempty" example:"For work"`
	Bio      string `json:"bio,omitempty" example:"Some bio"`
	Comments string `json:"comments,omitempty" example:"Some comments"`
	Server   string `json:"server,omitempty" example:"manyface.net"`
	UserID   string `json:"user_id,omitempty" example:"10"`
}
