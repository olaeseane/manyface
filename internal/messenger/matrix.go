package messenger

import (
	sync "sync"

	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

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

func startMtrxSyncer(wg *sync.WaitGroup, c *Conn, logger *zap.SugaredLogger) { // TODO: remove WG?
	defer func() {
		// defer wg.Done()
		c.cli.StopSync()
		close(c.ch)
	}()
	syncer := c.cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		// fmt.Printf("<%s> %s (%s/%s) - %s - %s\n", cyan(string(evt.Sender)), green(evt.Content.AsMessage().Body), evt.Type.String(), evt.ID, evt.RoomID, evt.ToDeviceID)
		/*select {
		case c.ch <- EventMessage{
			sender:    string(evt.Sender),
			content:   evt.Content.AsMessage().Body,
			roomID:    string(evt.RoomID),
			timestamp: evt.Timestamp,
		}:
			// case <- done: TODO: instead of default:
			// return
		default: // TODO: remove?
		}*/

		c.ch <- EventMessage{
			sender:    string(evt.Sender),
			content:   evt.Content.AsMessage().Body,
			roomID:    string(evt.RoomID),
			mtype:     evt.Type.String(),
			timestamp: evt.Timestamp,
		}
	})

	// TODO: maybe use channel somehow?
	c.cli.Sync()
}
