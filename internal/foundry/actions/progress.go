package actions

import (
	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

type WsProgressMsg struct {
	Id         string  `json:"id"`
	Message    string  `json:"message"`
	Pct        float64 `json:"pct"`
	HasChanged bool    `json:"hasChanged,omitempty"`
	Act        string  `json:"action"`
	Step       string  `json:"step"`
}

func (msg WsProgressMsg) Action(tr *transport.FoundryTransport, status *types.FoundryStatus) error {
	if !msg.HasChanged {
		tr.Logger.Debug("World has been launched", "world", msg.Id)
		select {
		case <-tr.ExchangeChan.ProgressMsg[msg.Id]:
			tr.Logger.Warn("Channel for the world is closed", "world", msg.Id)
		default:
			close(tr.ExchangeChan.ProgressMsg[msg.Id])
		}

		if status.IsActive != true {
			newStatus, err := tr.Http.GetStatus()
			if err != nil {
				return err
			}
			status.Update(newStatus)
			tr.Logger.Info("Status has been changed", "status", status)
		}
	}
	return nil
}
