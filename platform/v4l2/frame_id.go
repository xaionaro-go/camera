package v4l2

import (
	"github.com/xaionaro-go/camera"
)

type FrameID uint32

func (f FrameID) FrameID() camera.FrameID {
	return f
}
