package camera

import (
	"fmt"
	"image"
)

type frameDecoderHEIC struct {
}

func newFrameDecoderHEIC(Format) (*frameDecoderHEIC, error) {
	return nil, fmt.Errorf("not implemented")
}

func (frameDecoderHEIC) Close() error {
	return fmt.Errorf("not implemented")
}
func (frameDecoderHEIC) NewImage() image.Image {
	return nil
}
func (frameDecoderHEIC) WriteFrames(FramesData) error {
	return fmt.Errorf("not implemented")
}
func (frameDecoderHEIC) DecodeFrame(image.Image) (image.Image, error) {
	return nil, fmt.Errorf("not implemented")
}
