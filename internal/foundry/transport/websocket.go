package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lbenedar/fresta/internal/foundry/requests"
	"github.com/lbenedar/fresta/internal/foundry/types"
	json_models "github.com/lbenedar/fresta/internal/models/json"
)

func (tr *FoundryTransport) InitWebSocketConnection() error {
	if !tr.IsAuth {
		return ErrorNotAuth
	}

	query := url.Values{}
	query.Add("session", *tr.Http.SessionID)
	query.Add("EIO", "4")
	query.Add("transport", "websocket")

	u := url.URL{Scheme: "ws", Host: tr.Http.Host, Path: fmt.Sprintf("%s/", requests.SocketPath), RawQuery: query.Encode()}

	wsHeader := http.Header{}
	wsHeader.Set("Cookie", fmt.Sprintf("session=%s", *tr.Http.SessionID))
	wsHeader.Set("Pragma", "no-cache")
	wsHeader.Set("Cache-Control", "no-cache")

	dialer := websocket.DefaultDialer
	dialer.ReadBufferSize = 128 * 1024 * 1024
	dialer.WriteBufferSize = 128 * 1024 * 1024

	wsConn, _, err := dialer.Dial(u.String(), wsHeader)
	if err != nil {
		return err
	}

	tr.WsConn = wsConn
	tr.CurrWsId = 0
	tr.Logger.Info("Successefully connected to Foundry Websocket")
	return nil
}

func (tr *FoundryTransport) HandleWebsocketRequest(msg *types.WsMessage) ([]byte, error) {
	tr.Logger.Debug("WS: Data has been send\n", "msg", msg.ToString())
	err := tr.WsConn.WriteMessage(websocket.TextMessage, msg.ToByteSlice())
	if err != nil {
		return nil, err
	}

	data, err := tr.ReceiveMessage(msg.Id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tr *FoundryTransport) ReceiveMessage(id int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case msg, ok := <-tr.ExchangeChan.Msgs[id]:
			if !ok {
				time.Sleep(5 * time.Microsecond)
				continue
			}
			tr.CloseMsgChannel(tr.ExchangeChan.Msgs[id], id, 0)
			return msg, nil
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil {
				return nil, err
			}
			return nil, ErrorTimeout
		default:
			time.Sleep(5 * time.Microsecond)
		}
	}
}

func (tr *FoundryTransport) ListenWebSocket() {
	for {
		_, message, err := tr.WsConn.ReadMessage()
		if err != nil {
			select {
			case <-tr.ReadChan.Done():
				return
			case <-tr.ReadChan.Reconnect():
				return
			default:
				tr.ReadChan.Err() <- &types.FoundryError{Direction: types.ReaderCode, Err: err, Type: types.WebSocketCode, IsFatal: true}
			}
			return
		}
		data, err := types.ParseWsRespMessage(string(message), types.RequestCodes)
		if err != nil {
			tr.ReadChan.Err() <- &types.FoundryError{Direction: types.ReaderCode, Err: err, Type: types.WebSocketCode}
			continue
		}
		tr.ReadChan.Msg() <- data
	}
}

func (tr *FoundryTransport) ConnectToWorld(userId string, userPass string) error {
	err := tr.LogInToWorld(userId, userPass)
	if err != nil {
		return err
	}
	tr.Logger.Info("World is started. Succesfully logged into world")
	tr.LoggedInChan = make(chan bool)

	close(tr.ReadChan.Reconnect())

	val, ok := <-tr.LoggedInChan
	if !ok {
		return ErrorChannelIsClosed
	}
	if !val {
		return ErrorUserIsNotConnected
	}

	return nil
}

func (tr *FoundryTransport) GetUserId(username string) (string, error) {
	msgJson, err := tr.GetJsonDataByType(requests.JoinPath)
	if err != nil {
		return "", err
	}

	var modelUsers []struct {
		Users []json_models.User `json:"users"`
	}

	err = json.Unmarshal(msgJson, &modelUsers)
	if err != nil {
		return "", err
	}

	tr.Logger.Debug("User data", "msg in json format", modelUsers)

	users := modelUsers[0].Users
	for i := range users {
		if strings.EqualFold(users[i].Name, username) {
			return users[i].ID, nil
		}
	}

	return "", ErrorUserDoesntExist
}

func (tr *FoundryTransport) GetUserIdAndPass() (string, string, error) {
	status, err := tr.Http.GetStatus()
	if err != nil {
		return "", "", err
	}

	var worldData *types.WorldData
	for i := range tr.InitWorlds {
		if strings.EqualFold(tr.InitWorlds[i].Name, status.World) {
			worldData = &tr.InitWorlds[i]
			break
		}
	}

	if worldData == nil {
		return "", "", err
	}

	userId, err := tr.GetUserId(worldData.Username)
	if err != nil {
		return "", "", err
	}
	return userId, worldData.UserPass, nil
}

func (tr *FoundryTransport) SendOnlyCodeRequest(code string) error {
	tr.Logger.Debug("WS: Data has been send\n", "msg", types.CodesRespToReq[code])

	return tr.WsConn.WriteMessage(websocket.TextMessage, []byte(types.CodesRespToReq[code]))
}

func (tr *FoundryTransport) GetJsonDataByType(stateType string) ([]byte, error) {
	msg := types.NewWsMessageByPage(stateType, tr.CurrWsId)
	tr.CurrWsId++

	return tr.HandleWebsocketRequest(msg)
}

func (tr *FoundryTransport) GetJsonData(msg string) ([]byte, error) {
	wsMsg := types.NewWsMessage(msg, tr.CurrWsId)
	tr.CurrWsId++

	return tr.HandleWebsocketRequest(wsMsg)
}

func (tr *FoundryTransport) CloseWebSocketConn() error {
	tr.Logger.Debug("Websocket has been closed")
	return tr.WsConn.Close()
}
