package hub

import (
	"context"

	"github.com/mkke/vu1/conn"
)

type Hub struct {
	conn conn.Conn
}

func New(conn conn.Conn) *Hub {
	return &Hub{
		conn: conn,
	}
}

func (h *Hub) RescanBus(ctx context.Context) error {
	return h.sendCommand(ctx, CmdRescanBus, DataNone)
}

type DevicesMap DataList

func (dm DevicesMap) OnlineDials() map[uint8]struct{} {
	onlineDials := make(map[uint8]struct{})
	for idx, v := range dm {
		if v == 1 {
			onlineDials[uint8(idx)] = struct{}{}
		}
	}
	return onlineDials
}

func (h *Hub) GetDevicesMap(ctx context.Context) (DevicesMap, error) {
	data, err := h.sendCommandAndRecv(ctx, CmdGetDevicesMap, DataNone)
	if err != nil {
		return nil, err
	}

	return data.(DataList), nil
}

func (h *Hub) SetDialPower(ctx context.Context, powerOn bool) error {
	value := 0
	if powerOn {
		value = 1
	}
	return h.sendCommand(ctx, CmdDialPower, DataSingleValue, value)
}

func (h *Hub) ProvisionDials(ctx context.Context) error {
	return h.sendCommand(ctx, CmdProvisionDevice, DataNone)
}

func (h *Hub) ResetAllDevices(ctx context.Context) error {
	return h.sendCommand(ctx, CmdResetAllDevices, DataNone)
}

func (h *Hub) DebugI2CScan(ctx context.Context) (any, error) {
	data, err := h.sendCommandAndRecv(ctx, CmdDebugI2CScan, DataNone)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *Hub) DialMultipleSetPercent(ctx context.Context, dialIndexPercent map[uint8]uint8) error {
	var data []any
	for dialIndex, percent := range dialIndexPercent {
		data = append(data, dialIndex, percent)
	}
	return h.sendCommand(ctx, CmdSetDialPercMultiple, DataKeyValuePair, data...)
}
