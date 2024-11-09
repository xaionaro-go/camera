package camera

import (
	"fmt"
	"image"
	"io"

	"github.com/mattn/go-mjpeg"
)

type frameDecompressorMJPEG struct {
	Format       Format
	FrameWriter  io.WriteCloser
	MJPEGDecoder *mjpeg.Decoder
}

var _ FrameDecompressor = (*frameDecompressorMJPEG)(nil)

func newFrameDecompressorMJPEG() *frameDecompressorMJPEG {
	r, w := io.Pipe()
	return &frameDecompressorMJPEG{
		FrameWriter:  w,
		MJPEGDecoder: mjpeg.NewDecoder(r, "image/jpeg"),
	}
}

func (d *frameDecompressorMJPEG) Close() error {
	return d.FrameWriter.Close()
}

func (d *frameDecompressorMJPEG) NewImage() image.Image {
	return nil
}

func (d *frameDecompressorMJPEG) WriteCompressed(
	compressed FramesCompressed,
) error {
	_, err := d.FrameWriter.Write(compressed.Bytes())
	return err
}

func (d *frameDecompressorMJPEG) DecompressNext() (Frame, error) {
	img, err := d.MJPEGDecoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("unable to decode the frame: %w", err)
	}

	return imageWrapper{img}, nil
}

func (d *frameDecompressorMJPEG) ReleaseFrame(
	frame Frame,
) {

}
