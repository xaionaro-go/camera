package ximage

import (
	"unsafe"
)

const (
	y0CbY1CrSize = unsafe.Sizeof(Y0CbY1Cr{})
)

func (p *YUYV) setY0CbY1CrBytes(b []byte) {
	sliceLen := len(b) / int(y0CbY1CrSize)
	p.Y0CbY1Cr = (unsafe.Slice((*Y0CbY1Cr)(unsafe.Pointer(unsafe.SliceData(b))), sliceLen))
}
