package camera

import (
	"fmt"
	"image"
)

type frameDecompressorHEIC struct {
}

var _ FrameDecompressor = (*frameDecompressorHEIC)(nil)

func newFrameDecompressorHEIC() (*frameDecompressorHEIC, error) {
	return nil, fmt.Errorf("not implemented")
}

func (frameDecompressorHEIC) Close() error {
	return fmt.Errorf("not implemented")
}
func (frameDecompressorHEIC) NewImage() image.Image {
	return nil
}
func (frameDecompressorHEIC) WriteCompressed(FramesCompressed) error {
	return fmt.Errorf("not implemented")
}
func (frameDecompressorHEIC) DecompressNext() (Frame, error) {
	return nil, fmt.Errorf("not implemented")
}
func (frameDecompressorHEIC) ReleaseFrame(Frame) {
}
