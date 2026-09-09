package foundry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/lbenedar/fresta/internal/foundry/actions"
	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
	"github.com/lbenedar/fresta/internal/models/db"
	"golang.org/x/sync/errgroup"
)

var (
	ListenIsDone         = errors.New("Listen for websocket data in foundry is stopped")
	ReconnectToWebSocket = errors.New("Reconnect to websocket")
	ChannelIsClosed      = errors.New("Channel is closed")

	CloseTimeoutExceed = errors.New("Timeout of websocket close is exceed")
)

type FoundryApi struct {
	//TODO: make check of admin's authentication
	Transport   *transport.FoundryTransport
	IsAvailable bool

	// Logger *slog.Logger
	Status types.FoundryStatus
	wg     sync.WaitGroup
}

func NewFoundry() *FoundryApi {
	return &FoundryApi{}
}

func (foundry *FoundryApi) SetTransport(tr *transport.FoundryTransport) {
	foundry.Transport = tr
}

func (foundry *FoundryApi) ListenAndServeWS() error {
	var err error

	foundry.Transport.ReadChan = *types.InitWsChannels()
	defer foundry.Transport.ReadChan.Close()

	foundry.background(foundry.Transport.ListenWebSocket)

	err = foundry.ServeWebSocket()
	if err != nil {
		foundry.IsAvailable = false
	}
	return err
}

func (foundry *FoundryApi) ServeWebSocket() error {
	var err error
	wsChannels := &foundry.Transport.ReadChan

	for {
		select {
		case message, ok := <-wsChannels.Msg():
			if !ok {
				foundry.Transport.Logger.Debug("Read channel is closed", "type", types.WebSocketCode)
				foundry.Transport.Http.SessionID = nil
				return ChannelIsClosed
			}
			foundry.Transport.Logger.Debug("Data has been received\n", "msg", message.Code, "type", types.WebSocketCode) //, "msg", message.MsgJson)

			switch message.Code {
			case types.RespPingCode, types.RespSessionDataCode:
				err := foundry.Transport.SendOnlyCodeRequest(message.Code)
				if err != nil {
					wsChannels.Err() <- &types.FoundryError{Direction: types.WriterCode, Err: err, Type: types.WebSocketCode}
					continue
				}
			case types.RespServerChangeCode:
				foundry.IsAvailable = true
				data, err := actions.SplitToTypeAndData([]byte(message.MsgJson))
				if err != nil {
					wsChannels.Err() <- &types.FoundryError{Direction: types.ReaderCode, Err: err, Type: types.WebSocketCode}
					continue
				}

				foundry.Transport.Logger.Debug("RespServerChangeCode", "data", data)
				if data == nil {
					continue
				}

				go func() {
					errAction := data.Action(foundry.Transport, &foundry.Status)

					if errAction != nil {
						wsChannels.Err() <- &types.FoundryError{Direction: types.ReaderCode, Err: errAction, Type: types.WebSocketCode}
					}
				}()
			case types.RespDataCode:
				tr := foundry.Transport
				func() {
					tr.ChanMutex.Lock()
					defer tr.ChanMutex.Unlock()

					tr.ExchangeChan.Msgs[message.Id] = make(chan []byte, 1)
					tr.ExchangeChan.Msgs[message.Id] <- []byte(message.MsgJson)
					go tr.CloseMsgChannel(tr.ExchangeChan.Msgs[message.Id], message.Id, 5*time.Second)

				}()
			default:
			}
		case err = <-wsChannels.Err():
			var foundryErr *types.FoundryError
			if errors.As(err, &foundryErr) {
				if foundryErr.IsFatal {
					foundry.Transport.CloseWebSocketConn()
					foundry.Transport.Http.SessionID = nil
					return foundryErr
				}
				foundry.Transport.Logger.Warn("Got error when listening or served", "err", err.Error(), "type", foundryErr.Type, "direction", foundryErr.Direction)
			}
		case <-wsChannels.Done():
			foundry.Transport.CloseWebSocketConn()
			foundry.Transport.Http.SessionID = nil
			return ListenIsDone
		case <-wsChannels.Reconnect():
			foundry.Transport.CloseWebSocketConn()
			return ReconnectToWebSocket
		}
	}
}

