package hub

import (
	"strings"
	"testing"
)

func TestCommandStringNamesEveryCommand(t *testing.T) {
	commands := []Command{
		CmdSetDialRawSingle, CmdSetDialRawMultiple, CmdSetDialPercSingle,
		CmdSetDialPercMultiple, CmdSetDialCalibrateMax, CmdSetDialCalibrateHalf,
		CmdGetDevicesMap, CmdProvisionDevice, CmdResetAllDevices, CmdDialPower,
		CmdGetDeviceUid, CmdRescanBus, CmdDisplayClear, CmdDisplayGotoXY,
		CmdDisplayImgData, CmdDisplayShowImg, CmdRxBufferSize, CmdResetCfg,
		CmdSetRGBBacklight, CmdSetDialEasingStep, CmdSetDialEasingPeriod,
		CmdSetBacklightEasingStep, CmdSetBacklightEasingPeriod,
		CmdGetEasingConfig, CmdGetBuildInfo, CmdGetFwInfo, CmdGetHwInfo,
		CmdGetProtocolInfo, CmdDebugI2CScan, CmdHubBTLJumpToBootloader,
		CmdDialBTLJumpToBootloader, CmdDialBTLGetInfo, CmdDialBTLGetCRC,
		CmdDialBTLEraseApp, CmdDialBTLFwupSendPackage, CmdDialBTLFwupFinished,
		CmdDialBTLExit, CmdDialBTLRestartFwUpload, CmdDialBTLReadLastStatus,
	}
	for _, c := range commands {
		if name := c.String(); strings.HasPrefix(name, "command<") {
			t.Errorf("command %02x has no name: %s", byte(c), name)
		}
	}
}

func TestCommandStringUnknown(t *testing.T) {
	if got, want := Command(0xEE).String(), "command<ee>"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
