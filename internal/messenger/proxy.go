package messenger

import (
	"database/sql"
	sync "sync"

	"go.uber.org/zap"
)

func NewProxy(db *sql.DB, logger *zap.SugaredLogger, mtrxServer string) *Proxy {
	return &Proxy{
		wg:         &sync.WaitGroup{},
		db:         db,
		logger:     logger,
		mtrxServer: mtrxServer,
		// clients: make(map[string]MtrxClient),
		conns: make(map[string]*Conn),
	}
}
