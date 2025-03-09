package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"slices"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/combin"
)

// Unit represents a physical unit, with a name and a list of coefficients.
// These coefficients express the unit as a combination of SI base units.
type Unit struct {
	Name         string    `json:"name"`
	Coefficients []float64 `json:"coeffs"`
}

type Units struct {
	SiUnits []string `json:"si_units"`
	Units   []Unit   `json:"units"`
}

var inputFlag = flag.String("i", "", "description of units")

func main() {
	flag.Parse()

	if *inputFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
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
		if len(unit.Coefficients) != siUnitCount {
			log.Fatalf("unit %s has should have %d coefficients, but has %d", unit.Name, siUnitCount, len(unit.Coefficients))
		}
	}
	if len(units.Units) < siUnitCount {
		log.Fatalf("there should be at least as many units as SI units (%d), but there are %d", siUnitCount, len(units.Units))
	}

	var bestUnits []string
	var bestScore int
	gen := combin.NewCombinationGenerator(len(units.Units), siUnitCount)
combinations:
	for gen.Next() {
		score := 0
		baseUnits := gen.Combination(nil)
		unitNames := make([]string, siUnitCount)
		m := mat.NewDense(siUnitCount, siUnitCount, nil)
		for i, ix := range baseUnits {
			unitNames[i] = units.Units[ix].Name
			for j, coeff := range units.Units[ix].Coefficients {
				m.Set(j, i, coeff)
			}
		}
		// For every unit not in the candidate base set, try to express it in terms of the candidate base units.
		for i, unit := range units.Units {
			if slices.Contains(baseUnits, i) {
				continue
			}
			coeffs := mat.NewDense(siUnitCount, 1, unit.Coefficients)
			resultSlice := make([]float64, siUnitCount)
			result := mat.NewDense(siUnitCount, 1, resultSlice)
			// Solve the linear system: m * x = coeffs.
			// If the system does not have a solution, the candidate base units are not independent.
			if err := result.Solve(m, coeffs); err != nil {
				log.Print("Units ", unitNames, " are not independent")
				continue combinations
			}
			for _, resultCoeff := range resultSlice {
				if resultCoeff != 0 {
					score += 1
				}
			}
		}
		log.Print("Score for units ", unitNames, " is ", score)
		if bestUnits == nil || score < bestScore {
			bestUnits = unitNames
			bestScore = score
		}
	}
	log.Print("Best units: ", bestUnits, " score: ", bestScore)
}
