package ximage

import (
	"fmt"
	"image"
	"image/color"
)

type Y0CbY1Cr struct {
	Y0 uint8
	Cb uint8
	Y1 uint8
	Cr uint8
}

type YUYV struct {
	Y0CbY1Cr []Y0CbY1Cr
	YStride  int
	Rect     image.Rectangle
}

func NewYUYV(r image.Rectangle) *YUYV {
	p := &YUYV{
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

func NewYUYVNoAlloc(r image.Rectangle) *YUYV {
	p := &YUYV{
		YStride: r.Max.X - r.Min.X,
		Rect:    r,
	}
	_, bytesExpected := p.sizes()
	if bytesExpected < 0 {
		panic("ximage: Rectangle has huge or negative dimensions")
	}
	return p
}

func (p *YUYV) sizes() (int, int) {
	w := p.Rect.Max.X - p.Rect.Min.X
	h := p.Rect.Max.Y - p.Rect.Min.Y
	pixelCount := w * h
	bytesExpected := pixelCount * 2
	return pixelCount, bytesExpected
}

func (p *YUYV) SetBytes(b []byte) error {
	// see https://www.kernel.org/doc/html/v4.10/media/uapi/v4l/pixfmt-yuyv.html
	if err := p.SetY0CbY1CrBytes(b); err != nil {
		return fmt.Errorf("unable to set the Y0CbY1Cr bytes: %w", err)
	}
	return nil
}

func (p *YUYV) SetY0CbY1CrBytes(b []byte) error {
	// see https://www.kernel.org/doc/html/v4.10/media/uapi/v4l/pixfmt-yuyv.html
	_, bytesExpected := p.sizes()
	if len(b) != bytesExpected {
		return fmt.Errorf("the size the provided slice does not match the expected size: expected:%d, received:%d", bytesExpected, len(b))
	}
	if len(b)%int(y0CbY1CrSize) != 0 {
		return fmt.Errorf("the size of the bytes slice is not a multiple of the Y0CbY1Cr struct: %d %% %d == %d", len(b), y0CbY1CrSize, len(b)%int(y0CbY1CrSize))
	}
	p.setY0CbY1CrBytes(b)
	return nil
}

func (p *YUYV) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *YUYV) Bounds() image.Rectangle {
	return p.Rect
}

func (p *YUYV) At(x, y int) color.Color {
	return p.YCbCrAt(x, y)
}

func (p *YUYV) RGBA64At(x, y int) color.RGBA64 {
	r, g, b, a := p.YCbCrAt(x, y).RGBA()
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (p *YUYV) YCbCrAt(x, y int) color.YCbCr {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}

	offset := p.Y0CbY1CrOffset(x, y)
	y0CbY1Cr := p.Y0CbY1Cr[offset]

	odd := uint8(x) & 1
	return color.YCbCr{
		y0CbY1Cr.Y0*(1-odd) + y0CbY1Cr.Y1*odd,
		y0CbY1Cr.Cb,
		y0CbY1Cr.Cr,
	}
}

func (p *YUYV) Y0CbY1CrOffset(x, y int) int {
	return ((y-p.Rect.Min.Y)*p.YStride + (x - p.Rect.Min.X)) / 2
}

// COffset returns the index of the first element of Cb or Cr that corresponds
// to the pixel at (x, y).
func (p *YUYV) COffset(x, y int) int {
	return (y/2-p.Rect.Min.Y/2)*p.YStride/2 + (x/2 - p.Rect.Min.X/2)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *YUYV) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &YUYV{}
	}

	offset := p.Y0CbY1CrOffset(r.Min.X, r.Min.Y)
	return &YUYV{
		Y0CbY1Cr: p.Y0CbY1Cr[offset:],
		YStride:  p.YStride,
		Rect:     r,
	}
}

func (p *YUYV) Opaque() bool {
	return true
}
