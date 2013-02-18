package goastro

import (
	"math"
	"testing"
)

func TestSunPosition(t *testing.T) {
	time := TD{Date{1992, 10, 13}, 0}
	got := SunPosition(time)
	wantRA := Degrees(-161.61917)
	wantDecl := Degrees(-7.78507)
	if math.Abs((got.RA-wantRA).Degrees()) > 0.0001 {
		t.Errorf("SunPosition(%v).RA == %v, want %v", time, got.RA, wantRA)
	}
	if math.Abs((got.Decl-wantDecl).Degrees()) > 0.0001 {
		t.Errorf("SunPosition(%v).Decl == %v, want %v", time, got.Decl, wantDecl)
	}
}
