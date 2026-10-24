package hub

import "fmt"

// Command is the command byte of a request or reply.
type Command byte

// Commands understood by the hub. The DialBTL and HubBTL commands drive the
// bootloaders used for firmware updates.
const (
	CmdSetDialRawSingle         Command = 0x01
	CmdSetDialRawMultiple       Command = 0x02
	CmdSetDialPercSingle        Command = 0x03
	CmdSetDialPercMultiple      Command = 0x04
	CmdSetDialCalibrateMax      Command = 0x05
	CmdSetDialCalibrateHalf     Command = 0x06
	CmdGetDevicesMap            Command = 0x07
	CmdProvisionDevice          Command = 0x08
	CmdResetAllDevices          Command = 0x09
	CmdDialPower                Command = 0x0A
	CmdGetDeviceUid             Command = 0x0B
	CmdRescanBus                Command = 0x0C
	CmdDisplayClear             Command = 0x0D
	CmdDisplayGotoXY            Command = 0x0E
	CmdDisplayImgData           Command = 0x0F
	CmdDisplayShowImg           Command = 0x10
	CmdRxBufferSize             Command = 0x11
	CmdResetCfg                 Command = 0x12
	CmdSetRGBBacklight          Command = 0x13
	CmdSetDialEasingStep        Command = 0x14
	CmdSetDialEasingPeriod      Command = 0x15
	CmdSetBacklightEasingStep   Command = 0x16
	CmdSetBacklightEasingPeriod Command = 0x17
	CmdGetEasingConfig          Command = 0x18
	CmdGetBuildInfo             Command = 0x19
	CmdGetFwInfo                Command = 0x20
	CmdGetHwInfo                Command = 0x21
	CmdGetProtocolInfo          Command = 0x22
	CmdDebugI2CScan             Command = 0xF3
	CmdHubBTLJumpToBootloader   Command = 0xF4
	CmdDialBTLJumpToBootloader  Command = 0xF5
	CmdDialBTLGetInfo           Command = 0xF6
	CmdDialBTLGetCRC            Command = 0xF7
	CmdDialBTLEraseApp          Command = 0xF8
	CmdDialBTLFwupSendPackage   Command = 0xF9
	CmdDialBTLFwupFinished      Command = 0xFA
	CmdDialBTLExit              Command = 0xFB
	CmdDialBTLRestartFwUpload   Command = 0xFC
	CmdDialBTLReadLastStatus    Command = 0xFD
)

func (c Command) String() string {
	switch c {
	case CmdSetDialRawSingle:
		return "SetDialRawSingle"
	case CmdSetDialRawMultiple:
		return "SetDialRawMultiple"
	case CmdSetDialPercSingle:
		return "SetDialPercSingle"
	case CmdSetDialPercMultiple:
		return "SetDialPercMultiple"
	case CmdSetDialCalibrateMax:
		return "SetDialCalibrateMax"
	case CmdSetDialCalibrateHalf:
		return "SetDialCalibrateHalf"
	case CmdGetDevicesMap:
		return "GetDevicesMap"
	case CmdProvisionDevice:
		return "ProvisionDevice"
	case CmdResetAllDevices:
		return "ResetAllDevices"
	case CmdDialPower:
		return "DialPower"
	case CmdGetDeviceUid:
		return "GetDeviceUid"
	case CmdRescanBus:
		return "RescanBus"
	case CmdDisplayClear:
		return "DisplayClear"
	case CmdDisplayGotoXY:
		return "DisplayGotoXY"
	case CmdDisplayImgData:
		return "DisplayImgData"
	case CmdDisplayShowImg:
		return "DisplayShowImg"
	case CmdRxBufferSize:
		return "RxBufferSize"
	case CmdResetCfg:
		return "ResetCfg"
	case CmdSetRGBBacklight:
		return "SetRGBBacklight"
	case CmdSetDialEasingStep:
		return "SetDialEasingStep"
	case CmdSetDialEasingPeriod:
		return "SetDialEasingPeriod"
	case CmdSetBacklightEasingStep:
		return "SetBacklightEasingStep"
	case CmdSetBacklightEasingPeriod:
		return "SetBacklightEasingPeriod"
	case CmdGetEasingConfig:
		return "GetEasingConfig"
	case CmdGetBuildInfo:
		return "GetBuildInfo"
	case CmdGetFwInfo:
		return "GetFwInfo"
	case CmdGetHwInfo:
		return "GetHwInfo"
	case CmdGetProtocolInfo:
		return "GetProtocolInfo"
	case CmdDebugI2CScan:
		return "DebugI2CScan"
	case CmdHubBTLJumpToBootloader:
		return "HubBTLJumpToBootloader"
	case CmdDialBTLJumpToBootloader:
		return "DialBTLJumpToBootloader"
	case CmdDialBTLGetInfo:
		return "DialBTLGetInfo"
	case CmdDialBTLGetCRC:
		return "DialBTLGetCRC"
	case CmdDialBTLEraseApp:
		return "DialBTLEraseApp"
	case CmdDialBTLFwupSendPackage:
		return "DialBTLFwupSendPackage"
	case CmdDialBTLFwupFinished:
		return "DialBTLFwupFinished"
	case CmdDialBTLExit:
		return "DialBTLExit"
	case CmdDialBTLRestartFwUpload:
		return "DialBTLRestartFwUpload"
	case CmdDialBTLReadLastStatus:
		return "DialBTLReadLastStatus"
	default:
		return fmt.Sprintf("command<%02x>", byte(c))
	}
}
