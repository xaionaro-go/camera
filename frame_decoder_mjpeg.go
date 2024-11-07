package camera

import (
	"fmt"
	"image"
	"io"

	"github.com/mattn/go-mjpeg"
)

type frameDecoderMJPEG struct {
	Format       Format
	FrameWriter  io.WriteCloser
	MJPEGDecoder *mjpeg.Decoder
}

var _ FrameDecoder = (*frameDecoderMJPEG)(nil)

func newFrameDecoderMJPEG(format Format) *frameDecoderMJPEG {
	r, w := io.Pipe()
	return &frameDecoderMJPEG{
		Format:       format,
		FrameWriter:  w,
		MJPEGDecoder: mjpeg.NewDecoder(r, "image/jpeg"),
	}
}

func (d *frameDecoderMJPEG) Close() error {
	return d.FrameWriter.Close()
}

func (d *frameDecoderMJPEG) NewImage() image.Image {
	return nil
}

func (d *frameDecoderMJPEG) WriteFrames(frames []byte) error {
	_, err := d.FrameWriter.Write(frames)
	return err
}

func (d *frameDecoderMJPEG) DecodeFrame(
	dstImg image.Image,
) (image.Image, error) {
	img, err := d.MJPEGDecoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("unable to decode the frame: %w", err)
	}

	return img, nil
}
