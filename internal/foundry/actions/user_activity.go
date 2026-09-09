package actions

import (
	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

type WsUserActivityMsg struct {
	Cursor WsUserActivityCursor `json:"cursor"`
}

type WsUserActivityCursor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (msg WsUserActivityMsg) Action(tr *transport.FoundryTransport, status *types.FoundryStatus) error {
	tr.Logger.Debug("WsUserActivityMsg")
	return nil
}
