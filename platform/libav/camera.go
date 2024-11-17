package libav

import (
	"context"
	"fmt"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
	"github.com/xaionaro-go/camera"
)

type Camera struct {
	*astikit.Closer
	Input  *Input
	Format camera.Format
}

var _ camera.Camera = (*Camera)(nil)

func (c *Camera) StartStreaming() error {
	return nil
}

func (c *Camera) StopStreaming() error {
	return nil
}

func (c *Camera) GetFormat() camera.Format {
	return c.Format
}

func (c *Camera) GetFrame(
	ctx context.Context,
) (camera.Frame, error) {
	packet := astiav.AllocPacket()
	for tryCount := 0; tryCount < 10*int(c.Format.FPS.Float64()); tryCount++ {
		err := c.Input.FormatContext.ReadFrame(packet)
		if err != nil {
			return nil, fmt.Errorf("unable to read a frame: %w", err)
		}
		if len(packet.Data()) != 0 {
			break
		}
	}
	if len(packet.Data()) == 0 {
		return nil, fmt.Errorf("the packet is empty")
	}

	return &Frame{
		Packet: packet,
		Camera: c,
	}, nil
}

func (c *Camera) ReleaseFrame(frame camera.Frame) error {
	return frame.(*Frame).Close()
}
