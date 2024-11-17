//go:build !android
// +build !android

package allplatforms

import (
	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/platform/libav"
	"github.com/xaionaro-go/camera/platform/v4l2"
)

func Get(platID string) camera.Platform {
	switch platID {
	case "libav":
		return libav.Platform{}
	case "v4l2":
		return v4l2.Platform{}
	default:
		return nil
	}
}
