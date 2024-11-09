package camera

import (
	"fmt"
	"image"
	"io"
)

type FrameDecoder interface {
	io.Closer
	NewImage() image.Image
	WriteFrames(FramesData) error
	DecodeFrame(image.Image) (image.Image, error)
}

func NewFrameDecoder(format Format) (_ FrameDecoder, _err error) {
	defer func() {
		if _err != nil {
			_err = fmt.Errorf("unable to initialize a frame decoder for %s/%s: %w", format.Compression, format.PixelFormat, _err)
		}
	}()

	switch format.Compression {
	case CompressionMJPEG:
		if format.PixelFormat != PixelFormatAuto {
			return nil, fmt.Errorf("a pixel format cannot be forced, when MJPEG compression is used")
		}
		return newFrameDecoderMJPEG(format), nil
	case CompressionHEIC:
		return newFrameDecoderHEIC(format)
	}

	decoder, err := newFrameDecoderRaw(format)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize a decoder of raw frames: %w", err)
	}
	return decoder, nil
}
