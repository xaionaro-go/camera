package camera

import (
	"fmt"
	"image"
	"io"
)

type FrameDecoder interface {
	io.Closer
	AllocateImage() image.Image
	WriteFrames([]byte) error
	DecodeFrame(image.Image) (image.Image, error)
}

func NewFrameDecoder(format Format) (FrameDecoder, error) {
	switch format.PixelFormat {
	case PixelFormatMJPEG:
		return newFrameDecoderMJPEG(format), nil
	}

	switch {
	case format.PixelFormat.IsRaw():
		decoder, err := newFrameDecoderRaw(format)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize a decoder of raw frames: %w", err)
		}
		return decoder, nil
	}

	return nil, fmt.Errorf("the support of pixel format '%v' is not implemented, yet", format.PixelFormat)
}
