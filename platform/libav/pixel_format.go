package libav

import (
	"strings"

	"github.com/xaionaro-go/camera"
)

func PixelFormatToLibAV(pixFmt camera.PixelFormat) string {
	return strings.ToLower(string(pixFmt))
}
