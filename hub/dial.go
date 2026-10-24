package hub

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
)

// Dial sends commands to one dial attached to a Hub.
type Dial struct {
	hub       *Hub
	dialIndex uint8
}

// NewDial returns the Dial at dialIndex, an index that
// DevicesMap.OnlineDials reports. It does not check that the dial exists.
func (h *Hub) NewDial(dialIndex uint8) *Dial {
	return &Dial{
		hub:       h,
		dialIndex: dialIndex,
	}
}

// GetUid returns the dial's unique id.
func (d *Dial) GetUid(ctx context.Context) ([]byte, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetDeviceUid, DataSingleValue, d.dialIndex); err != nil {
		return nil, err
	} else {
		return data.([]byte), nil
	}
}

// GetFirmwareHash returns the build hash of the dial's firmware.
func (d *Dial) GetFirmwareHash(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetBuildInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

// GetFirmwareVersion returns the version of the dial's firmware.
func (d *Dial) GetFirmwareVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetFwInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

// GetHardwareVersion returns the dial's hardware revision.
func (d *Dial) GetHardwareVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetHwInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

// GetProtocolVersion returns the protocol version the dial speaks.
func (d *Dial) GetProtocolVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetProtocolInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

// CalibratePosition names the needle position a calibration value is for.
type CalibratePosition int

// Calibration positions: the middle of the scale and its end.
const (
	CalibrateHalf CalibratePosition = iota
	CalibrateFull
)

// Calibrate stores value as the raw needle drive value for position.
func (d *Dial) Calibrate(ctx context.Context, position CalibratePosition, value uint32) error {
	// we have to make sure full 4 bytes are sent
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, value)

	switch position {
	case CalibrateHalf:
		return d.hub.sendCommand(ctx, CmdSetDialCalibrateHalf, DataKeyValuePair, data)
	case CalibrateFull:
		return d.hub.sendCommand(ctx, CmdSetDialCalibrateMax, DataKeyValuePair, data)
	default:
		return errors.New("invalid calibration position")
	}
}

// SetDialEasingStep sets how far the needle moves per easing period on its
// way to a new position.
func (d *Dial) SetDialEasingStep(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetDialEasingStep, DataSingleValue, data)
}

// SetDialEasingPeriod sets the interval between needle easing steps.
func (d *Dial) SetDialEasingPeriod(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetDialEasingPeriod, DataSingleValue, data)
}

// SetBacklightEasingStep sets how far the backlight moves per easing period
// on its way to a new color.
func (d *Dial) SetBacklightEasingStep(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetBacklightEasingStep, DataSingleValue, data)
}

// SetBacklightEasingPeriod sets the interval between backlight easing steps.
func (d *Dial) SetBacklightEasingPeriod(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetBacklightEasingPeriod, DataSingleValue, data)
}

// EasingConfig holds the values the Set*Easing* methods set.
type EasingConfig struct {
	DialStep        uint32
	DialPeriod      uint32
	BacklightStep   uint32
	BacklightPeriod uint32
}

// GetEasingConfig returns the dial's current easing settings.
func (d *Dial) GetEasingConfig(ctx context.Context) (*EasingConfig, error) {
	resp, err := d.hub.sendCommandAndRecv(ctx, CmdGetEasingConfig, DataSingleValue, d.dialIndex)
	if err != nil {
		return nil, err
	}

	data := resp.([]byte)
	if len(data) < 16 {
		return nil, fmt.Errorf("expected %d easing config bytes, but got %d", 16, len(data))
	}

	return &EasingConfig{
		DialStep:        binary.BigEndian.Uint32(data[0:4]),
		DialPeriod:      binary.BigEndian.Uint32(data[4:8]),
		BacklightStep:   binary.BigEndian.Uint32(data[8:12]),
		BacklightPeriod: binary.BigEndian.Uint32(data[12:16]),
	}, nil
}

// SetRaw moves the needle to a raw drive value, sent as 16 bits.
func (d *Dial) SetRaw(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint16([]byte{d.dialIndex}, uint16(value))
	return d.hub.sendCommand(ctx, CmdSetDialRawSingle, DataKeyValuePair, data)
}

// SetPercent moves the needle to value, from 0 to 100.
func (d *Dial) SetPercent(ctx context.Context, value uint8) error {
	return d.hub.sendCommand(ctx, CmdSetDialPercSingle, DataKeyValuePair, d.dialIndex, value)
}

// DisplayClear fills the display white, or black when whiteBackground is
// false.
func (d *Dial) DisplayClear(ctx context.Context, whiteBackground bool) error {
	value := 0
	if !whiteBackground {
		value = 1
	}
	return d.hub.sendCommand(ctx, CmdDisplayClear, DataSingleValue, d.dialIndex, value)
}

// DisplayGotoXY sets the display position the next DisplaySendImageData
// writes to.
func (d *Dial) DisplayGotoXY(ctx context.Context, x, y int) error {
	data := []byte{d.dialIndex}
	data = binary.BigEndian.AppendUint16(data, uint16(x))
	data = binary.BigEndian.AppendUint16(data, uint16(y))
	return d.hub.sendCommand(ctx, CmdDisplayGotoXY, DataSingleValue, data)
}

// DisplaySendImageData converts image with ConvertImageData and uploads it.
// The display does not change until DisplayShow.
func (d *Dial) DisplaySendImageData(ctx context.Context, image image.Image) error {
	imageData, err := ConvertImageData(image)
	if err != nil {
		return err
	}

	const chunkSize = 1000
	for idx := 0; idx < len(imageData); idx += chunkSize {
		chunk := imageData[idx:min(idx+chunkSize, len(imageData))]
		if err := d.hub.sendCommand(ctx, CmdDisplayImgData, DataSingleValue, append([]byte{d.dialIndex}, chunk...)); err != nil {
			return err
		}
	}
	return nil
}

// GetRXBufferSize returns the size of the dial's receive buffer in bytes.
func (d *Dial) GetRXBufferSize(ctx context.Context) (int, error) {
	resp, err := d.hub.sendCommandAndRecv(ctx, CmdRxBufferSize, DataSingleValue, d.dialIndex)
	if err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint32(resp.([]byte))), nil
}

// DisplayShow puts the uploaded image on the display.
func (d *Dial) DisplayShow(ctx context.Context) error {
	return d.hub.sendCommand(ctx, CmdDisplayShowImg, DataSingleValue, d.dialIndex)
}

// SetBacklight sets the backlight's red, green, blue and white channels,
// each from 0 to 255.
func (d *Dial) SetBacklight(ctx context.Context, red, green, blue, white uint8) error {
	return d.hub.sendCommand(ctx, CmdSetRGBBacklight, DataMultipleValue, d.dialIndex, red, green, blue, white)
}
