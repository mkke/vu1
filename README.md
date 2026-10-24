Go client for communicating with a [VU-1 hub](https://vudials.com/) over a USB serial port.

## Packages

- `conn` -- the line-based request/response exchange over the serial port.
- `hub` -- the hub's command protocol: hub-wide commands on `Hub`, per-dial
  commands (needle, backlight, display) on `Dial`, and conversion of images
  to the dial display's 1-bit format.
- `graduation` -- draws needle scales (arcs, ticks, labels) for dial face
  images.

Each package's doc comment describes its call order and wire format; `go doc`
shows them.

## Usage

```go
c, err := conn.NewSerialConn(port) // port opened at 115200 baud, 8N1
if err != nil {
	return err
}
defer c.Close()

h := hub.New(c)
if err := h.ProvisionDials(ctx); err != nil {
	return err
}
if err := h.RescanBus(ctx); err != nil {
	return err
}
devices, err := h.GetDevicesMap(ctx)
if err != nil {
	return err
}
for idx := range devices.OnlineDials() {
	if err := h.NewDial(idx).SetPercent(ctx, 50); err != nil {
		return err
	}
}
```

`git.mallorn.de/vu1ctl` is a complete program built on these packages.

Serial communication protocol is derived from https://github.com/SasaKaranovic/VU-Server.

Dial image generation is derived from https://github.com/shabaz123/MeterDialCreator,
see also https://community.element14.com/technologies/open-source-hardware/b/blog/posts/meter-dial-creator-toolkit.