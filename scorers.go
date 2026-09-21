package main

import "math/big"

type Scorer interface {
	Add(exponents []*big.Rat)
	Total() float64
	Reset()
}

type countScorer struct {
	count int
}

func NewCountScorer() Scorer {
	return &countScorer{}
}

func (s *countScorer) Add(exponents []*big.Rat) {
	for _, exponent := range exponents {
		if exponent.Num().Sign() != 0 {
			s.count += 1
		}
	}
}

func (s *countScorer) Total() float64 {
	return float64(s.count)
}

func (s *countScorer) Reset() {
	s.count = 0
}

type sumScorer struct {
	sum *big.Rat
}

func NewSumScorer() Scorer {
	scorer := &sumScorer{}
	scorer.sum = new(big.Rat)
	return scorer
}

func (s *sumScorer) Add(exponents []*big.Rat) {
	abs := new(big.Rat)
	for _, exponent := range exponents {
		s.sum.Add(s.sum, abs.Abs(exponent))
	}
}

func (s *sumScorer) Total() float64 {
	result, _ := s.sum.Float64()
	return result
}

func (s *sumScorer) Reset() {
	s.sum.SetInt64(0)
}
