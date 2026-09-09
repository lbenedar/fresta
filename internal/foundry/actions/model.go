package actions

import (
	"encoding/json"

	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
)

const (
	WsSessionType      = "session"
	WsProgressType     = "progress"
	WsUserActivityType = "userActivity"
	WsShutdownType     = "shutdown"
)

type WsActions interface {
	Action(tr *transport.FoundryTransport, status *types.FoundryStatus) error
}

var WsMsgActions = map[string]func() WsActions{
	"session":  func() WsActions { return &WsSessionMsg{} },
	"progress": func() WsActions { return &WsProgressMsg{} },
	// "userActivity": WsUserActivityCursor{},
	"shutdown": func() WsActions { return &WsShutdownMsg{} },
}

func GetWsMsgAction(dataType string) WsActions {
	if factory, ok := WsMsgActions[dataType]; ok {
		return factory()
	}
	return nil
}

func SplitToTypeAndData(dataJson []byte) (WsActions, error) {
	var rawItems []json.RawMessage
	err := json.Unmarshal(dataJson, &rawItems)
	if err != nil {
		return nil, err
	}

	var dataType string
	err = json.Unmarshal(rawItems[0], &dataType)
	if err != nil {
		return nil, err
	}

	if dataType == WsUserActivityType {
		return nil, nil
	}

	var data map[string]any
	err = json.Unmarshal(rawItems[1], &data)
	if err != nil {
		return nil, err
	}

	actionDataPtr := GetWsMsgAction(dataType)
	if actionDataPtr == nil {
		return nil, ErrActionNotFound
	}

	err = FillStruct(data, actionDataPtr)
	if err != nil {
		return nil, err
	}

	return actionDataPtr, nil
}
