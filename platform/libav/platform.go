package libav

import (
	"fmt"

	"github.com/asticode/go-astikit"
	"github.com/xaionaro-go/camera"
)

type Platform struct{}

func NewPlatform() Platform {
	return Platform{}
}

func (Platform) OpenCameraCompressed(
	devicePath camera.DevicePath,
	format camera.Format,
	compression camera.Compression,
	compressionQuality camera.CompressionQuality,
) (camera.CameraCompressed, error) {
	return nil, fmt.Errorf("not supported")
}

func (Platform) OpenCamera(
	devicePath string,
	format camera.Format,
) (camera.Camera, error) {
	inputString, err := InputStringFromDevicePath(devicePath)
	if err != nil {
		return nil, err
	}

	input, err := NewInput(InputFormat, inputString, format)
	if err != nil {
		return nil, fmt.Errorf("unable to open the camera: %w", err)
	}

	c := &Camera{
		Closer: astikit.NewCloser(),
		Input:  input,
		Format: format,
	}
	c.Closer.Add(input.Free)
	return c, nil
}
