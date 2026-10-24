package conn

import (
	"bytes"
	"context"
	"errors"
	"io"
	"slices"
	"testing"
	"time"

	bugserial "go.bug.st/serial"
)

// fakePort serves Read from a fixed input and records what is written. The
// embedded Port is nil, so a call to a method not overridden here panics.
type fakePort struct {
	bugserial.Port
	in          *bytes.Reader
	out         bytes.Buffer
	readTimeout time.Duration
	resets      int
	closes      int
}

func newFakePort(input string) *fakePort {
	return &fakePort{in: bytes.NewReader([]byte(input))}
}

func (p *fakePort) Read(b []byte) (int, error) { return p.in.Read(b) }

func (p *fakePort) Write(b []byte) (int, error) { return p.out.Write(b) }

func (p *fakePort) Drain() error { return nil }

func (p *fakePort) ResetInputBuffer() error {
	p.resets++
	return nil
}

func (p *fakePort) ResetOutputBuffer() error {
	p.resets++
	return nil
}

func (p *fakePort) Close() error {
	p.closes++
	return nil
}

func (p *fakePort) SetReadTimeout(t time.Duration) error {
	p.readTimeout = t
	return nil
}

func TestNewSerialConnPreparesPort(t *testing.T) {
	port := newFakePort("")
	if _, err := NewSerialConn(port); err != nil {
		t.Fatal(err)
	}
	if port.resets != 2 || port.readTimeout != 500*time.Millisecond {
		t.Errorf("resets = %d, read timeout = %s, want 2 and 500ms", port.resets, port.readTimeout)
	}
}

func TestSendAndRecvReadsUpToReply(t *testing.T) {
	port := newFakePort("log line\r\n\r\n<0C0500020000\r\nnext\r\n")
	c, err := NewSerialConn(port)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.SendAndRecv(context.Background(), ">0C010000")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"log line", "<0C0500020000"}; !slices.Equal(resp, want) {
		t.Errorf("resp = %q, want %q", resp, want)
	}
	if got, want := port.out.String(), ">0C010000\r\n"; got != want {
		t.Errorf("written = %q, want %q", got, want)
	}
	if got, want := port.in.Len(), len("next\r\n"); got != want {
		t.Errorf("%d bytes left unread, want %d", got, want)
	}
}

func TestSendAndRecvCanceled(t *testing.T) {
	c, err := NewSerialConn(newFakePort("<0C0500020000\r\n"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := c.SendAndRecv(ctx, ">0C010000"); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want %v", err, context.Canceled)
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	port := newFakePort("")
	c, err := NewSerialConn(port)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Errorf("second Close = %v, want nil", err)
	}
	if port.closes != 1 {
		t.Errorf("port closed %d times, want 1", port.closes)
	}
	if err := c.Send(context.Background(), ">0C010000"); !errors.Is(err, io.EOF) {
		t.Errorf("Send after Close = %v, want %v", err, io.EOF)
	}
	if _, err := c.SendAndRecv(context.Background(), ">0C010000"); !errors.Is(err, io.EOF) {
		t.Errorf("SendAndRecv after Close = %v, want %v", err, io.EOF)
	}
}
