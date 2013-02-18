package goastro

import (
	"errors"
	"math"
)

type EquatorialPos struct {
	RA   Angle
	Decl Angle
}

type EarthPos struct {
	Lat, Long Angle
}

type HorizontalPos struct {
	Azi, Alt Angle
}

type Positioner interface {
	Position(t TD) EquatorialPos
}

type DailyPositioner interface {
	// Returns position for 0h TD
	Position(d Date) EquatorialPos
}

// Wraps a DailyPositioner to make a fully-featured Positioner
type InterpolatedPositioner struct {
	dp DailyPositioner
}

func (ip InterpolatedPositioner) Position(t TD) EquatorialPos {
	p1 := ip.dp.Position(t.date.AddDays(-1))
	p2 := ip.dp.Position(t.date)
	p3 := ip.dp.Position(t.date.AddDays(1))
	ra := Interpolate3(float64(p1.RA), float64(p2.RA), float64(p3.RA), t.hours/24)
	decl := Interpolate3(float64(p1.Decl), float64(p2.Decl), float64(p3.Decl), t.hours/24)
	return EquatorialPos{Angle(ra), Angle(decl)}
}

type rstType int

const (
	risingT rstType = iota
	settingT
	transitT
)

// Ch 15 p.102
func risingSetting(p Positioner, h0 Angle, ep EarthPos, d Date, rst rstType) (UT, error) {
	//log.Print(rst)
	pos2 := p.Position(TD{d, 0})
	α2 := pos2.RA
	δ2 := pos2.Decl
	//log.Print("pos2 =", pos2)
	φ := ep.Lat
	L := -ep.Long
	cosH0 := (sin(h0) - sin(φ)*sin(δ2)) / (cos(φ) * cos(δ2))
	//log.Print("cosH0 =", cosH0)
	if cosH0 < -1 || cosH0 > 1 {
		return UT{}, errors.New("no rising on day")
	}
	H0 := acos(cosH0)
	//log.Print("H0 =", H0)
	dayΘ := ApparentSiderealTime(UT{d, 0})
	//log.Print("dayΘ =", dayΘ)
	m0 := (α2 + L - dayΘ).Degrees() / 360 // transit as day fraction
	//log.Print("m0 =", m0)
	var m float64
	switch rst {
	case risingT:
		m = m0 - H0.Degrees()/360 // rising as day fraction
	case settingT:
		m = m0 + H0.Degrees()/360 // setting as day fraction
	case transitT:
		m = m0
	}
	//log.Print("m =", m)
	if m < 0 {
		m++
	}
	if m >= 1 {
		m--
	}
	//log.Print("m =", m)
	//log.Print("time =", UT{d, 24*m})
	for {
		timeθ := (dayΘ + Degrees(360.985647*m)).Normalize() // sidereal time
		//log.Print("timeθ =", timeθ)
		pos := p.Position(UT{d, 24 * m}.TD())
		//log.Print("pos =", pos)
		α := pos.RA
		δ := pos.Decl
		H := (timeθ - L - α).Normalize180()
		//log.Print("H =", H)
		h := asin(sin(φ)*sin(δ) + cos(φ)*cos(δ)*cos(H))
		//log.Print("h =", h)
		var Δm float64
		if rst == transitT {
			Δm = -H.Degrees() / 360
		} else {
			Δm = (h - h0).Degrees() / (360 * cos(δ) * cos(φ) * sin(H))
		}
		//log.Print("Δm =", Δm)
		m += Δm
		//log.Print("m =", m)
		//log.Print("time =", UT{d, 24*m})
		if math.Abs(Δm) < 0.00001 {
			break
		}
	}
	return UT{d, 24 * m}, nil
}

func Rising(p Positioner, h0 Angle, ep EarthPos, d Date) (UT, error) {
	return risingSetting(p, h0, ep, d, risingT)
}

func Setting(p Positioner, h0 Angle, ep EarthPos, d Date) (UT, error) {
	return risingSetting(p, h0, ep, d, settingT)
}

func Transit(p Positioner, ep EarthPos, d Date) (UT, error) {
	return risingSetting(p, 0, ep, d, transitT)
}

func (p EquatorialPos) HorizontalPos(θ0 Angle, ep EarthPos) HorizontalPos {
	α := p.RA
	δ := p.Decl
	L := -ep.Long
	φ := ep.Lat
	H := θ0 - L - α
	A := atan2(sin(H), cos(H)*sin(φ)-tan(δ)*cos(φ))
	h := asin(sin(φ)*sin(δ) + cos(φ)*cos(δ)*cos(H))
	return HorizontalPos{A, h}
}

func (p EquatorialPos) TransitAltitude(lat Angle) Angle {
	δ := p.Decl
	φ := lat
	return asin(sin(φ)*sin(δ) + cos(φ)*cos(δ))
}

func (p HorizontalPos) EquatorialPos(θ0 Angle, ep EarthPos) EquatorialPos {
	A := p.Azi
	h := p.Alt
	L := -ep.Long
	φ := ep.Lat
	H := atan2(sin(A), cos(A)*sin(φ)+tan(h)*cos(φ))
	δ := asin(sin(φ)*sin(h) - cos(φ)*cos(h)*cos(A))
	α := θ0 - L - H
	return EquatorialPos{α, δ}
}

func (p EquatorialPos) HourAngle(h, lat Angle) Angle {
    φ := lat
    δ := p.Decl
    cosH := (sin(h) - sin(φ) * sin(δ)) / (cos(φ) * cos(δ))
    return acos(cosH)
}