func (foundry *FoundryApi) HandleWSRequest(msgType string) ([]byte, error) {
	msg := types.NewWsMessageByPage(msgType, foundry.Transport.CurrWsId)

	return foundry.Transport.HandleWebsocketRequest(msg)
}

func (foundry *FoundryApi) Shutdown() error {
	if !foundry.IsAvailable {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	foundry.Transport.Logger.Info("Completing foundry background tasks")

	foundryClosed := make(chan struct{})
	go func() {
		types.CloseChannel(foundry.Transport.ReadChan.Done())

		foundry.wg.Wait()
		foundry.Transport.ExchangeChan.Close()
		foundry.Transport.ReadChan.Close()

		types.CloseChannel(foundryClosed)
	}()

	select {
	case <-foundryClosed:
		foundry.Transport.Logger.Info("Stopped foundry server")
	case <-ctx.Done():
		return CloseTimeoutExceed
	}

	return nil
}

func (foundry *FoundryApi) ConnectToWebSocket(ctx context.Context) (bool, error) {
	group, _ := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		if !foundry.Transport.HasSessionId() {
			err = foundry.Transport.InitSessionId()
			if err != nil {
				return err
			}

			err = foundry.Transport.ConnectToFoundry()
			if err != nil {
				return err
			}
		}

		err = foundry.Transport.InitWebSocketConnection()
		if err != nil {
			return err
		}

		status, err := foundry.Transport.Http.GetStatus()
		if err != nil {
			return err
		}
		foundry.Status = *types.NewFoundryStatus(status)

		return nil
	})

	err := group.Wait()
	if err != nil {
		return false, err
	}

	err = foundry.ListenAndServeWS()
	if err != nil && err != ListenIsDone {
		return true, err
	}

	return true, nil
}

func (foundry *FoundryApi) StartListenFoundry() {
	foundry.background(func() {
		timeInterval := 1 * time.Second

		for {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			ok, err := foundry.ConnectToWebSocket(ctx)
			if err != nil {
				if errors.Is(err, ReconnectToWebSocket) {
					if foundry.Transport.ReconnectNum >= foundry.Transport.ReconnectNumMax {
						return
					}
					time.Sleep(foundry.Transport.ReconnectTimeout)
					foundry.Transport.ReconnectNum++
					foundry.Transport.Logger.Info("Reconnecting to WebSocket", "times", foundry.Transport.ReconnectNum)
					continue
				}

				foundry.Transport.Logger.Error("Error raised", "err", err.Error())

				if ok {
					timeInterval = 1 * time.Second
				}
				foundry.Transport.Logger.Info("Trying to reconnect", "timer", timeInterval.String())

				time.Sleep(timeInterval)
				timeInterval = min(timeInterval*2, 15*time.Second)
				continue
			}
			return
		}
	})
}

func (foundry *FoundryApi) PrepareDB() error {
	dbConn := foundry.Transport.DB

	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() error { return db.DeleteSetupAll(dbConn) })
	group.Go(func() error { return db.DeleteGameAll(dbConn) })
	group.Go(func() error { return db.DeleteSeqAll(dbConn) })

	err := group.Wait()
	if err != nil && !errors.Is(err, db.ErrorRecordNotFound) {
		return err
	}

	return nil
}

// func (foundry *FoundryApi) Test() (*json_models.Game, error) {
// 	wsMsg := types.NewWsMessage("world", foundry.Transport.CurrWsId)
// 	foundry.Transport.CurrWsId++

// 	answer, err := foundry.Transport.HandleWebsocketRequest(wsMsg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var game []json_models.Game

// 	start := time.Now()
// 	err = json.Unmarshal(answer, &game)
// 	if err != nil {
// 		return nil, err
// 	}
// 	elapsed := time.Since(start)

// 	foundry.Transport.Logger.Info("World data successfully received and parsed", "parseTime", elapsed)
// 	return &game[0], nil
// }
