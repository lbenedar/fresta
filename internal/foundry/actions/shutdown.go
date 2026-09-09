package actions

import (
	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

type WsShutdownMsg struct {
	World  string  `json:"world"`
	UserId *string `json:"userId,omitempty"`
}

func (msg WsShutdownMsg) Action(tr *transport.FoundryTransport, status *types.FoundryStatus) error {
	if status.IsActive != false {
		newStatus, err := tr.Http.GetStatus()
		if err != nil {
			return err
		}
		status.Update(newStatus)
		tr.Logger.Info("Status has been changed", "status", status)
	}

	if tr.IsDbInit {
		return nil
	}
	tr.IsDbInit = true

	return tr.FillDBWithFoundryData()
}
