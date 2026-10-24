package hub

import (
	"context"

	"github.com/mkke/vu1/conn"
)

// Hub sends commands to a VU-1 hub.
type Hub struct {
	conn conn.Conn
}

// New returns a Hub that talks over conn. The Hub does not close conn.
func New(conn conn.Conn) *Hub {
	return &Hub{
		conn: conn,
	}
}

// RescanBus makes the hub search its bus for attached dials.
func (h *Hub) RescanBus(ctx context.Context) error {
	return h.sendCommand(ctx, CmdRescanBus, DataNone)
}

// DevicesMap holds one entry per dial index. A 1 marks a dial that is
// online.
type DevicesMap DataList

// OnlineDials returns the indexes of the dials that are online.
func (dm DevicesMap) OnlineDials() map[uint8]struct{} {
	onlineDials := make(map[uint8]struct{})
	for idx, v := range dm {
		if v == 1 {
			onlineDials[uint8(idx)] = struct{}{}
		}
	}
	return onlineDials
}

// GetDevicesMap returns which dials the hub found at its last bus scan.
func (h *Hub) GetDevicesMap(ctx context.Context) (DevicesMap, error) {
	data, err := h.sendCommandAndRecv(ctx, CmdGetDevicesMap, DataNone)
	if err != nil {
		return nil, err
	}

	return data.(DataList), nil
}

// SetDialPower switches the power of the dials on or off.
func (h *Hub) SetDialPower(ctx context.Context, powerOn bool) error {
	value := 0
	if powerOn {
		value = 1
	}
	return h.sendCommand(ctx, CmdDialPower, DataSingleValue, value)
}

// ProvisionDials sends CmdProvisionDevice, which prepares the attached
// dials. Call it before RescanBus.
func (h *Hub) ProvisionDials(ctx context.Context) error {
	return h.sendCommand(ctx, CmdProvisionDevice, DataNone)
}

// ResetAllDevices resets every dial on the bus.
func (h *Hub) ResetAllDevices(ctx context.Context) error {
	return h.sendCommand(ctx, CmdResetAllDevices, DataNone)
}

// DebugI2CScan makes the hub scan its I2C bus and returns the decoded reply
// payload.
func (h *Hub) DebugI2CScan(ctx context.Context) (any, error) {
	data, err := h.sendCommandAndRecv(ctx, CmdDebugI2CScan, DataNone)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// DialMultipleSetPercent moves the needles of several dials in one command.
// dialIndexPercent maps a dial index to a needle position from 0 to 100.
func (h *Hub) DialMultipleSetPercent(ctx context.Context, dialIndexPercent map[uint8]uint8) error {
	var data []any
	for dialIndex, percent := range dialIndexPercent {
		data = append(data, dialIndex, percent)
	}
	return h.sendCommand(ctx, CmdSetDialPercMultiple, DataKeyValuePair, data...)
}
