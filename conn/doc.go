// Package conn carries the line-based request/response exchange with a VU-1
// hub over a serial port.
//
// A Conn sends one request line and, with SendAndRecv, collects the lines
// the hub sends back until one starts with '<'. That line holds the reply.
// Any lines read before it are returned with it. Empty lines are dropped.
// Package hub builds the request lines and decodes the replies.
//
// SerialConn wraps a port opened with go.bug.st/serial. The hub expects
// 115200 baud, 8 data bits, no parity, one stop bit:
//
//	port, err := serial.Open(name, &serial.Mode{
//		BaudRate: 115200,
//		DataBits: 8,
//		Parity:   serial.NoParity,
//		StopBits: serial.OneStopBit,
//	})
//	if err != nil {
//		return err
//	}
//	c, err := conn.NewSerialConn(port)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	lines, err := c.SendAndRecv(ctx, ">0C010000")
//
// Behaviour across calls:
//
//   - A SerialConn is safe for concurrent use. Calls are serialised, so one
//     request and its reply complete before the next request is written.
//   - SendAndRecv gives up after 5 seconds, even when ctx allows longer.
//   - A request is written no sooner than 10 ms after the previous reply was
//     read.
//   - Close may be called more than once. After Close, Send and SendAndRecv
//     return io.EOF.
package conn
