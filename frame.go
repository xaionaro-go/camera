package camera

import (
	"image"
)

type FramesCompressed interface {
	Bytes() []byte
}

type Frame interface {
	Image() image.Image
}

type imageWrapper struct {
	Img image.Image
}

func (w imageWrapper) Image() image.Image {
	return w.Img
}
