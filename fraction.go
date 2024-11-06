package camera

type Fraction struct {
	Numerator   uint
	Denominator uint
}

func (f Fraction) Float64() float64 {
	return float64(f.Numerator) / float64(f.Denominator)
}

func (f Fraction) Float32() float32 {
	return float32(f.Numerator) / float32(f.Denominator)
}
