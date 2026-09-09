package transport

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/lbenedar/fresta/internal/foundry/requests"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

type FoundryTransport struct {
	WsConn       *websocket.Conn
	CurrWsId     int
	ChanMutex    sync.Mutex
	ReadChan     types.ReadChannels
	ExchangeChan types.ExchangeChannels
	LoggedInChan chan bool

	Http             *requests.FoundryHttpRequest
	ReconnectTimeout time.Duration
	ReconnectNum     int
	ReconnectNumMax  int
	IsAuth           bool
	IsLogin          bool

	// Models   *db.Models
	DB         *sqlx.DB
	IsDbInit   bool
	InitWorlds []types.WorldData

	Logger *slog.Logger
}

type FoundryTransportData struct {
	DbConn     *sqlx.DB
	Logger     *slog.Logger
	HttpConfig *requests.FoundryHttpRequest
	Worlds     []types.WorldData
}

func NewFoundryTransport(data *FoundryTransportData) *FoundryTransport {
	return &FoundryTransport{
		CurrWsId: 0,
		ExchangeChan: types.ExchangeChannels{
			Msgs:        make(map[int]chan []byte),
			ProgressMsg: map[string]chan struct{}{},
		},

		Http:             data.HttpConfig,
		ReconnectTimeout: 500 * time.Millisecond,
		IsAuth:           false,
		IsLogin:          false,
		ReconnectNumMax:  1,
		ReconnectNum:     0,

		DB:         data.DbConn,
		IsDbInit:   false,
		InitWorlds: data.Worlds,

		Logger: data.Logger,
	}
}

func (t *FoundryTransport) CloseMsgChannel(msgChan chan []byte, id int, timeout time.Duration) bool {
	time.Sleep(timeout)

	t.ChanMutex.Lock()
	defer t.ChanMutex.Unlock()

	ok := true
	select {
	case _, ok = <-msgChan:
		if ok {
			if msgChan == t.ExchangeChan.Msgs[id] {
				delete(t.ExchangeChan.Msgs, id)
			}
			close(msgChan)
		}
	default:
		if msgChan == t.ExchangeChan.Msgs[id] {
			delete(t.ExchangeChan.Msgs, id)
		}
		close(msgChan)
	}

	return ok
}
