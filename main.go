/*
unit_systems finds the "optimal" set of physical units that can express all
the remaining ones.

For calculating the score for candidate set of units it:
1. Checks if every unit can be expressed as a linear combination of powers
   of units from the candidate set.
   If it cannot the candidate set is ignored.
2. For each unit counts how many units from the candidate set are required
   to express it, and adds the count to the score for the candidate set
3. Prints the candidates with the minimum score
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

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

var inputFlag = flag.String("i", "", "path to json file with description of units")
var verboseFlag = flag.Bool("v", false, "print score for every candidate set")
var veryVerboseFlag = flag.Bool("vv", false, "print partial score for every candidate set")

const epsilon = 0.0001

func main() {
	flag.Parse()

	if *inputFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	veryVerbose := *veryVerboseFlag
	verbose := veryVerbose || *verboseFlag

	units := ReadUnitsFromFile(*inputFlag)
	siUnitCount := len(units.SiUnits)

	var bestUnits [][]string
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
		// For every unit, try to express it in terms of the candidate base units.
		if veryVerbose {
			fmt.Println("Scoring units", unitNames)
		}
		for _, unit := range units.Units {
			coeffs := mat.NewDense(siUnitCount, 1, unit.Coefficients)
			resultSlice := make([]float64, siUnitCount)
			result := mat.NewDense(siUnitCount, 1, resultSlice)
			// Solve the linear system: m * x = coeffs.
			// If the system does not have a solution, the candidate base units are not independent.
			if err := result.Solve(m, coeffs); err != nil {
				if verbose {
					fmt.Println("Units", unitNames, "are not independent")
				}
				continue combinations
			}
			partialScore := 0
			for _, resultCoeff := range resultSlice {
				if math.Abs(resultCoeff) > epsilon {
					partialScore += 1
				}
			}
			score += partialScore
			if veryVerbose {
				fmt.Print("  +", partialScore, ", ", unit.Name, " = ", unitNames, " ⋅ ", resultSlice, "\n")
			}
		}
		if verbose {
			fmt.Println("Score for units", unitNames, "is", score)
		}
		if bestUnits == nil || score < bestScore {
			bestUnits = [][]string{unitNames}
			bestScore = score
		} else if score == bestScore {
			bestUnits = append(bestUnits, unitNames)
		}
	}
	fmt.Println("Best units:", bestUnits, "score:", bestScore)
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
		if len(unit.Coefficients) != siUnitCount {
			log.Fatalf("unit %s has should have %d coefficients, but has %d", unit.Name, siUnitCount, len(unit.Coefficients))
		}
	}
	if len(units.Units) < siUnitCount {
		log.Fatalf("there should be at least as many units as SI units (%d), but there are %d", siUnitCount, len(units.Units))
	}

	return units
}
