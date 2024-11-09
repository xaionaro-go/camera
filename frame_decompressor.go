package camera

import (
	"fmt"
	"io"
)

type FrameDecompressor interface {
	io.Closer

	WriteCompressed(FramesCompressed) error
	DecompressNext() (Frame, error)
	ReleaseFrame(Frame)
}

func NewFrameDecompressor(
	compression Compression,
) (FrameDecompressor, error) {
	switch compression {
	case CompressionHEIC:
		return newFrameDecompressorHEIC()
	case CompressionMJPEG:
		return newFrameDecompressorMJPEG(), nil
	default:
		return nil, fmt.Errorf("compression '%s' is not supported", compression)
	}
}
