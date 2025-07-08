package hub

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
)

type Dial struct {
	hub       *Hub
	dialIndex uint8
}

func (h *Hub) NewDial(dialIndex uint8) *Dial {
	return &Dial{
		hub:       h,
		dialIndex: dialIndex,
	}
}

func (d *Dial) GetUid(ctx context.Context) ([]byte, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetDeviceUid, DataSingleValue, d.dialIndex); err != nil {
		return nil, err
	} else {
		return data.([]byte), nil
	}
}

func (d *Dial) GetFirmwareHash(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetBuildInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

func (d *Dial) GetFirmwareVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetFwInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

func (d *Dial) GetHardwareVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetHwInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

func (d *Dial) GetProtocolVersion(ctx context.Context) (string, error) {
	if data, err := d.hub.sendCommandAndRecv(ctx, CmdGetProtocolInfo, DataSingleValue, d.dialIndex); err != nil {
		return "", err
	} else {
		return string(data.([]byte)), nil
	}
}

type CalibratePosition int

const (
	CalibrateHalf CalibratePosition = iota
	CalibrateFull
)

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

func (d *Dial) SetDialEasingStep(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetDialEasingStep, DataSingleValue, data)
}

func (d *Dial) SetDialEasingPeriod(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetDialEasingPeriod, DataSingleValue, data)
}

func (d *Dial) SetBacklightEasingStep(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetBacklightEasingStep, DataSingleValue, data)
}

func (d *Dial) SetBacklightEasingPeriod(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint32([]byte{d.dialIndex}, uint32(value))
	return d.hub.sendCommand(ctx, CmdSetBacklightEasingPeriod, DataSingleValue, data)
}

type EasingConfig struct {
	DialStep        uint32
	DialPeriod      uint32
	BacklightStep   uint32
	BacklightPeriod uint32
}

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

func (d *Dial) SetRaw(ctx context.Context, value int) error {
	data := binary.BigEndian.AppendUint16([]byte{d.dialIndex}, uint16(value))
	return d.hub.sendCommand(ctx, CmdSetDialRawSingle, DataKeyValuePair, data)
}

func (d *Dial) SetPercent(ctx context.Context, value uint8) error {
	return d.hub.sendCommand(ctx, CmdSetDialPercSingle, DataKeyValuePair, d.dialIndex, value)
}

func (d *Dial) DisplayClear(ctx context.Context, whiteBackground bool) error {
	value := 0
	if !whiteBackground {
		value = 1
	}
	return d.hub.sendCommand(ctx, CmdDisplayClear, DataSingleValue, d.dialIndex, value)
}

func (d *Dial) DisplayGotoXY(ctx context.Context, x, y int) error {
	data := []byte{d.dialIndex}
	data = binary.BigEndian.AppendUint16(data, uint16(x))
	data = binary.BigEndian.AppendUint16(data, uint16(y))
	return d.hub.sendCommand(ctx, CmdDisplayGotoXY, DataSingleValue, data)
}

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

func (d *Dial) GetRXBufferSize(ctx context.Context) (int, error) {
	resp, err := d.hub.sendCommandAndRecv(ctx, CmdRxBufferSize, DataSingleValue, d.dialIndex)
	if err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint32(resp.([]byte))), nil
}

func (d *Dial) DisplayShow(ctx context.Context) error {
	return d.hub.sendCommand(ctx, CmdDisplayShowImg, DataSingleValue, d.dialIndex)
}

func (d *Dial) SetBacklight(ctx context.Context, red, green, blue, white uint8) error {
	return d.hub.sendCommand(ctx, CmdSetRGBBacklight, DataMultipleValue, d.dialIndex, red, green, blue, white)
}
