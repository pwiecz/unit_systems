package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
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

type ratMat [][]*big.Rat

func newMatRat(rows, cols int) ratMat {
	mat := make([][]*big.Rat, rows)
	for i := range rows {
		mat[i] = make([]*big.Rat, cols)
		for j := range cols {
			mat[i][j] = new(big.Rat)
		}
	}
	return mat
}

// Gauss elimination for matrices of big.Rat
func solve(mat ratMat) ([]*big.Rat, error) {
	rows := len(mat)
	if rows == 0 {
		log.Fatal("Empty matrix")
	}
	cols := len(mat[0])
	if cols != rows+1 {
		log.Fatalf("Expected column count (%d) to be 1 greater than the row count (%d)", cols, rows)
	}
	for i := range rows {
		pivotRow := i
		for r := i + 1; r < rows; r++ {
			absR := new(big.Rat).Abs(mat[r][i])
			absPivot := new(big.Rat).Abs(mat[pivotRow][i])
			if absR.Cmp(absPivot) > 0 {
				pivotRow = r
			}
		}

		if mat[pivotRow][i].Sign() == 0 {
			return nil, fmt.Errorf("no single solution")
		}

		if pivotRow != i {
			mat[i], mat[pivotRow] = mat[pivotRow], mat[i]
		}

		pivotVal := new(big.Rat).Set(mat[i][i])
		for j := i; j < cols; j++ {
			mat[i][j].Quo(mat[i][j], pivotVal)
		}

		for r := range rows {
			if r != i && mat[r][i].Sign() != 0 {
				factor := new(big.Rat).Set(mat[r][i])
				for j := i; j < cols; j++ {
					term := new(big.Rat).Mul(factor, mat[i][j])
					mat[r][j].Sub(mat[r][j], term)
				}
			}
		}
	}

	solution := make([]*big.Rat, rows)
	for i := range rows {
		solution[i] = mat[i][cols-1]
	}

	return solution, nil
}

// ExpressUnit tries to express given physical unit as a product of powers of units from
// the baseUnits.
// If it succeeds it returns a slice which for each corresponding base unit contains
// exponent it has in the representation.
// If it cannot express the unit in terms of baseUnits it returns nil.
// All the passed units have to use the SI unit system (same number of exponents),
// and number of base units have to be equal to the number of SI units.
func ExpressUnit(unit Unit, baseUnits []Unit) []*big.Rat {
	siUnitCount := len(unit.Exponents)
	if len(baseUnits) != siUnitCount {
		log.Fatal("Expected", siUnitCount, "base units, got", len(baseUnits))
	}
	m := newMatRat(siUnitCount, siUnitCount+1)
	for i, baseUnit := range baseUnits {
		if len(baseUnit.Exponents) != siUnitCount {
			log.Fatal("Expected", siUnitCount, "base units for unit", baseUnit.Name, ", got", len(baseUnit.Exponents))
		}
		for j, coeff := range baseUnit.Exponents {
			m[j][i].SetFloat64(coeff)
		}
	}
	for i, siUnit := range unit.Exponents {
		m[i][siUnitCount].SetFloat64(siUnit)
	}

	// Solve the linear system: m * x = coeffs.
	// If the system does not have a solution, the candidate base units are not independent.
	if result, err := solve(m); err != nil {
		return nil
	} else {
		return result
	}
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
