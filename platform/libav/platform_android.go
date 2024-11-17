package libav

import (
	"fmt"

	"github.com/xaionaro-go/camera"
)

const (
	InputFormat = "android_camera"
)

const (
	DevicePathFront = camera.DevicePath("FRONT")
	DevicePathBack  = camera.DevicePath("BACK")
)

func DevicePathToIndex(devicePath camera.DevicePath) int {
	switch devicePath { // see https://ffmpeg.org/ffmpeg-devices.html#android_005fcamera
	case DevicePathBack:
		return 0
	case DevicePathFront:
		return 1
	}
	return -1 // not found
}

func InputStringFromDevicePath(devicePath camera.DevicePath) (string, error) {
	index := DevicePathToIndex(devicePath)
	if index < 0 {
		return "", fmt.Errorf("invalid device path: '%s'", devicePath)
	}
	return fmt.Sprintf("%d", index), nil
}

func (Platform) ListCameras() ([]camera.DevicePath, error) {
	return []camera.DevicePath{DevicePathBack, DevicePathFront}, nil
}

func (Platform) ListFormats(
	devicePath string,
) (camera.Formats, error) {
	if DevicePathToIndex(devicePath) < 0 {
		return nil, fmt.Errorf("invalid path: '%s'", devicePath)
	}

	return camera.Formats{{
		Width:       1920,
		Height:      1080,
		PixelFormat: camera.PixelFormatNV12,
		FPS: camera.Fraction{
			Numerator:   30,
			Denominator: 1,
		},
	}}, nil
}
