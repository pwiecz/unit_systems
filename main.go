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

type Unit struct {
	Name         string `json:"name"`
	Coefficients []float64  `json:"coeffs"`
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
			log.Fatal("unit %s has should have %d coefficients, but has %d", unit.Name, siUnitCount, len(unit.Coefficients))
		}
	}
	if len(units.Units) < siUnitCount {
		log.Fatal("there should be at least as many units as SI units (%d), but there are %d", siUnitCount, len(units.Units))
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
		for i, unit := range units.Units {
			if slices.Contains(baseUnits, i) {
				continue
			}
			coeffs := mat.NewDense(siUnitCount, 1, unit.Coefficients)
			resultSlice := make([]float64, 3)
			result := mat.NewDense(siUnitCount, 1, resultSlice)
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
