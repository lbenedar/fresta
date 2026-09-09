package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lbenedar/fresta/internal/foundry/requests"
)

var (
	ErrorMsgNotHaveNumber = errors.New("Message doesn't have dataCode and id")
)

type WsMessage struct {
	Code    string
	Id      int
	MsgJson string
}

func (w WsMessage) ToString() string {
	return fmt.Sprintf("%s%d%s", w.Code, w.Id, w.MsgJson)
}

func (w WsMessage) ToByteSlice() []byte {
	return fmt.Appendf([]byte{}, "%s%d%s", w.Code, w.Id, w.MsgJson)
}

func parseCode(msg *string, requestCodes []string) string {
	for j := range requestCodes {
		if !strings.HasPrefix(*msg, requestCodes[j]) {
			continue
		}
		return requestCodes[j]
	}
	return ""
}

func parseId(msg *string, start int) (int, int) {
	msgLen := len(*msg)
	j := start
	for ; j < msgLen; j++ {
		if (*msg)[j] < '0' || (*msg)[j] > '9' {
			break
		}
	}
	if j >= msgLen || j <= 0 {
		return 0, 0
	}

	msgId, err := strconv.Atoi((*msg)[start:j])
	if err != nil {
		return 0, 0
	}
	return msgId, j
}

func ParseWsRespMessage(msg string, requestCodes []string) (*WsMessage, error) {
	data := &WsMessage{}

	data.Code = parseCode(&msg, requestCodes)

	i := len(data.Code)
	if i == 0 {
		return nil, ErrorMsgNotHaveNumber
	}

	data.Id, i = parseId(&msg, i)
	if i == 0 {
		i = len(data.Code)
	}
	// if i == 0 {
	// 	return nil, ErrorMsgNotHaveNumber
	// }

	data.MsgJson = msg[i:]
	return data, nil
}

func NewWsMessageByPage(page string, currWsId int) *WsMessage {
	return &WsMessage{Code: CodesRespToReq[RespServerChangeCode], Id: currWsId, MsgJson: fmt.Sprintf("[\"%s\"]", requests.WsTypeData[page])}
}

func NewWsMessage(msg string, currWsId int) *WsMessage {
	return &WsMessage{Code: CodesRespToReq[RespServerChangeCode], Id: currWsId, MsgJson: fmt.Sprintf("[\"%s\"]", msg)}
}
