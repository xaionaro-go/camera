package v4l2

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/blackjack/webcam"
	"github.com/xaionaro-go/camera"
)

type Platform struct{}

func NewPlatform() Platform {
	return Platform{}
}

func (Platform) ListCameras() ([]camera.DevicePath, error) {
	const devDir = "/dev/"
	entries, err := os.ReadDir(devDir)
	if err != nil {
		return nil, fmt.Errorf("unable list the available devices: %w", err)
	}

	var result []camera.DevicePath
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "video") {
			result = append(result, path.Join(devDir, entry.Name()))
		}
	}

	return result, nil
}
func (Platform) ListFormats(
	devicePath string,
) (camera.Formats, error) {
	webCam, err := webcam.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open '%s' as V4L2 camera: %w", devicePath, err)
	}
	defer webCam.Close()

	var result []camera.Format
	for pixelFormat := range webCam.GetSupportedFormats() {
		for _, frameSize := range webCam.GetSupportedFrameSizes(pixelFormat) {
			width := frameSize.MaxWidth
			height := frameSize.MaxHeight
			for _, fps := range webCam.GetSupportedFramerates(pixelFormat, width, height) {
				fpsValue := camera.Fraction{
					Numerator:   uint(fps.MaxDenominator),
					Denominator: uint(fps.MaxNumerator),
				}
				result = append(result, camera.Format{
					Width:       uint64(width),
					Height:      uint64(height),
					PixelFormat: PixelFormatFromV4L2(pixelFormat),
					FPS:         fpsValue,
				})
			}
		}
	}
	return result, nil
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
	webCam, err := webcam.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open '%s' as V4L2 camera: %w", devicePath, err)
	}

	pixFmt, width, height, err := webCam.SetImageFormat(
		PixelFormatToV4L2(format.PixelFormat),
		uint32(format.Width),
		uint32(format.Height),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to configure the image format: %w", err)
	}

	err = webCam.SetFramerate(format.FPS.Float32())
	if err != nil {
		return nil, fmt.Errorf("unable to configure the frame rate: %w", err)
	}

	return &Camera{
		Camera: webCam,
		Format: camera.Format{
			Width:       uint64(width),
			Height:      uint64(height),
			PixelFormat: PixelFormatFromV4L2(pixFmt),
			FPS:         format.FPS,
		},
	}, nil
}
