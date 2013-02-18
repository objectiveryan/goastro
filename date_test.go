package goastro

import (
    "testing"
    "time"
    "math"
)

func date(y, m, d int) time.Time {
    return time.Date(y, time.Month(m), d, 0, 0, 0,0, time.UTC)
}

func dayFrac(f float64) time.Duration {
    nsec := int64(f * float64(24 * time.Hour))
    return time.Duration(nsec)
}

func TestJulianDay(t *testing.T) {
    cases := []struct {
        y, m int
        d float64
        want float64
    }{
        {1957, 10, 4.81, 2436116.31},
        {333, 1, 27.5, 1842713},
        {1977, 4, 26.4, 2443259.9},
        {2000, 1, 1.5, 2451545},
        {1987, 1, 27, 2446822.5},
        {1987, 6, 19.5, 2446966},
        {1988, 1, 27, 2447187.5},
        {1988, 6, 19.5, 2447332},
        {1900, 1, 1, 2415020.5},
        {1600, 1, 1, 2305447.5},
        {1600, 12, 31, 2305812.5},
        {837, 4, 10.3, 2026871.8},
        {-1000, 7, 12.5, 1356001},
        //{-1000, 2, 29, 1355866.5}, // a leap year acc to julian but not gregorian?
        {-1001, 8, 17.9, 1355671.4},
        {-4712, 1, 1.5, 0},
    }

    for _, c := range cases {
        d, df := math.Modf(c.d)
        date := UT{Date{c.y, c.m, int(d)}, 24*df}
        got := MakeJulianDay(date)
        if math.Abs(float64(got) - c.want) > 0.0001 {
            t.Errorf("JulianDay(%v) = %f, want %f", date, got, c.want)
        }
    }
}

func TestDayOfYear(t *testing.T) {
    cases := []struct {
        y, m, d int
        want int
    }{
        {1978, 11, 14, 318},
        {1988, 4, 22, 113},
    }

    for _, c := range cases {
        date := time.Date(c.y, time.Month(c.m), c.d, 0, 0, 0,0, time.UTC)
        got := DayOfYear(date)
        if got != c.want {
            t.Errorf("DayOfYear(%v) = %d, want %d", date, got, c.want)
        }
    }
}
