package libav

import (
	"fmt"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
	"github.com/xaionaro-go/camera"
)

func init() {
	astiav.RegisterAllDevices()
}

type Input struct {
	*astikit.Closer
	*astiav.FormatContext
}

func NewInput(
	formatString string,
	inputString string,
	frameFormat camera.Format,
) (_ *Input, _err error) {
	input := &Input{
		Closer: astikit.NewCloser(),
	}
	defer func() {
		if _err != nil {
			input.Closer.Close()
		}
	}()

	inputFormat := astiav.FindInputFormat(formatString)
	if inputFormat == nil {
		return nil, fmt.Errorf("format '%s' not found", formatString)
	}

	input.FormatContext = astiav.AllocFormatContext()
	if input.FormatContext == nil {
		// TODO: is there a way to extract the actual error code or something?
		return nil, fmt.Errorf("unable to allocate a format context")
	}
	input.Closer.Add(input.FormatContext.Free)

	dict := astiav.NewDictionary()
	input.Closer.Add(dict.Free)

	if err := dict.Set("video_size", fmt.Sprintf("%dx%d", frameFormat.Width, frameFormat.Height), 0); err != nil {
		return nil, fmt.Errorf("unable to set the video_size in the dictionary: %w", err)
	}
	if err := dict.Set("pixel_format", PixelFormatToLibAV(frameFormat.PixelFormat), 0); err != nil {
		return nil, fmt.Errorf("unable to set the pixel_format in the dictionary: %w", err)
	}
	if err := dict.Set("framerate", fmt.Sprintf("%f", frameFormat.FPS.Float64()), 0); err != nil {
		return nil, fmt.Errorf("unable to set the framerate in the dictionary: %w", err)
	}

	if err := input.FormatContext.OpenInput(inputString, inputFormat, dict); err != nil {
		return nil, fmt.Errorf("unable to open input '%s':'%s': %w", formatString, inputString, err)
	}
	input.Closer.Add(input.FormatContext.CloseInput)

	if err := input.FormatContext.FindStreamInfo(nil); err != nil {
		return nil, fmt.Errorf("unable to get stream info: %w", err)
	}
	return input, nil
}
