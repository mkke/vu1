package hub

import "fmt"

// StatusCode represents a hub status code
type StatusCode uint16

const (
	StatusOk                StatusCode = 0x0000
	StatusFail              StatusCode = 0x0001
	StatusBusy              StatusCode = 0x0002
	StatusTimeout           StatusCode = 0x0003
	StatusBadData           StatusCode = 0x0004
	StatusProtocolError     StatusCode = 0x0005
	StatusNoMemory          StatusCode = 0x0006
	StatusInvalidArgument   StatusCode = 0x0007
	StatusBadAddress        StatusCode = 0x0008
	StatusForbidden         StatusCode = 0x0009
	StatusAlreadyExists     StatusCode = 0x000B
	StatusUnsupported       StatusCode = 0x000C
	StatusNotImplemented    StatusCode = 0x000D
	StatusMalformedPackage  StatusCode = 0x000E
	StatusRecursiveCall     StatusCode = 0x0010
	StatusDataMismatch      StatusCode = 0x0011
	StatusDeviceOffline     StatusCode = 0x0012
	StatusModuleNotInit     StatusCode = 0x0013
	StatusI2CError          StatusCode = 0x0014
	StatusUSASTError        StatusCode = 0x0015
	StatusSPIError          StatusCode = 0x0016
	StatusBTLNoDevice       StatusCode = 0xE001
	StatusBTLInvalidState   StatusCode = 0xE002
	StatusBTLInvalidRequest StatusCode = 0xE003
)

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return "OK"
	case StatusFail:
		return "Fail"
	case StatusBusy:
		return "Busy"
	case StatusTimeout:
		return "Timeout"
	case StatusBadData:
		return "Bad Data"
	case StatusProtocolError:
		return "Protocol Error"
	case StatusNoMemory:
		return "No Memory"
	case StatusInvalidArgument:
		return "Invalid Argument"
	case StatusBadAddress:
		return "Bad Address"
	case StatusForbidden:
		return "Forbidden"
	case StatusAlreadyExists:
		return "Already Exists"
	case StatusUnsupported:
		return "Unsupported"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusMalformedPackage:
		return "Malformed Package"
	case StatusRecursiveCall:
		return "Recursive Call"
	case StatusDataMismatch:
		return "Data Mismatch"
	case StatusDeviceOffline:
		return "Device Offline"
	case StatusModuleNotInit:
		return "Module Not Init"
	case StatusI2CError:
		return "I2C Error"
	case StatusUSASTError:
		return "USART Error"
	case StatusSPIError:
		return "SPI Error"
	case StatusBTLNoDevice:
		return "BTL No Device"
	case StatusBTLInvalidState:
		return "BTL Invalid State"
	case StatusBTLInvalidRequest:
		return "BTL Invalid Request"
	default:
		return fmt.Sprintf("status<%02x>", int(s))
	}
}
