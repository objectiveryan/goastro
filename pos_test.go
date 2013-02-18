package goastro

import (
	"fmt"
	"math"
	"testing"
	"time"
)

type Venus15a struct{}

func (v Venus15a) Position(d Date) EquatorialPos {
	switch d {
	case Date{1988, 3, 19}:
		return EquatorialPos{Degrees(40.68021), Degrees(18.04761)}
	case Date{1988, 3, 20}:
		return EquatorialPos{Degrees(41.73129), Degrees(18.44092)}
	case Date{1988, 3, 21}:
		return EquatorialPos{Degrees(42.78204), Degrees(18.82742)}
	}
	panic("Venus15a.Position() date not supported")
}

func TestRising(t *testing.T) {
	var venus Venus15a
	h0 := Degrees(-0.5667)
	ep := EarthPos{Degrees(42 + 20/60.), -Degrees(71 + 5/60.)}
	d := Date{1988, 3, 20}
	want := 12 + 25/60.
	got, err := Rising(InterpolatedPositioner{venus}, h0, ep, d)
	if err != nil {
		t.Error(err)
	}
	if got.date != d {
		t.Errorf("Rising returned wrong date: got %v, expected %v", got, d)
	}
	if math.Abs(got.hours-want)*60 > 1 {
		t.Errorf("Rising returned wrong time: got %v, want %v", TimeOfDay(got.hours), TimeOfDay(want))
	}
}

func TestSetting(t *testing.T) {
	var venus Venus15a
	h0 := Degrees(-0.5667)
	ep := EarthPos{Degrees(42 + 20/60.), -Degrees(71 + 5/60.)}
	d := Date{1988, 3, 20}
	want := 2 + 55/60.
	got, err := Setting(InterpolatedPositioner{venus}, h0, ep, d)
	if err != nil {
		t.Error(err)
	}
	if got.date != d {
		t.Errorf("Setting returned wrong date: got %v, expected %v", got, d)
	}
	if math.Abs(got.hours-want)*60 > 1 {
		t.Errorf("Setting returned wrong time: got %v, want %v", TimeOfDay(got.hours), TimeOfDay(want))
	}
}

func TestTransit(t *testing.T) {
	var venus Venus15a
	ep := EarthPos{Degrees(42 + 20/60.), -Degrees(71 + 5/60.)}
	d := Date{1988, 3, 20}
	want := 19 + 41/60.
	got, err := Transit(InterpolatedPositioner{venus}, ep, d)
	if err != nil {
		t.Error(err)
	}
	if got.date != d {
		t.Errorf("Transit returned wrong date: got %v, expected %v", got, d)
	}
	if math.Abs(got.hours-want)*60 > 1 {
		t.Errorf("Transit returned wrong time: got %v, want %v", TimeOfDay(got.hours), TimeOfDay(want))
	}
}

type SunPositioner struct{}

func (sp SunPositioner) Position(t TD) EquatorialPos {
	return SunPosition(t)
}

func GetTime(t UT, err error) UT {
	if err != nil {
		panic(err)
	}
	return t
}

var ET *time.Location

func init() {
	var err error
	ET, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
}

func DateTime(d Date, h, m int) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, h, m, 0, 0, ET)
}

func Asr(p Positioner, ep EarthPos, dhuhr UT) (UT, error) {
	θ0 := ApparentSiderealTime(dhuhr)
	dhuhrAlt := p.Position(dhuhr.TD()).HorizontalPos(θ0, ep).Alt
	dhuhrShadow := 1 / tan(dhuhrAlt)
	asrShadow := dhuhrShadow + 1
	asrAlt := atan(1 / asrShadow)
	return Setting(p, asrAlt, ep, dhuhr.date)
}

func TestPrayers(t *testing.T) {
	var sun SunPositioner
	d := Date{2012, 12, 4}
	twilight := Degrees(-15)
	horizon := ArcMinutes(-50)
	ep := EarthPos{Degrees(42.36462), Degrees(-71.11518)}
	fajr := GetTime(Rising(sun, twilight, ep, d))
	sunrise := GetTime(Rising(sun, horizon, ep, d))
	dhuhr := GetTime(Transit(sun, ep, d))
	asr := GetTime(Asr(sun, ep, dhuhr))
	maghrib := GetTime(Setting(sun, horizon, ep, d))
	isha := GetTime(Setting(sun, twilight, ep, d))
	got := []UT{fajr, sunrise, dhuhr, asr, maghrib, isha}
	want := []struct{ h, m int }{{5, 34}, {6, 57}, {11, 35}, {13, 54}, {16, 13}, {17, 36}}
	for i := range got {
		if got[i].date != d {
			t.Errorf("got[%d].date == %v, want %v", i, got[i].date, d)
		}
		wantHours := float64(5+want[i].h) + float64(want[i].m)/60

		if math.Abs(got[i].hours-wantHours) > 1/60. {
			t.Errorf("got[%d] == %v, want %v", i, TimeOfDay(got[i].hours), TimeOfDay(wantHours))
		}
	}
}

func TestEquatorialToHorizontalPos(t *testing.T) {
	ut := UT{Date{1987, 4, 10}, 19 + 21/60.}
	eq := EquatorialPos{Hours(23 + ms(9, 16.641)), -Degrees(6 + ms(43, 11.61))}
	θ0 := ApparentSiderealTime(ut)
	ep := EarthPos{Degrees(38 + ms(55, 17)), -Degrees(77 + ms(3, 56))}
	got := eq.HorizontalPos(θ0, ep)
	want := HorizontalPos{Degrees(68.0337), Degrees(15.1249)}

	if math.Abs((got.Azi - want.Azi).Degrees()) > 0.0002 {
		t.Errorf("got.Azi == %f, want %f", got.Azi, want.Azi)
	}

	if math.Abs((got.Alt - want.Alt).Degrees()) > 0.0001 {
		t.Errorf("got.Alt == %f, want %f", got.Alt, want.Alt)
	}
}
