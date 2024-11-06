package camera

import (
	"context"
	"io"
)

type FrameID int

type Camera interface {
	io.Closer
	StartStreaming() error
	StopStreaming() error
	GetFormat() Format
	GetRawFrames(context.Context, []byte) ([]byte, FrameID, error)
	ReleaseFrames(FrameID) error
}
