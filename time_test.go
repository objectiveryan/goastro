package goastro

import (
	"testing"
	"math"
	"fmt"
)

func ms(m int, s float64) float64 {
	return (float64(m) + s / 60) / 60
}

func timeSecondDifference(a, b Angle) float64 {
	return math.Abs((a - b).Hours() * 3600)
}

type HMS Angle

func (hms HMS) String() string {
	h, hf := math.Modf(Angle(hms).Hours())
	m, mf := math.Modf(hf * 60)
	s, sf := math.Modf(mf * 60)
	return fmt.Sprintf("%dh%dm%ds.%d", int(h), int(m), int(s), int(sf * 10000))
}

func TestMeanSiderealTime(t *testing.T) {
	cases := []struct {
		time UT
		want Angle
	}{
		{UT{Date{1987, 4, 10}, 0}, Hours(13 + ms(10, 46.3668))},
		{UT{Date{1987, 4, 10}, 19+21/60.}, Hours(8 + ms(34, 57.0896))},
	}
	for _, c := range cases {
		got := MeanSiderealTime(c.time)
		if timeSecondDifference(c.want, got) > 0.0001 {
			t.Errorf("MeanSiderealTime(%v) == %v, want %v", c.time, HMS(got), HMS(c.want))
		}
	}
}

func TestApparentSiderealTime(t *testing.T) {
	time := UT{Date{1987, 4, 10}, 0}
	want := Hours(13 + ms(10, 46.1351))
	got := ApparentSiderealTime(time)
	if timeSecondDifference(want, got) > 0.005 {
		t.Log("Diff =", timeSecondDifference(want, got))
		t.Errorf("ApparentSiderealTime(%v) == %v, want %v", time, HMS(got), HMS(want))
	}
}

func TestDeltaT(t *testing.T) {
	cases := []struct {
		year int
		want float64
	}{
		{1900, -2.8},
		{1910, 10.4},
		{1920, 21.1},
		{1930, 24},
		{1940, 24.3},
		{1950, 29.1},
		{1960, 33.1},
		{1970, 40.2},
		{1980, 50.5},
		{1990, 56.9},
		{1996, 61.6},
		//{2000, 63.8},
		//{2004, 64.6},
		//{2008, 65.5},
		//{2010, 66.1},
	}

	for _, c := range cases {
		d := Date{c.year, 1, 1}
		got := DeltaT(d)
		if math.Abs(got - c.want) > 0.9 {
			t.Errorf("DeltaT(%v) == %f, want %f", d, got, c.want)
		}
	}
}
