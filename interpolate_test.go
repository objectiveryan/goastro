package goastro

import (
	"math"
	"testing"
)

func TestInterpolate3(t *testing.T) {
	cases := []struct {
		y1, y2, y3, n float64
		want          float64
	}{
		{0.884226, 0.877366, 0.870531, 4.35 / 24, 0.876125},
	}

	for _, c := range cases {
		got := Interpolate3(c.y1, c.y2, c.y3, c.n)
		if math.Abs(c.want-got) > 0.000001 {
			t.Errorf("Interpolate3(%f, %f, %f, %f) == %f, want %f", c.y1, c.y2, c.y3, c.n, got, c.want)
		}
	}
}
