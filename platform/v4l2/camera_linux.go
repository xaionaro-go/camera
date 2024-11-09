package v4l2

import (
	"context"
	"fmt"
	"time"

	"github.com/blackjack/webcam"
	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/rawimage"
)

type Camera struct {
	Camera *webcam.Webcam
	Format camera.Format
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
	return c.Format
}

func (c *Camera) GetFrame(
	ctx context.Context,
) (camera.Frame, error) {
	for tryCount := 0; tryCount < 10*int(c.Format.FPS.Float64()); tryCount++ {
		if err := c.WaitForFrame(ctx); err != nil {
			return nil, fmt.Errorf("unable to wait for a frame: %w", err)
		}

		b, frameID, err := c.Camera.GetFrame()
		if err != nil {
			return nil, fmt.Errorf("unable to read a frame: %w", err)
		}

		if len(b) != 0 {
			img, err := rawimage.NewRawImage(&c.Format, b)
			if err != nil {
				return nil, fmt.Errorf("unable to parse the image: %w", err)
			}

			return &Frame{
				FrameID: frameID,
				Frame:   img,
			}, nil
		}
		if err := c.Camera.ReleaseFrame(frameID); err != nil {
			return nil, fmt.Errorf("cannot release an allocated frame (%d): %w", frameID, err)
		}

		time.Sleep(time.Duration(float64(time.Second) / c.Format.FPS.Float64()))
	}

	return nil, fmt.Errorf("internal error: we always get a zero-sized frame")
}

func (c *Camera) ReleaseFrame(frame camera.Frame) error {
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
