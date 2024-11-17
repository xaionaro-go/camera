package libav

import (
	"fmt"
	"image"

	"github.com/asticode/go-astiav"
	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/rawimage"
)

type Frame struct {
	Packet *astiav.Packet
	Camera *Camera
}

var _ camera.Frame = (*Frame)(nil)

func (f *Frame) Image() image.Image {
	frameBytes := f.Packet.Data()
	img, err := rawimage.NewRawImage(&f.Camera.Format, frameBytes)
	if err != nil {
		panic(fmt.Errorf("unable to parse the image: %w", err))
	}
	return img
}

func (f *Frame) Close() error {
	f.Packet.Free()
	return nil
}
