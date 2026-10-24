// Package hub implements the command protocol of a VU-1 hub and its dials.
//
// A Hub sends commands over a conn.Conn. Commands that address the hub as a
// whole are methods of Hub. Commands that address one dial are methods of the
// Dial returned by Hub.NewDial. Every method waits for the hub's reply and
// returns an error wrapping ErrCommandFailed when the hub answers with a
// status code other than StatusOk.
//
// # Wire format
//
// A request is one line: '>', then the command, the data type and the
// payload length in bytes as 2, 2 and 4 upper-case hex digits, then the
// payload as hex. A reply has the same layout with '<' in place of '>'.
// MarshallCommand builds a request and UnmarshallCommand decodes a reply.
//
// # Call order
//
// The hub learns which dials are attached when it is provisioned and its bus
// is rescanned, so a session starts with:
//
//	h := hub.New(c)
//	if err := h.ProvisionDials(ctx); err != nil {
//		return err
//	}
//	if err := h.RescanBus(ctx); err != nil {
//		return err
//	}
//	devices, err := h.GetDevicesMap(ctx)
//	if err != nil {
//		return err
//	}
//	for idx := range devices.OnlineDials() {
//		dial := h.NewDial(idx)
//		...
//	}
//
// # Dial display
//
// A dial's display is 200 by 144 pixels (DialImageWidth, DialImageHeight)
// at one bit per pixel. Showing an image takes three calls:
// Dial.DisplayGotoXY to set the write position, Dial.DisplaySendImageData to
// upload the pixels, and Dial.DisplayShow to put them on screen.
// DrawImage runs a draw function on a white canvas of that size and returns
// the result.
//
// ConvertImageData reduces an image to black and white and packs it column
// by column, left to right. Each byte holds 8 vertically adjacent pixels,
// the topmost in the most significant bit. A white pixel is a 1 bit.
// DisplaySendImageData sends the packed data in chunks of 1000 bytes.
package hub
