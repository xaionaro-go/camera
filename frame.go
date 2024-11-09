package camera

import (
	"image"
)

type FramesData interface {
	Bytes() []byte
}

type Frame interface {
	Image() image.Image
}
