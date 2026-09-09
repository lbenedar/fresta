package types

type DirectionCode int
type TransportCode int

const (
	WriterCode = DirectionCode(0)
	ReaderCode = DirectionCode(1)

	DbCode        = TransportCode(0)
	HttpCode      = TransportCode(1)
	WebSocketCode = TransportCode(2)
)

const (
	RespSessionDataCode  = "0"
	RespPingCode         = "2"
	RespSessionIdCode    = "40"
	RespServerChangeCode = "42"
	RespDataCode         = "43"

	ReqPongCode          = "3"
	ReqCreateSessionCode = "40"
	ReqDataCode          = "42"
)

var RequestCodes = []string{
	RespSessionDataCode,
	RespPingCode,
	RespSessionIdCode,
	RespServerChangeCode,
	RespDataCode,
}

var CodesRespToReq = map[string]string{
	RespSessionDataCode:  ReqCreateSessionCode,
	RespPingCode:         ReqPongCode,
	RespServerChangeCode: ReqDataCode,
	//ReqDataCode:  RespServerChangeCode,
}
