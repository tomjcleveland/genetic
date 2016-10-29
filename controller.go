package genetic

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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

// Params holds all of the parameters for the
// genetic algorithm.
type Params struct {
	// 0-len(population); How many of the top Individuals in the current
	// population should make it into the next generation unchanged?
	Elitism int

	// 0-1; What portion of the surviving population should undergo mutation?
	Mutation float32

	// 0-1; What portion of the surviving population should undergo crossover?
	Crossover float32

	// What is the target fitness score, at which point the
	// algorithm will terminate?
	TargetFitness float64

	// SelectionMethod is the method by which the genetic algorithm
	// chooses a partner for crossover.
	SelectionMethod SelectionMethod

	// The initial population. Must be a slice of Individuals.
	InitPop interface{}
}

type indWithScore struct {
	Individual
	score float64
}

type pairs []indWithScore

func (p pairs) Len() int           { return len(p) }
func (p pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairs) Less(i, j int) bool { return p[i].score < p[j].score }

// Controller coordinates the running of the
// genetic algorithm.
type Controller struct {
	params     Params
	population *Population
}

// NewController is the constructor for Controller. It returns an error
// if any of the input parameters are not allowed.
func NewController(params Params) (*Controller, error) {
	if !isFrac(params.Crossover) {
		return nil, errors.New("crossover factor must be between 0 and 1")
	}
	if params.Elitism < 0 {
		return nil, errors.New("elitism factor must be greater than 0")
	}
	if !isFrac(params.Mutation) {
		return nil, errors.New("mutation factor must be between 0 and 1")
	}
	if params.SelectionMethod == nil {
		return nil, errors.New("selection method cannot be nil")
	}
	pop, err := NewPopulation(params.InitPop)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize population: %s", err)
	}
	return &Controller{
		params:     params,
		population: pop,
	}, nil
}

func isFrac(d float32) bool {
	return d >= 0 && d <= 1
}

// Run runs the genetic algorithm until an individual with
// the target fitness level is found. It returns only after
// finding this individual.
func (c *Controller) Run() error {
	// Score initial population
	err := c.population.scoreAndSort()
	if err != nil {
		return err
	}

	// Loop through generations until target fitness is acheived.
	for !c.population.TargetMet(c.params.TargetFitness) {
		_, err := c.population.FittestScore()
		if err != nil {
			return err
		}
		if err := c.performCrossovers(); err != nil {
			return fmt.Errorf("crossover step failed: %s", err)
		}
		if err := c.performMutations(); err != nil {
			return fmt.Errorf("mutation step failed: %s", err)
		}
		err = c.population.scoreAndSort()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) performCrossovers() error {
	var err error
	newGeneration := make([]Individual, len(c.population.pop))
	for i, ind := range c.population.pop {

		// Should we perform crossover on this individual?
		if c.params.Crossover > rand.Float32() && i >= c.params.Elitism {
			parent, err := c.params.SelectionMethod(ind.Individual, c.population)
			if err != nil {
				return err
			}
			child, err := ind.Crossover(parent)
			if err != nil {
				return err
			}
			newGeneration[i] = child
		} else {
			// If not, it goes to the next generation unchanged
			newGeneration[i] = ind.Individual
		}

	}
	c.population, err = NewPopulation(newGeneration)
	return err
}

func (c *Controller) performMutations() error {
	var err error
	newGeneration := make([]Individual, len(c.population.pop))
	for i, ind := range c.population.pop {

		// Should we perform mutation on this individual?
		if c.params.Mutation > rand.Float32() && i >= c.params.Elitism {
			mutated, err := ind.Mutate()
			if err != nil {
				return err
			}
			newGeneration[i] = mutated
		} else {
			// If not, it goes to the next generation unchanged
			newGeneration[i] = ind.Individual
		}

	}
	c.population, err = NewPopulation(newGeneration)
	return err
}
