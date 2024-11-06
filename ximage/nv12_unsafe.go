package ximage

import (
	"unsafe"
)

const (
	cbCrSize = unsafe.Sizeof(CbCr{})
)

func (p *NV12) setCbCrBytes(b []byte) {
	sliceLen := len(b) / int(cbCrSize)
	p.CbCr = (unsafe.Slice((*CbCr)(unsafe.Pointer(unsafe.SliceData(b))), sliceLen))
}
