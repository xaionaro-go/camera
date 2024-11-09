package camera

import (
	"context"
	"io"
)

type CameraCommon interface {
	io.Closer
	StartStreaming() error
	StopStreaming() error
	GetFormat() Format
}

type Camera interface {
	CameraCommon
	GetFrame(context.Context) (Frame, error)
	ReleaseFrame(Frame) error
}

type CameraCompressed interface {
	CameraCommon
	GetCompressedFrames(context.Context) (FramesCompressed, error)
	ReleaseFrames(FramesCompressed) error
}
