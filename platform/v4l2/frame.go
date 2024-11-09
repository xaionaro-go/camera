package v4l2

import (
	"image"

	"github.com/xaionaro-go/camera"
)

type Frame struct {
	FrameID uint32
	Frame   image.Image
}

var _ camera.Frame = (*Frame)(nil)

func (f *Frame) Image() image.Image {
	return f.Frame
}
