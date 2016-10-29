package genetic

import "math/rand"

// SelectionMethod is used to choose a second Individual
// during crossover.
type SelectionMethod func(Individual, *Population) (Individual, error)

// Roulette is a SelectionMethod that picks a partner for an individual
// at random, weighting the likelihood of picking a particular partner
// with that partner's fitness.
func Roulette(ind Individual, pop *Population) (Individual, error) {
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
