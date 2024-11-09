package camera

import (
	"fmt"
	"image"

	"github.com/xaionaro-go/camera/ximage"
)

type frameDecoderRaw struct {
	Format    Format
	FrameData FramesData
}

var _ FrameDecoder = (*frameDecoderRaw)(nil)

func newFrameDecoderRaw(format Format) (*frameDecoderRaw, error) {
	switch format.PixelFormat {
	case PixelFormatYUYV, PixelFormatNV12:
	default:
		return nil, fmt.Errorf("the support of pixel format '%v' is not implemented, yet", format.PixelFormat)
	}

	return &frameDecoderRaw{
		Format: format,
	}, nil
}

func (d *frameDecoderRaw) Close() error {
	return nil
}

func (d *frameDecoderRaw) NewImage() image.Image {
	size := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: int(d.Format.Width),
			Y: int(d.Format.Height),
		},
	}
	switch d.Format.PixelFormat {
	case PixelFormatNV12:
		return ximage.NewNV12(size)
	case PixelFormatYUYV:
		return image.NewYCbCr(size, image.YCbCrSubsampleRatio422)
	default:
		panic(fmt.Errorf("internal error, this panic is supposed to be unreachable (see function newFrameDecoderRaw): the support of pixel format '%v' is not implemented, yet", d.Format.PixelFormat))
	}
}

func (d *frameDecoderRaw) WriteFrames(framesData FramesData) error {
	d.FrameData = framesData
	return nil
}

func (d *frameDecoderRaw) DecodeFrame(
	dstImg image.Image,
) (_ret image.Image, _err error) {
	defer func() {
		if _err != nil {
			_err = fmt.Errorf("pixel format %v: %w", d.Format.PixelFormat, _err)
			return
		}
		d.FrameData = nil
	}()
	if frame, ok := d.FrameData.(Frame); ok {
		return frame.Image(), nil
	}

	switch d.Format.PixelFormat {
	case PixelFormatYUYV:
		return d.decodeFrameYUYV(typeAssertOrZero[*ximage.YUYV](dstImg))
	case PixelFormatNV12:
		return d.decodeFrameNV12(typeAssertOrZero[*ximage.NV12](dstImg))
	default:
		return nil, fmt.Errorf("unexpected pixel")
	}
}

func (d *frameDecoderRaw) decodeFrameNV12(
	dstImg *ximage.NV12,
) (*ximage.NV12, error) {
	if dstImg == nil {
		dstImg = ximage.NewNV12(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: int(d.Format.Width),
				Y: int(d.Format.Height),
			},
		})
	}

	if err := dstImg.SetBytes(d.FrameData.Bytes()); err != nil {
		return nil, fmt.Errorf("unable to set bytes: %w", err)
	}
	return dstImg, nil
}

func (d *frameDecoderRaw) decodeFrameYUYV(
	dstImg *ximage.YUYV,
) (*ximage.YUYV, error) {
	if dstImg == nil {
		dstImg = ximage.NewYUYV(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: int(d.Format.Width),
				Y: int(d.Format.Height),
			},
		})
	}

	if err := dstImg.SetBytes(d.FrameData.Bytes()); err != nil {
		return nil, fmt.Errorf("unable to set bytes: %w", err)
	}
	return dstImg, nil
}
