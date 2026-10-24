package hub

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

var (
	// ErrUnsupportedDataType is returned by MarshallCommand for a data value
	// that is not a byte, int, uint16 or uint32.
	ErrUnsupportedDataType     = errors.New("unsupported data type")
	// ErrCommandResponseNotFound is returned by UnmarshallCommand when no
	// line holds a reply.
	ErrCommandResponseNotFound = errors.New("command response not found")
	// ErrResponseTooShort is returned by UnmarshallCommand when a reply
	// holds fewer data bytes than its length field announces.
	ErrResponseTooShort        = errors.New("response shorter than its length field")
	// ErrCommandFailed is wrapped by the error returned when the hub answers
	// with a status code other than StatusOk.
	ErrCommandFailed           = errors.New("command failed")
)

func marshalMultiLenUint[T constraints.Integer](value T) string {
	switch {
	case value < any(256).(T):
		return fmt.Sprintf("%02X", value)
	case value < any(256*256).(T):
		return fmt.Sprintf("%04X", value)
	default:
		return fmt.Sprintf("%08X", value)
	}
}

// MarshallCommand builds a request line. Each data value is a byte, int,
// uint16 or uint32. A byte is written as one byte. The others are written
// as 1, 2 or 4 bytes, the fewest that hold the value.
func MarshallCommand(command Command, dataType DataType, data ...any) (string, error) {

	switch len(data) {
	case 0:
		return fmt.Sprintf(">%02X%02X%04X", int(command), int(dataType), 0), nil
	case 1:
		switch data := data[0].(type) {
		case byte:
			return fmt.Sprintf(">%02X%02X%04X%02X", int(command), int(dataType), 1, data), nil
		case int:
			hexData := marshalMultiLenUint[int](data)
			return fmt.Sprintf(">%02X%02X%04X%s", int(command), int(dataType), len(hexData)/2, hexData), nil
		case uint16:
			hexData := marshalMultiLenUint[uint16](data)
			return fmt.Sprintf(">%02X%02X%04X%s", int(command), int(dataType), len(hexData)/2, hexData), nil
		case uint32:
			hexData := marshalMultiLenUint[uint32](data)
			return fmt.Sprintf(">%02X%02X%04X%s", int(command), int(dataType), len(hexData)/2, hexData), nil
		default:
			return "", fmt.Errorf("%w '%s'", ErrUnsupportedDataType, reflect.TypeOf(data))
		}
	default:
		var hexData strings.Builder
		for _, d := range data {
			switch d := d.(type) {
			case byte:
				hexData.WriteString(fmt.Sprintf("%02X", d))
			case int:
				hexData.WriteString(marshalMultiLenUint[int](d))
			case uint16:
				hexData.WriteString(marshalMultiLenUint[uint16](d))
			case uint32:
				hexData.WriteString(marshalMultiLenUint[uint32](d))
			default:
				return "", fmt.Errorf("%w '%s'", ErrUnsupportedDataType, reflect.TypeOf(data))
			}
		}
		return fmt.Sprintf(">%02X%02X%04X", int(command), int(dataType), hexData.Len()/2) + hexData.String(), nil
	}
}

var respRegexp = regexp.MustCompile(`(?i)^<([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{4})([0-9a-f]+)$`)

// Decoded reply payloads.
type (
	// DataList is the payload of a DataSingleValue or DataMultipleValue
	// reply.
	DataList = []uint8
	// DataMap is the payload of a DataKeyValuePair reply.
	DataMap  = map[uint8]uint8
)

