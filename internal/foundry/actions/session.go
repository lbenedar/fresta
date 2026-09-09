package actions

import (
	"time"

	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

type WsSessionMsg struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId,omitempty"`
}

func (msg WsSessionMsg) Action(tr *transport.FoundryTransport, status *types.FoundryStatus) error {
	if status.IsActive {
		return msg.OnActiveWorld(tr)
	}

	if tr.IsDbInit {
		return nil
	}
	tr.IsDbInit = true

	return tr.FillDBWithFoundryData()
}

func (msg WsSessionMsg) OnActiveWorld(tr *transport.FoundryTransport) error {
	tr.Logger.Info("Session msg", "userid", msg.UserId)
	if msg.UserId != "" {
		return msg.HandleLoggedInUser(tr)
	}

	userId, userPass, err := tr.GetUserIdAndPass()
	if err != nil {
		return err
	}

	err = tr.LogInToWorld(userId, userPass)
	if err != nil {
		return err
	}

	tr.Logger.Info("World is started. Succesfully logged into world")
	close(tr.ReadChan.Reconnect())
	return nil
}

func (msg WsSessionMsg) HandleLoggedInUser(tr *transport.FoundryTransport) error {
	if tr.LoggedInChan == nil {
		tr.Logger.Info("World had been started before application was started. Run insertion of world data")
		return tr.InsertGameToDB()
	}

	select {
	case <-tr.LoggedInChan:
		return ErrUserChannelIsClosed
	default:
		tr.LoggedInChan <- true
	}

	time.Sleep(3 * time.Second)
	types.CloseChannel(tr.LoggedInChan)
	return nil
}
