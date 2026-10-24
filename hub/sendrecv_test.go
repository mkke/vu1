package hub

import (
	"errors"
	"maps"
	"testing"
)

func TestUnmarshallCommandKeyValuePairs(t *testing.T) {
	// Two pairs: dial 0 at 0x32, dial 1 at 0x64.
	_, dataType, data, err := UnmarshallCommand([]string{"<040400040032" + "0164"})
	if err != nil {
		t.Fatal(err)
	}
	if dataType != DataKeyValuePair {
		t.Fatalf("data type = %s, want %s", dataType, DataKeyValuePair)
	}
	want := DataMap{0: 0x32, 1: 0x64}
	if got := data.(DataMap); !maps.Equal(got, want) {
		t.Errorf("data = %v, want %v", got, want)
	}
}

func TestUnmarshallCommandShortData(t *testing.T) {
	// The length field announces 4 bytes, the line carries 1.
	_, _, _, err := UnmarshallCommand([]string{"<0302000400"})
	if !errors.Is(err, ErrResponseTooShort) {
		t.Errorf("err = %v, want %v", err, ErrResponseTooShort)
	}
}

func TestUnmarshallCommandSkipsLogLines(t *testing.T) {
	command, _, data, err := UnmarshallCommand([]string{"boot ok", "<0C0500020000"})
	if err != nil {
		t.Fatal(err)
	}
	if command != CmdRescanBus || data != StatusOk {
		t.Errorf("got %s %v, want %s %v", command, data, CmdRescanBus, StatusOk)
	}
}
