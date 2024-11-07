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
	GetRawFrames(context.Context, []byte) (RawFrames, FrameID, error)
	ReleaseRawFrames(FrameID) error
}
