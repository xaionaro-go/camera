package camera

import (
	"context"
	"io"
)

type Camera interface {
	io.Closer
	StartStreaming() error
	StopStreaming() error
	GetFormat() Format
	GetFrames(context.Context) (FramesData, error)
	ReleaseFrames(FramesData) error
}
