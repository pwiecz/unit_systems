package main

import (
	"encoding/json"
	"log"
	"os"

	"gonum.org/v1/gonum/mat"
)

// Unit represents a physical unit, with a name and a list of exponents.
// These exponents represent the unit as a product of powers of SI base units.
type Unit struct {
	Name      string    `json:"name"`
	Exponents []float64 `json:"exponents"`
}

type Units struct {
	SiUnits []string `json:"si_units"`
	Units   []Unit   `json:"units"`
}

// ExpressUnit tries to express given physical unit as a product of powers of units from
// the baseUnits.
// If it succeeds it returns a slice which for each corresponding base unit contains
// exponent it has in the representation.
// If it cannot express the unit in terms of baseUnits it returns nil.
// All the passed units have to use the SI unit system (same number of exponents),
// and number of base units have to be equal to the number of SI units.
func ExpressUnit(unit Unit, baseUnits []Unit) []float64 {
	siUnitCount := len(unit.Exponents)
	if len(baseUnits) != siUnitCount {
		log.Fatal("Expected", siUnitCount, "base units, got", len(baseUnits))
	}
	m := mat.NewDense(siUnitCount, siUnitCount, nil)
	for i, baseUnit := range baseUnits {
		if len(baseUnit.Exponents) != siUnitCount {
			log.Fatal("Expected", siUnitCount, "base units for unit", baseUnit.Name, ", got", len(baseUnit.Exponents))
		}
		for j, coeff := range baseUnit.Exponents {
			m.Set(j, i, coeff)
		}
	}

	coeffs := mat.NewDense(siUnitCount, 1, unit.Exponents)
	resultSlice := make([]float64, siUnitCount)
	result := mat.NewDense(siUnitCount, 1, resultSlice)
	// Solve the linear system: m * x = coeffs.
	// If the system does not have a solution, the candidate base units are not independent.
	if err := result.Solve(m, coeffs); err != nil {
		return nil
	}
	return resultSlice
}

func ReadUnitsFromFile(path string) Units {
	f, err := os.Open(*inputFlag)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(f)
	var units Units
	if err = decoder.Decode(&units); err != nil {
		log.Fatal(err)
	}

	siUnitCount := len(units.SiUnits)
	for _, unit := range units.Units {
		if len(unit.Exponents) != siUnitCount {
			log.Fatalf("unit %s has should have %d exponents, but has %d", unit.Name, siUnitCount, len(unit.Exponents))
		}
	}
	if len(units.Units) < siUnitCount {
		log.Fatalf("there should be at least as many units as SI units (%d), but there are %d", siUnitCount, len(units.Units))
	}

	return units
}
