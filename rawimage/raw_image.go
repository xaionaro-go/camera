package rawimage

import (
	"fmt"
	"image"

	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/ximage"
)

func NewRawImage(
	format *camera.Format,
	frameBytes []byte,
) (_ret image.Image, _err error) {
	defer func() {
		if _err != nil {
			_err = fmt.Errorf("pixel format %v: %w", format.PixelFormat, _err)
			return
		}
	}()

	switch format.PixelFormat {
	case camera.PixelFormatYUYV:
		return NewRawImageYUYV(frameBytes, uint(format.Width), uint(format.Height))
	case camera.PixelFormatNV12:
		return NewRawImageNV12(frameBytes, uint(format.Width), uint(format.Height))
	default:
		return nil, fmt.Errorf("unexpected pixel")
	}
}

func NewRawImageNV12(
	frameBytes []byte,
	width, height uint,
) (*ximage.NV12, error) {
	dstImg := ximage.NewNV12NoAlloc(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: int(width),
			Y: int(height),
		},
	})

	if err := dstImg.SetBytes(frameBytes); err != nil {
		return nil, fmt.Errorf("unable to set bytes: %w", err)
	}
	return dstImg, nil
}

func NewRawImageYUYV(
	frameBytes []byte,
	width, height uint,
) (*ximage.YUYV, error) {
	dstImg := ximage.NewYUYVNoAlloc(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: int(width),
			Y: int(height),
		},
	})

	if err := dstImg.SetBytes(frameBytes); err != nil {
		return nil, fmt.Errorf("unable to set bytes: %w", err)
	}
	return dstImg, nil
}
