package genetic

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"context"
)

// ErrContextCancelled is returned by Run() or Wait() when the search
// has been prematurely cancelled by the context before an acceptable
// solution was found.
var ErrContextCancelled = errors.New("search cancelled by context")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Params holds all of the parameters for the
// genetic algorithm.
type Params struct {
	// 0-len(population); How many of the top Individuals in the current
	// population should make it into the next generation unchanged?
	Elitism int

	// 0-1; What portion of the population should undergo mutation?
	Mutation float64

	// Should we change the mutation rate based on search convergence?
	AdaptiveMutation bool

	// 0-1; What portion of the population should undergo crossover?
	Crossover float64

	// What is the target fitness score, at which point the
	// algorithm will terminate?
	TargetFitness float64

	// Parallelism dictates how many goroutines will be used to calculate
	// the fitness of a population. The default is one.
	Parallelism int

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

	err chan error
}

// NewController is the constructor for Controller. It returns an error
// if any of the input parameters are not allowed.
func NewController(params Params) (*Controller, error) {
	if !isProb(params.Crossover) {
		return nil, errors.New("crossover factor must be between 0 and 1, inclusive")
	}
	if params.Elitism < 0 {
		return nil, errors.New("elitism factor must be greater than 0")
	}
	if !isProb(params.Mutation) {
		return nil, errors.New("mutation factor must be between 0 and 1, inclusive")
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
		err:        make(chan error),
	}, nil
}

func isProb(d float64) bool {
	return d >= 0 && d <= 1
}

func (c *Controller) run(ctx context.Context) error {
	// Score initial population
	err := c.population.scoreAndSort(c.params.Parallelism)
	if err != nil {
		return err
	}

	// Loop through generations until target fitness is acheived.
	for !c.population.TargetMet(c.params.TargetFitness) {
		select {
		case <-ctx.Done():
			return ErrContextCancelled
		default:
		}
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
		err = c.population.scoreAndSort(c.params.Parallelism)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run runs the genetic algorithm until an individual with
// the target fitness level is found. It returns only after
// finding this individual.
func (c *Controller) Run() error {
	c.Start(context.Background())
	return c.Wait()
}

// Start begins the genetic algorithm in a new goroutine, and
// returns immediately. The context parameter can be used to
// prematurely cancel a long-running search.
func (c *Controller) Start(ctx context.Context) {
	go func() {
		c.err <- c.run(ctx)
	}()
}

// Wait blocks until the target fitness has been acheived.
func (c *Controller) Wait() error {
	if err, ok := <-c.err; ok {
		return err
	}
	return errors.New("error channel is closed")
}

// Fittest returns the fittest individual in the current population.
func (c *Controller) Fittest() (Individual, error) {
	return c.population.Fittest()
}

func (c *Controller) performCrossovers() error {
	var err error
	newGeneration := make([]Individual, len(c.population.pop))
	for i, ind := range c.population.pop {

		// Should we perform crossover on this individual?
		if c.params.Crossover > rand.Float64() && i >= c.params.Elitism {
			parent, err := c.params.SelectionMethod(c.population)
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
		mutatationRate := c.params.Mutation

		// Are we using an adaptive mutation rate?
		if c.params.AdaptiveMutation {
			currScore, err := ind.Fitness()
			if err != nil {
				return err
			}
			avgFitness, err := c.population.AvgFitness()
			if err != nil {
				return err
			}
			if currScore > avgFitness {
				fittestScore, err := c.population.FittestScore()
				if err != nil {
					return err
				}

				delta1 := fittestScore - currScore
				delta2 := fittestScore - avgFitness
				if delta2 == 0 {
					mutatationRate = 1 // If the average and the fittest are the same, we need some mutation
				} else {
					mutatationRate = delta1 / delta2 * c.params.Mutation
				}
			}
		}

		// Should we perform mutation on this individual?
		if mutatationRate > rand.Float64() && i >= c.params.Elitism {
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
