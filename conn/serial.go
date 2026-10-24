package conn

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	bugserial "go.bug.st/serial"
)

// SerialConn is a Conn over a serial port. See the package documentation
// for its timing and concurrency behaviour.
type SerialConn struct {
	port     bugserial.Port
	lock     sync.Mutex
	closed   bool
	minIdle  time.Duration
	lastRecv time.Time
}

// NewSerialConn discards anything buffered on port, sets its read timeout to
// 500 ms and returns a Conn that owns it.
func NewSerialConn(port bugserial.Port) (Conn, error) {
	if err := port.ResetInputBuffer(); err != nil {
		return nil, err
	}
	if err := port.ResetOutputBuffer(); err != nil {
		return nil, err
	}
	if err := port.SetReadTimeout(500 * time.Millisecond); err != nil {
		return nil, err
	}

	c := &SerialConn{port: port, minIdle: 10 * time.Millisecond}
	return c, nil
}

// Close closes the port. Calls after the first return nil.
func (c *SerialConn) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return c.port.Close()
}

func (c *SerialConn) readLine(ctx context.Context) (string, error) {
	// we don't want to consume any bytes after the first LF
	readBuf := make([]byte, 1)
	lineBuf := bytes.Buffer{}
nextLine:
	for {
		if err := ctx.Err(); err != nil {
			return string(lineBuf.Bytes()), err
		}

		if n, err := c.port.Read(readBuf); err != nil {
			return string(lineBuf.Bytes()), err
		} else if n == 0 {
			continue
		}
		switch readBuf[0] {
		case '\r':
			// discard
		case '\n':
			if lineBuf.Len() == 0 {
				// don't return empty lines
				goto nextLine
			}
			return string(lineBuf.Bytes()), nil
		default:
			lineBuf.Write(readBuf)
		}
	}
}

func (c *SerialConn) send(_ context.Context, req string) error {
	if idleDuration := time.Since(c.lastRecv); idleDuration < c.minIdle {
		time.Sleep(c.minIdle - idleDuration)
	}

	if _, err := c.port.Write([]byte(req + "\r\n")); err != nil {
		return err
	}

	return c.port.Drain()
}

func (c *SerialConn) recv(ctx context.Context) ([]string, error) {
	defer func() {
		c.lastRecv = time.Now()
	}()

	var resp []string
	for ctx.Err() == nil {
		line, err := c.readLine(ctx)
		if err != nil {
			return nil, err
		}
		resp = append(resp, line)
		if strings.HasPrefix(line, "<") {
			return resp, nil
		}
	}
	return resp, context.Cause(ctx)
}

// Send writes req followed by CR LF and waits until the port has sent it.
func (c *SerialConn) Send(ctx context.Context, req string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return io.EOF
	}

	return c.send(ctx, req)
}

// SendAndRecv writes req and reads the reply lines, giving up after 5
// seconds.
func (c *SerialConn) SendAndRecv(ctx context.Context, req string) (resp []string, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, io.EOF
	}

	ctx, cancel := context.WithTimeoutCause(ctx, 5*time.Second, errors.New("timeout waiting for send/recv"))
	defer cancel()

	if err := c.send(ctx, req); err != nil {
		return nil, err
	}

	return c.recv(ctx)
}
