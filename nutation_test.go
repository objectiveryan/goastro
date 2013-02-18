package goastro

import (
	"testing"
	"math"
)

func arcSecondDifference(a, b Angle) float64 {
	return math.Abs((a - b).ArcSeconds())
}

func TestLongitudeNutation(t *testing.T) {
	date := Date{1987, 4, 10}
	want := ArcSeconds(-3.788)
	got := LongitudeNutation(TD{date, 0})
	if arcSecondDifference(want, got) > 0.01 {
		t.Errorf("LongitudeNutation(%v) == %v\", want %v\"", date, got, want)
	}
}

func TestObliquityNutation(t *testing.T) {
	date := Date{1987, 4, 10}
	want := ArcSeconds(9.443)
	got := ObliquityNutation(TD{date, 0})
	if arcSecondDifference(want, got) > 0.01 {
		t.Errorf("ObliquityNutation(%v) == %v\", want %v\"", date, got, want)
	}
}

func TestMeanObliquity(t *testing.T) {
	date := Date{1987, 4, 10}
	want := Degrees(23 + ms(26, 27.407))
	got := MeanObliquity(TD{date, 0})
	if arcSecondDifference(want, got) > 0.001 {
		t.Errorf("MeanObliquity(%v) == %v\", want %v\"", date, got, want)
	}
}

func TestTrueObliquity(t *testing.T) {
	date := Date{1987, 4, 10}
	want := Degrees(23 + ms(26, 36.850))
	got := TrueObliquity(TD{date, 0})
	if arcSecondDifference(want, got) > 0.01 {
		t.Errorf("TrueObliquity(%v) == %v, want %v", date, got, want)
	}
}
