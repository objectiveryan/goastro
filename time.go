package goastro

import (
	"fmt"
	"math"
	"time"
)

// type Date defined in date.go

// Dynamical Time
type TD struct {
	date  Date
	hours float64
}

// Universal Time
type UT struct {
	date  Date
	hours float64
}

type TimeOfDay float64

func (t TimeOfDay) String() string {
	h, hf := math.Modf(float64(t))
	m, mf := math.Modf(hf * 60)
	s, sf := math.Modf(mf * 60)
	return fmt.Sprintf("%02d:%02d:%02d.%01d", int(h), int(m), int(s), int(sf*10))
}

func dateHours(t time.Time) float64 {
	s := float64(t.Second()) + float64(t.Nanosecond())/1e9
	m := float64(t.Minute()) + s/60
	return float64(t.Hour()) + m/60
}

func MakeUT(t time.Time) UT {
	t = t.UTC()
	return UT{MakeDate(t), dateHours(t)}
}

// Chapter 12 p.87
func MeanSiderealTime(t UT) Angle {
	T := (float64(MakeJulianDay(UT{t.date, 0})) - 2451545) / 36525
	JD := float64(MakeJulianDay(t))
        return Degrees(280.46061837 + 360.98564736629 * (JD - 2451545) + T*T * (0.000387933 - T / 38710000)).Normalize()
}

// Chapter 12 p.88
func ApparentSiderealTime(t UT) Angle {
	// XXX This mixes calculations with UT and TD, but that's what the book
	// does (intentionally), since it doesn't make much of a difference.
	td := TD{t.date, t.hours}
	correction := Angle(float64(LongitudeNutation(td)) * cos(TrueObliquity(td)))
	return MeanSiderealTime(t) + correction
}

// ΔT = TD - UT (unit = seconds)
// Chapter 10 p.80
func DeltaT(d Date) float64 {
	if d.Year >= 1900 && d.Year <= 1997 {
		dtime := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
		epoch := time.Date(1900, time.January, 0, 0, 0, 0, 0, time.UTC)
		nanoDiff := int64(dtime.Sub(epoch))
		daysDiff := nanoDiff / 1e9 / 60 / 60 / 24
		t := float64(daysDiff) / 36525.
		return ((((((((58353.42*t-232424.66)*t+372919.88)*t-303191.19)*t+124906.15)*t-18756.33)*t-2637.80)*t+815.20)*t+87.24)*t - 2.44
	} else if d.Year >= 2000 {
		t := float64(d.Year-2000) / 100
		ΔT := 102 + 102*t + 25.3*t*t
		if d.Year < 2100 {
			ΔT += 0.37 * float64(d.Year-2100)
		}
		return ΔT
	}
	panic("DeltaT: years before 2000 unsupported")
}

func (t TD) Date() Date {
	return t.date
}

func (t TD) Hours() float64 {
	return t.hours
}

func (t TD) String() string {
	return fmt.Sprintf("TD %v %v", t.date, TimeOfDay(t.hours))
}

func (t TD) UT() UT {
	h := t.hours - DeltaT(t.date)/60/60
	if h < 0 {
		return UT{t.date.AddDays(-1), h + 24}
	} else if h >= 24 {
		return UT{t.date.AddDays(1), h - 24}
	}
	return UT{t.date, h}
}

func (t UT) Date() Date {
	return t.date
}

func (t UT) Hours() float64 {
	return t.hours
}

func (t UT) String() string {
	return fmt.Sprintf("UT %v %v", t.date, TimeOfDay(t.hours))
}

func (t UT) TD() TD {
	h := t.hours + DeltaT(t.date)/60/60
	if h >= 24 {
		return TD{t.date.AddDays(1), h - 24}
	} else if h < 0 {
		return TD{t.date.AddDays(-1), h + 24}
	}
	return TD{t.date, h}
}

func (t UT) Time() time.Time {
	h, hf := math.Modf(t.hours)
	m, mf := math.Modf(hf * 60)
	s, sf := math.Modf(mf * 60)
	return time.Date(t.date.Year, time.Month(t.date.Month), t.date.Day,
		int(h), int(m), int(s), int(sf*1e9), time.UTC)
}
