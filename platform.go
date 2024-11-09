package camera

type DevicePath = string

type Platform interface {
	ListCameras() ([]DevicePath, error)

	OpenCamera(
		devicePath DevicePath,
		format Format,
	) (Camera, error)

	OpenCameraCompressed(
		devicePath DevicePath,
		format Format,
		compression Compression,
		compressionQuality CompressionQuality,
	) (CameraCompressed, error)

	ListFormats(
		devicePath string,
	) (Formats, error)
}

func ListCameras() ([]DevicePathAndPlatform, error) {
	return DefaultRegistry().ListCameras()
}
