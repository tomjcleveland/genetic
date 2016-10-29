package genetic

import (
	"errors"
	"math"
	"math/rand"
)

// SelectionMethod is used to choose a second Individual
// during crossover.
type SelectionMethod func(*Population) (Individual, error)

// Roulette returns a SelectionMethod that picks a partner for an individual
// at random, weighting the likelihood of picking a particular partner
// with that partner's fitness.
func Roulette() SelectionMethod {
	return func(pop *Population) (Individual, error) {
		totalFitness, err := pop.TotalFitness()
		if err != nil {
			return nil, err
		}
		position := rand.Float64() * totalFitness
		spinWheel := float64(0)
		for _, curr := range pop.pop {
			fitness, err := curr.Fitness()
			if err != nil {
				return nil, err
			}
			spinWheel += fitness
			if spinWheel >= position {
				return curr.Individual, nil
			}
		}
		return pop.pop[len(pop.pop)-1].Individual, nil
	}
}

// Tournament returns the tournament SelectionMethod, which uses
// tournaments of size n to select a partner for crossover.
func Tournament(n int) SelectionMethod {
	return func(pop *Population) (Individual, error) {
		if len(pop.pop) < n {
			return nil, errors.New("tournament size is larger than population")
		}
		var winner Individual
		max := -math.MaxFloat64
		for i := 0; i < n; i++ {
			if pop.pop[i].score > max {
				winner = pop.pop[i].Individual
				max = pop.pop[i].score
			}
		}
		return winner, nil
	}
}
