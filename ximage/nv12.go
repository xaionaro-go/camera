package ximage

import (
	"fmt"
	"image"
	"image/color"
)

type CbCr struct {
	Cb uint8
	Cr uint8
}

type NV12 struct {
	Y       []uint8
	CbCr    []CbCr
	YStride int
	Rect    image.Rectangle
}

func NewNV12(r image.Rectangle) *NV12 {
	p := &NV12{
		YStride: r.Max.X - r.Min.X,
		Rect:    r,
	}
	_, bytesExpected := p.sizes()
	if bytesExpected < 0 {
		panic("ximage: Rectangle has huge or negative dimensions")
	}
	if err := p.SetBytes(make([]byte, bytesExpected)); err != nil {
		panic(err)
	}
	return p
}

func (p *NV12) sizes() (int, int) {
	w := p.Rect.Max.X - p.Rect.Min.X
	h := p.Rect.Max.Y - p.Rect.Min.Y
	pixelCount := w * h
	bytesExpected := (pixelCount*3 + 1) / 2
	return pixelCount, bytesExpected
}

func (p *NV12) SetBytes(b []byte) error {
	// see https://www.kernel.org/doc/html/v4.10/media/uapi/v4l/pixfmt-nv12.html

	pixelCount, bytesExpected := p.sizes()
	if bytesExpected != len(b) {
		return fmt.Errorf("the size of the provided image does not match the expected size: expected:%d, received:%d", bytesExpected, len(b))
	}

	p.Y = b[:pixelCount:pixelCount]
	if err := p.SetCbCrBytes(b[pixelCount:bytesExpected:bytesExpected]); err != nil {
		return fmt.Errorf("unable to set the CbCr bytes: %w", err)
	}
	return nil
}

func (p *NV12) SetCbCrBytes(b []byte) error {
	// see https://www.kernel.org/doc/html/v4.10/media/uapi/v4l/pixfmt-nv12.html

	pixelCount, _ := p.sizes()
	bytesExpected := (pixelCount + 1) / 2
	if len(b) != bytesExpected {
		return fmt.Errorf("the size the provided slice does not match the expected size: expected:%d, received:%d", bytesExpected, len(b))
	}
	if len(b)%int(cbCrSize) != 0 {
		return fmt.Errorf("the size of the bytes slice is not a multiple of the CbCr struct: %d %% %d == %d", len(b), cbCrSize, len(b)%int(cbCrSize))
	}
	p.setCbCrBytes(b)
	return nil
}

func (p *NV12) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *NV12) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NV12) At(x, y int) color.Color {
	return p.YCbCrAt(x, y)
}

func (p *NV12) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.YCbCrAt(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (p *NV12) YCbCrAt(x, y int) color.YCbCr {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}
	yi := p.YOffset(x, y)
	ci := p.COffset(x, y)
	cbCr := p.CbCr[ci]
	return color.YCbCr{
		p.Y[yi],
		cbCr.Cb,
		cbCr.Cr,
	}
}

// YOffset returns the index of the first element of Y that corresponds to
// the pixel at (x, y).
func (p *NV12) YOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.YStride + (x - p.Rect.Min.X)
}

// COffset returns the index of the first element of Cb or Cr that corresponds
// to the pixel at (x, y).
func (p *NV12) COffset(x, y int) int {
	return (y/2-p.Rect.Min.Y/2)*p.YStride/2 + (x/2 - p.Rect.Min.X/2)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *NV12) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &NV12{}
	}

	yi := p.YOffset(r.Min.X, r.Min.Y)
	ci := p.COffset(r.Min.X, r.Min.Y)
	return &NV12{
		Y:       p.Y[yi:],
		CbCr:    p.CbCr[ci:],
		YStride: p.YStride,
		Rect:    r,
	}
}

func (p *NV12) Opaque() bool {
	return true
}
