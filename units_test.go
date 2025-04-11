package main

import (
	"math"
	"slices"
	"testing"
)

func almostEqual(lhs float64, rhs float64) bool {
	return math.Abs(lhs-rhs) < 0.00001

}

func TestExpressUnitAsSIUnits(t *testing.T) {
	length := Unit{"l", []float64{1, 0, 0}}
	mass := Unit{"m", []float64{0, 1, 0}}
	time := Unit{"t", []float64{0, 0, 1}}
	speed := Unit{"v", []float64{1, 0, -1}}
	momentum := Unit{"p", []float64{1, 1, -1}}
	angularMomentum := Unit{"L", []float64{2, 1, -1}}
	allUnits := []Unit{length, mass, time, speed, momentum, angularMomentum}
	siBaseUnits := []Unit{length, mass, time}
	for _, unit := range allUnits {
		exponents := ExpressUnit(unit, siBaseUnits)
		if !slices.EqualFunc(exponents, unit.Exponents, almostEqual) {
			t.Error("Unit ", unit.Name, "expected to have exponents", unit.Exponents, ", but has", exponents)
		}
	}
}

func TestExpressUnitAsComplexUnits(t *testing.T) {
	mass := Unit{"m", []float64{0, 1, 0}}
	time := Unit{"t", []float64{0, 0, 1}}
	angularMomentum := Unit{"L", []float64{2, 1, -1}}
	baseUnits := []Unit{mass, time, angularMomentum}
	{
		length := Unit{"l", []float64{1, 0, 0}}
		lengthExponents := ExpressUnit(length, baseUnits)
		expectedLengthExponents := []float64{-0.5, 0.5, 0.5}
		if !slices.EqualFunc(lengthExponents, expectedLengthExponents, almostEqual) {
			t.Error("Length expected to have exponents", expectedLengthExponents, ", but has", lengthExponents)
		}
	}
	{
		speed := Unit{"v", []float64{1, 0, -1}}
		speedExponents := ExpressUnit(speed, baseUnits)
		expectedSpeedExponents := []float64{-0.5, -0.5, 0.5}
		if !slices.EqualFunc(speedExponents, expectedSpeedExponents, almostEqual) {
			t.Error("Speed expected to have exponents", expectedSpeedExponents, ", but has", speedExponents)
		}
	}
	{
		momentum := Unit{"p", []float64{1, 1, -1}}
		momentumExponents := ExpressUnit(momentum, baseUnits)
		expectedMomentumExponents := []float64{0.5, -0.5, 0.5}
		if !slices.EqualFunc(momentumExponents, expectedMomentumExponents, almostEqual) {
			t.Error("Momentum expected to have exponents", expectedMomentumExponents, ", but has", momentumExponents)
		}
	}
}

func TestExpressUnitNotIndependent(t *testing.T) {
	length := Unit{"l", []float64{1, 0, 0}}
	mass := Unit{"m", []float64{0, 1, 0}}
	time := Unit{"t", []float64{0, 0, 1}}
	speed := Unit{"v", []float64{1, 0, -1}}
	baseUnits := []Unit{length, time, speed}
	{
		// Acceleration is expressible in terms of the base units, but not uniquely, so we should get nil as result
		acceleration := Unit{"a", []float64{1, 0, -2}}
		accelerationExponents := ExpressUnit(acceleration, baseUnits)
		if accelerationExponents != nil {
			t.Error("Units expected to be reported as not independent, but got exponents", accelerationExponents)
		}
	}
	{
		// Mass is not expressible in terms of the base units, so we should get nil as result
		massExponents := ExpressUnit(mass, baseUnits)
		if massExponents != nil {
			t.Error("Units expected to be reported as not independent, but got exponents", massExponents)
		}
	}
}
