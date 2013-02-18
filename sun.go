package goastro

import ()

func SunPosition(t TD) EquatorialPos {
	T := (float64(MakeJulianDay(t)) - 2451545) / 36525
	//log.Print("T = ", T)
	L0 := Degrees(280.46646 + T*(36000.76983+T*0.0003032)).Normalize()
	//log.Print("L0 = ", L0)
	M := Degrees(357.52911 + T*(35999.05029-T*0.0001537)).Normalize()
	//log.Print("M = ", M)
	C := Degrees((1.914602-T*(0.004817+T*0.000014))*sin(M) +
		(0.019993-T*0.000101)*sin(2*M) +
		0.000289*sin(3*M))
	//log.Print("C = ", C)
	long := L0 + C
	//log.Print("long = ", long)
	Ω := Degrees(125.04 - 1934.136*T)
	//log.Print("Ω = ", Ω)
	λ := long - Degrees(0.00569+0.00478*sin(Ω))
	//log.Print("λ = ", λ)
	ε := TrueObliquity(t)
	//log.Print("ε = ", ε)
	α := atan2(cos(ε)*sin(λ), cos(λ))
	//log.Print("α = ", α)
	δ := asin(sin(ε) * sin(λ))
	//log.Print("δ = ", δ)
	return EquatorialPos{α, δ}
}