// UnmarshallCommand decodes the first line of resp that holds a reply. The
// type of data follows dataType: nil for DataFirst and DataNone, DataList,
// DataMap, or StatusCode for DataStatusCode.
func UnmarshallCommand(resp []string) (command Command, dataType DataType, data any, err error) {
	for _, line := range resp {
		if m := respRegexp.FindStringSubmatch(line); m != nil {
			commandValue, _ := strconv.ParseUint(m[1], 16, 32)
			dataTypeValue, _ := strconv.ParseUint(m[2], 16, 32)
			dataLenValue, _ := strconv.ParseUint(m[3], 16, 32)
			lineData := m[4]
			if uint64(len(lineData)) < dataLenValue*2 {
				return 0, 0, nil, fmt.Errorf("%w: %d data bytes announced, %d hex digits present",
					ErrResponseTooShort, dataLenValue, len(lineData))
			}

			switch DataType(dataTypeValue) {
			case DataFirst, DataNone:
				return Command(commandValue), DataType(dataTypeValue), nil, nil
			case DataSingleValue, DataMultipleValue:
				var dataList DataList
				for i := range dataLenValue {
					hexValue := lineData[i*2 : (i+1)*2]
					dataValue, _ := strconv.ParseUint(hexValue, 16, 32)
					dataList = append(dataList, uint8(dataValue))
				}
				return Command(commandValue), DataType(dataTypeValue), dataList, nil
			case DataKeyValuePair /* same id as DataMultipleKeyValue */ :
				dataMap := make(DataMap, dataLenValue/2)
				for i := range dataLenValue / 2 {
					keyHexValue := lineData[i*4 : i*4+2]
					keyValue, _ := strconv.ParseUint(keyHexValue, 16, 32)
					dataHexValue := lineData[i*4+2 : i*4+4]
					dataValue, _ := strconv.ParseUint(dataHexValue, 16, 32)
					dataMap[uint8(keyValue)] = uint8(dataValue)
				}
				return Command(commandValue), DataType(dataTypeValue), dataMap, nil
			case DataStatusCode:
				statusCodeValue, _ := strconv.ParseUint(lineData, 16, 32)
				return Command(commandValue), DataType(dataTypeValue), StatusCode(statusCodeValue), nil
			}
		}
	}
	return 0, 0, nil, ErrCommandResponseNotFound
}

func sliceToAny[T any](slice []T) []any {
	anySlice := make([]any, len(slice))
	for idx, value := range slice {
		anySlice[idx] = value
	}
	return anySlice
}

func (h *Hub) sendCommand(ctx context.Context, command Command, dataType DataType, data ...any) error {
	if len(data) == 1 {
		switch typedData := data[0].(type) {
		case []byte:
			data = sliceToAny(typedData)
		}
	}

	req, err := MarshallCommand(command, dataType, data...)
	if err != nil {
		return err
	}

	resp, err := h.conn.SendAndRecv(ctx, req)
	if err != nil {
		fmt.Printf("recv error: %v\n", err)
		return err
	}

	_, respDataType, respData, err := UnmarshallCommand(resp)
	if err != nil {
		fmt.Printf("failed to unmarshal resp '%s': %s\n", resp, err)
		return err
	}

	if respDataType == DataStatusCode {
		if status := respData.(StatusCode); status != StatusOk {
			return fmt.Errorf("%w: %s %s", ErrCommandFailed, command, status.String())
		}
	}

	return nil
}

func (h *Hub) sendCommandAndRecv(ctx context.Context, command Command, dataType DataType, data ...any) (any, error) {
	req, err := MarshallCommand(command, dataType, data...)
	if err != nil {
		return nil, err
	}

	resp, err := h.conn.SendAndRecv(ctx, req)
	if err != nil {
		fmt.Printf("recv error: %v\n", err)
		return nil, err
	}

	_, respDataType, respData, err := UnmarshallCommand(resp)
	if err != nil {
		fmt.Printf("failed to unmarshal resp '%s': %s\n", resp, err)
		return nil, err
	}

	if respDataType == DataStatusCode {
		if status := respData.(StatusCode); status != StatusOk {
			return nil, fmt.Errorf("%w: %s", ErrCommandFailed, status.String())
		}
	}

	// fmt.Printf("recv: command %s dataType %s %s\n", respCommand, respDataType, spew.Sdump(respData))

	return respData, nil
}
