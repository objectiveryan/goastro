package goastro

import (
	"fmt"
	"math"
)

type Angle float64 // degrees

func (a Angle) Degrees() float64 {
	return float64(a)
}

func (a Angle) Radians() float64 {
	return float64(a) / 180 * math.Pi
}

func (a Angle) Hours() float64 {
	return a.Degrees() / 15
}

func (a Angle) ArcMinutes() float64 {
	return a.Degrees() * 60
}

func (a Angle) ArcSeconds() float64 {
	return a.Degrees() * 3600
}

func (a Angle) String() string {
	if a.IsNaN() {
		return "NaN°"
	}
	d, df := math.Modf(a.Degrees())
	m, mf := math.Modf(df * 60)
	s, sf := math.Modf(mf * 60)
	return fmt.Sprintf("%d°%d'%d\".%d/%.5f", int(d), int(m), int(s), int(sf*1000), a.Degrees())
}

func (a Angle) Sin() float64 {
	return sin(a)
}

func (a Angle) Cos() float64 {
	return cos(a)
}

func (a Angle) Tan() float64 {
	return tan(a)
}

func (a Angle) IsNaN() bool {
	return math.IsNaN(float64(a))
}

func Degrees(d float64) Angle {
	return Angle(d)
}

func ArcMinutes(s float64) Angle {
	return Angle(s / 60)
}

func ArcSeconds(s float64) Angle {
	return Angle(s / 3600)
}

func Radians(r float64) Angle {
	return Angle(r / math.Pi * 180)
}

func Hours(h float64) Angle {
	return Angle(h * 15)
}

func (a Angle) Normalize() Angle {
	d := math.Mod(float64(a), 360)
	if d < 0 {
		d += 360
	}
	return Angle(d)
}

func (a Angle) Normalize180() Angle {
	a = a.Normalize()
	if a.Degrees() >= 180 {
		return a - Degrees(360)
	}
	return a
}

func sin(a Angle) float64 {
	return math.Sin(a.Radians())
}

func cos(a Angle) float64 {
	return math.Cos(a.Radians())
}

func tan(a Angle) float64 {
	return math.Tan(a.Radians())
}

func asin(x float64) Angle {
	return Radians(math.Asin(x))
}

func acos(x float64) Angle {
	return Radians(math.Acos(x))
}

func atan(x float64) Angle {
	return Radians(math.Atan(x))
}

func Atan(x float64) Angle {
	return Radians(math.Atan(x))
}

func atan2(y float64, x float64) Angle {
	return Radians(math.Atan2(y, x))
}

func Atan2(y float64, x float64) Angle {
	return Radians(math.Atan2(y, x))
}
