package v4l2

import (
	"context"
	"fmt"
	"time"

	"github.com/blackjack/webcam"
	"github.com/xaionaro-go/camera"
)

type Camera struct {
	Camera      *webcam.Webcam
	Width       uint32
	Height      uint32
	PixelFormat webcam.PixelFormat
	FPS         camera.Fraction
}

var _ camera.Camera = (*Camera)(nil)

func (c *Camera) StartStreaming() error {
	return c.Camera.StartStreaming()
}

func (c *Camera) StopStreaming() error {
	return c.Camera.StopStreaming()
}

func (c *Camera) Close() error {
	return c.Camera.Close()
}

func (c *Camera) GetFormat() camera.Format {
	return camera.Format{
		Width:       uint64(c.Width),
		Height:      uint64(c.Height),
		PixelFormat: PixelFormatFromV4L2(c.PixelFormat),
		FPS:         c.FPS,
	}
}

func (c *Camera) GetFrames(
	ctx context.Context,
) (camera.FramesData, error) {
	for tryCount := 0; tryCount < 10*int(c.FPS.Float64()); tryCount++ {
		if err := c.WaitForFrame(ctx); err != nil {
			return nil, fmt.Errorf("unable to wait for a frame: %w", err)
		}

		b, frameID, err := c.Camera.GetFrame()
		if err != nil {
			return nil, fmt.Errorf("unable to read a frame: %w", err)
		}

		if len(b) != 0 {
			return &Frame{
				FrameID: frameID,
				Data:    b,
			}, nil
		}
		if err := c.Camera.ReleaseFrame(frameID); err != nil {
			return nil, fmt.Errorf("cannot release an allocated frame (%d): %w", frameID, err)
		}

		time.Sleep(time.Duration(float64(time.Second) / c.FPS.Float64()))
	}

	return nil, fmt.Errorf("internal error: we always get a zero-sized frame")
}

func (c *Camera) ReleaseFrames(frame camera.FramesData) error {
	return c.Camera.ReleaseFrame(frame.(*Frame).FrameID)
}

func (c *Camera) WaitForFrame(ctx context.Context) error {
	for {
		err := c.Camera.WaitForFrame(1)
		switch err.(type) {
		case *webcam.Timeout:
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				continue
			}
		case nil:
			return nil
		default:
			return err
		}
	}
}
