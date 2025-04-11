/*
unit_systems finds the "optimal" set of physical units that can express all
the remaining ones.

For calculating the score for candidate set of units it:
1. Checks if every unit can be expressed as a product of powers of units
   from the candidate set. If it cannot the candidate set is ignored.
2. For each unit counts how many units from the candidate set are required
   to express it, and adds the count to the score for the candidate set
3. Prints the candidates with the minimum score
*/

package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"gonum.org/v1/gonum/stat/combin"
)

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

	var bestUnits [][]string
	var bestScore int
	// Generate all the possible sets of base units and score each of them.
	gen := combin.NewCombinationGenerator(len(units.Units), len(units.SiUnits))
combinations:
	for gen.Next() {
		baseUnitIndices := gen.Combination(nil)
		baseUnits := []Unit{}
		unitNames := []string{}
		for _, ix := range baseUnitIndices {
			unit := units.Units[ix]
			baseUnits = append(baseUnits, unit)
			unitNames = append(unitNames, unit.Name)
		}
		if veryVerbose {
			fmt.Println("Scoring units", unitNames)
		}
		score := 0
		for _, unit := range units.Units {
			result := ExpressUnit(unit, baseUnits)
			if result == nil {
				if verbose {
					fmt.Println("Units", unitNames, "are not independent")
				}
				continue combinations
			}
			nonZeroExponents := 0
			for _, exponent := range result {
				if math.Abs(exponent) > epsilon {
					nonZeroExponents += 1
				}
			}
			score += nonZeroExponents
			if veryVerbose {
				fmt.Print("  +", nonZeroExponents, ", ", unit.Name, " = ", unitNames, " ⋅ ", result, "\n")
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
