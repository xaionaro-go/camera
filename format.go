package camera

import (
	"encoding/binary"
)

type PixelFormat string

const (
	PixelFormatUndefined = PixelFormat("")
	PixelFormatAuto      = PixelFormat("*")

	PixelFormatMJPEG = PixelFormat("MJPG")
	PixelFormatYUYV  = PixelFormat("YUYV")
	PixelFormatNV12  = PixelFormat("NV12")
)

func PixelFormatByName(pixFmtName string) PixelFormat {
	return PixelFormat(pixFmtName)
}

func (pixFmt PixelFormat) Uint32() uint32 {
	return binary.NativeEndian.Uint32([]byte(pixFmt))
}

func (pixFmt PixelFormat) rawBitSize() uint32 {
	switch pixFmt {
	case PixelFormatYUYV:
		return 16
	case PixelFormatNV12:
		return 12
	}
	return 0
}

func (pixFmt PixelFormat) IsRaw() bool {
	switch pixFmt {
	case PixelFormatYUYV,
		PixelFormatNV12:
		return true
	}
	return false
}

func PixelFormatFromUint32(v uint32) PixelFormat {
	var value [4]byte
	binary.NativeEndian.PutUint32(value[:], v)
	return PixelFormat(value[:])
}

type Format struct {
	Width       uint64
	Height      uint64
	PixelFormat PixelFormat
	FPS         Fraction
}

type Formats []Format

func (s Formats) FilterByPixelFormat(pixFmt PixelFormat) Formats {
	var result Formats

	for _, f := range s {
		if f.PixelFormat == pixFmt {
			result = append(result, f)
		}
	}
	return result
}

func (s Formats) FilterByWidth(width uint64) Formats {
	var result Formats

	for _, f := range s {
		if f.Width == width {
			result = append(result, f)
		}
	}
	return result
}

func (s Formats) FilterByFPS(fps float64) Formats {
	var result Formats

	for _, f := range s {
		if f.FPS.Float64() == fps {
			result = append(result, f)
		}
	}
	return result
}

func (s Formats) BestResolution() Format {
	var best Format
	for _, f := range s {
		if f.Width*f.Height != best.Width*best.Height {
			if f.Width*f.Height > best.Width*best.Height {
				best = f
			}
			continue
		}
		if f.FPS.Float64() != best.FPS.Float64() {
			if f.FPS.Float64() > best.FPS.Float64() {
				best = f
			}
			continue
		}
		if f.PixelFormat.rawBitSize() != best.PixelFormat.rawBitSize() {
			if f.PixelFormat.rawBitSize() > best.PixelFormat.rawBitSize() {
				best = f
			}
			continue
		}
	}
	return best
}
