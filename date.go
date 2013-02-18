package goastro

import (
	"time"
)

type Date struct {
	Year, Month, Day int
}

func MakeDate(t time.Time) Date {
	return Date{t.Year(), int(t.Month()), t.Day()}
}

func (d Date) Time() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func (d Date) AddDays(n int) Date {
	return MakeDate(d.Time().AddDate(0, 0, n))
}

func (d Date) compareTo(d2 Date) int {
	if d.Year != d2.Year {
		return d.Year - d2.Year
	}

	if d.Month != d2.Month {
		return d.Month - d2.Month
	}

	return d.Day - d2.Day
}

var gregorianStart = Date{1582, 10, 15}

type JulianDay float64

type Time interface {
	Date() Date
	Hours() float64
}

// Ch 7 p.60
func MakeJulianDay(t Time) JulianDay {
	y, m, d := t.Date().Year, t.Date().Month, t.Date().Day
	if m <= 2 {
		y--
		m += 12
	}
	b := 0
	if t.Date().compareTo(gregorianStart) >= 0 {
		a := y / 100
		b = 2 - a + (a / 4)
	}
	fd := float64(d) + t.Hours()/24
	jdint := int(365.25*float64(y+4716)) + int(30.6001*float64(m+1)) + b
	jd := float64(jdint) + fd - 1524.5
	return JulianDay(jd)
}

func DayOfYear(date time.Time) int {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	jan0 := time.Date(date.Year(), time.January, 0, 0, 0, 0, 0, date.Location())
	diff := date.Sub(jan0)
	return int(diff / (24 * time.Hour))
}
