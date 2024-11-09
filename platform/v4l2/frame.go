package v4l2

import "github.com/xaionaro-go/camera"

type Frame struct {
	FrameID uint32
	Data    []byte
}

var _ camera.FramesData = (*Frame)(nil)

func (f *Frame) Bytes() []byte {
	return f.Data
}
