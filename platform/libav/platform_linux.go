//go:build !android
// +build !android

package libav

import (
	"github.com/xaionaro-go/camera"
	"github.com/xaionaro-go/camera/platform/v4l2"
)

const (
	InputFormat = "v4l2"
)

func InputStringFromDevicePath(devicePath camera.DevicePath) (string, error) {
	return devicePath, nil
}

func (Platform) ListCameras() ([]camera.DevicePath, error) {
	return v4l2.NewPlatform().ListCameras()
}

func (Platform) ListFormats(
	devicePath string,
) (camera.Formats, error) {
	return v4l2.NewPlatform().ListFormats(devicePath)
}
