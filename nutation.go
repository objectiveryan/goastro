package goastro

import ()

// Ch 22 p.145
var nutationCoefficients = []struct {
	D, M, MM, F, Ω float64
	sinc, sint float64
	cosc, cost float64
}{
	{0,0,0,0,1, -171996,-174.2, 92025,8.9},
	{-2,0,0,2,2, -13187,-1.6, 5736,-3.1},
	{0,0,0,2,2, -2274,-0.2, 977,-0.5},
	{0,0,0,0,2, 2062,0.2, -895,0.5},
	{0,1,0,0,0, 1426,-3.4, 54,-0.1},
	{0,0,1,0,0, 712,0.1, -7,0},
	{-2,1,0,2,2, -517,1.2, 224,-0.6},
	{0,0,0,2,1, -386,-0.4, 200,0},
	{0,0,1,2,2, -301,0, 129,-0.1},
	{-2,-1,0,2,2, 217,-0.5, -95,0.3},
	{-2,0,1,0,0, -158,0, 0,0},
	// there are many more terms in the book!
}

// Ch 22 p.145
func nutationTerms(T float64) (D, M, MM, F, Ω float64) {
	// Mean elongation of the Moon from the Sun
	D = 297.85036 + T * (445267.111480 + T * (-0.0019142 + T / 189474))

	// Mean anomaly of the Sun (Earth)
	M = 357.52772 + T * (35999.050340 - T * (0.0001603 + T / 300000))

	// Mean anomaly of the Moon
	MM = 134.96298 + T * (477198.867398 + T * (0.0086972 + T / 56250))

	// Moon's argument of latitude
	F = 93.27191 + T * (483202.017538 + T * (-0.0036825 + T / 327270))

	// Longitude of the ascending node of the Moon's mean orbot on the
	// eliptic, measured from the mean equinox of the date
	Ω = 125.04452 + T * (-1934.136261 + T * (0.0020708 + T / 450000))
	return
}

// Ch 22 p.143
func LongitudeNutation(t TD) Angle {
	T := (float64(MakeJulianDay(t)) - 2451545) / 36525
	D, M, MM, F, Ω := nutationTerms(T)
	Δψ := 0.0
	for _, c := range nutationCoefficients {
		arg := D * c.D + M * c.M + MM * c.MM + F * c.F + Ω * c.Ω
		coeff := c.sinc + c.sint * T
		Δψ += coeff * sin(Degrees(arg))
	}
	return ArcSeconds(Δψ / 10000) // units are 0.0001"
}

// Ch 22 p.143
// Less accurate version
func LongitudeNutation0(t TD) Angle {
	jde := MakeJulianDay(t)
	T := (float64(jde) - 2451545) / 36525
	Ω := Degrees(125.04452 - 1934.136261*T + 0.0020708*T*T + T*T*T/450000)
	L := Degrees(280.4665 + 36000.7698*T)
	Lp := Degrees(218.3165 + 481267.8813*T)
	secs := -17.20*sin(Ω) - 1.32*sin(2*L) - 0.23*sin(2*Lp) + 0.21*sin(2*Ω)
	return ArcSeconds(secs)
}

// Ch 22 p.143
func ObliquityNutation(t TD) Angle {
	T := (float64(MakeJulianDay(t)) - 2451545) / 36525
	D, M, MM, F, Ω := nutationTerms(T)
	Δε := 0.0
	for _, c := range nutationCoefficients {
		arg := D * c.D + M * c.M + MM * c.MM + F * c.F + Ω * c.Ω
		coeff := c.cosc + c.cost * T
		Δε += coeff * cos(Degrees(arg))
	}
	return ArcSeconds(Δε / 10000) // units are 0.0001"
}


// Ch 22 p.143
// Less accurate version
func ObliquityNutation0(t TD) Angle {
	T := (float64(MakeJulianDay(t)) - 2451545) / 36525
	Ω := Degrees(125.04452 - 1934.136261*T + 0.0020708*T*T + T*T*T/450000)
	L := Degrees(280.4665 + 36000.7698*T)
	Lp := Degrees(218.3165 + 481267.8813*T)
	secs := 9.20*cos(Ω) + 0.57*cos(2*L) - 0.10*cos(2*Lp) - 0.09*cos(2*Ω)
	return ArcSeconds(secs)
}

// Ch 22 p.147
func MeanObliquity(t TD) Angle {
	T := (float64(MakeJulianDay(t)) - 2451545) / 36525
	secs := 21.448 - T * (46.8150 + T * (0.00059 - T * 0.001813))
	return Degrees(23 + (26+secs/60)/60)
}

// Ch 22 p.147
func TrueObliquity(t TD) Angle {
	return MeanObliquity(t) + ObliquityNutation(t)
}
