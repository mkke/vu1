package hub_test

import (
	"fmt"

	"github.com/mkke/vu1/hub"
)

func ExampleMarshallCommand() {
	// Move the needle of dial 0 to 50 percent.
	req, err := hub.MarshallCommand(hub.CmdSetDialPercSingle, hub.DataKeyValuePair, uint8(0), uint8(50))
	if err != nil {
		panic(err)
	}
	fmt.Println(req)
	// Output: >030400020032
}

func ExampleUnmarshallCommand() {
	command, dataType, data, err := hub.UnmarshallCommand([]string{"<0C0500020000"})
	if err != nil {
		panic(err)
	}
	fmt.Println(command, dataType, data)
	// Output: RescanBus StatusCode OK
}
