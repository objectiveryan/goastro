package goastro

func Interpolate3(y1, y2, y3, n float64) float64 {
	if n < -1 || n > 1 {
		panic("Interpolate: n not in [-1, 1]")
	}
	a := y2 - y1
	b := y3 - y2
	c := b - a
	return y2 + n * (a + b + n * c) / 2
}
