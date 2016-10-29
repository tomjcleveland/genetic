package genetic

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
)

// Population holds all the individuals in the
// current population, along with their scores.
type Population struct {
	pop []indWithScore
}

// NewPopulation constructs a Population with fitness
// scores of zero.
func NewPopulation(pop interface{}) (*Population, error) {
	out := &Population{}
	val := reflect.ValueOf(pop)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("value has type %q; expecting slice", val.Kind())
	}
	for i := 0; i < val.Len(); i++ {
		if !val.Index(i).CanInterface() {
			return nil, errors.New("input population can't interface{}; something's wrong")
		}
		ind, ok := val.Index(i).Interface().(Individual)
		if !ok {
			return nil, fmt.Errorf("%s does not satisfy interface Individual", val.Type())
		}
		out.pop = append(out.pop, indWithScore{ind, 0})
	}
	if len(out.pop) == 0 {
		return nil, errors.New("cannot pass in empty slice")
	}
	return out, nil
}

// TargetMet returns true if any individual in the population
// has met or exceeded the fitness target.
func (p *Population) TargetMet(t float64) bool {
	for _, ind := range p.pop {
		if ind.score >= t {
			log.Printf("Individual with score %.2f has exceeded target (%.2f)", ind.score, t)
			return true
		}
	}
	return false
}

// Fittest returns the individual with the highest fitness.
func (p *Population) Fittest() (Individual, error) {
	if len(p.pop) == 0 {
		return nil, errors.New("population is empty")
	}
	max := -math.MaxFloat64
	var fittest Individual
	for _, ind := range p.pop {
		score, err := ind.Fitness()
		if err != nil {
			return nil, err
		}
		if score > max {
			max = score
			fittest = ind.Individual
		}
	}
	return fittest, nil
}

// FittestScore returns the fitness score of the fittest individual.
func (p *Population) FittestScore() (float64, error) {
	fittest, err := p.Fittest()
	if err != nil {
		return 0, err
	}
	log.Printf("Fittest: %v", fittest)
	fitness, err := fittest.Fitness()
	if err != nil {
		return 0, err
	}
	return fitness, nil
}

// TotalFitness returns the sum of all the fitness scores
// of the population
func (p *Population) TotalFitness() (float64, error) {
	total := float64(0)
	for _, ind := range p.pop {
		fitness, err := ind.Fitness()
		if err != nil {
			return 0, err
		}
		total += fitness
	}
	return total, nil
}

// AvgFitness returns the average fitness of the population
func (p *Population) AvgFitness() (float64, error) {
	sum, err := p.TotalFitness()
	if err != nil {
		return 0, err
	}
	return sum / float64(len(p.pop)), nil
}

func (p *Population) scoreAndSort(workers int) error {
	newPop, err := calculateFitnessConcurrently(p.pop, workers)
	if err != nil {
		return err
	}
	p.pop = newPop
	return nil
}
