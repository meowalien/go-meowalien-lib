package websocket


const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)


//RFC_6455
func WebsocketCloseCodeNumberToString(errorCode int) string {
	switch errorCode {
	case CloseNormalClosure: //1000
		return "CloseNormalClosure"
	case CloseGoingAway: //1001
		return "CloseGoingAway"
	case CloseProtocolError: //1002
		return "CloseProtocolError"
	case CloseUnsupportedData: //1003
		return "CloseUnsupportedData"
	case CloseNoStatusReceived: //1005
		return "CloseNoStatusReceived"
	case CloseAbnormalClosure: //1006
		return "CloseAbnormalClosure"
	case CloseInvalidFramePayloadData: //1007
		return "CloseInvalidFramePayloadData"
	case ClosePolicyViolation: //1008
		return "ClosePolicyViolation"
	case CloseMessageTooBig: //1009
		return "CloseMessageTooBig"
	case CloseMandatoryExtension: //1010
		return "CloseMandatoryExtension"
	case CloseInternalServerErr: //1011
		return "CloseInternalServerErr"
	case CloseServiceRestart: //1012
		return "CloseServiceRestart"
	case CloseTryAgainLater: //1013
		return "CloseTryAgainLater"
	case CloseTLSHandshake: //1015
		return "CloseTLSHandshake"
	default:
		return "Unknown"
	}
}
