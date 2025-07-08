package hub

import "fmt"

type DataType byte

const (
	DataFirst            DataType = 0x00
	DataNone             DataType = 0x01
	DataSingleValue      DataType = 0x02
	DataMultipleValue    DataType = 0x03
	DataKeyValuePair     DataType = 0x04
	DataMultipleKeyValue DataType = 0x04
	DataStatusCode       DataType = 0x05
)

func (d DataType) String() string {
	switch d {
	case DataFirst:
		return "First"
	case DataNone:
		return "None"
	case DataSingleValue:
		return "SingleValue"
	case DataMultipleValue:
		return "MultipleValue"
	case DataKeyValuePair:
		return "KeyValuePair"
	case DataStatusCode:
		return "StatusCode"
	default:
		return fmt.Sprintf("datatype<%02x>", int(d))
	}
}
