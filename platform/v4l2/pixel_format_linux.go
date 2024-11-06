package v4l2

import (
	"github.com/blackjack/webcam"
	"github.com/xaionaro-go/camera"
)

func PixelFormatFromV4L2(
	v4l2PixFmt webcam.PixelFormat,
) camera.PixelFormat {
	return camera.PixelFormatFromUint32(uint32(v4l2PixFmt))
}

func PixelFormatToV4L2(
	pixFmt camera.PixelFormat,
) webcam.PixelFormat {
	return webcam.PixelFormat(pixFmt.Uint32())
}
