# github.com/tomjcleveland/genetic
This package simplifies the creation of [genetic algorithms](https://en.wikipedia.org/wiki/Genetic_algorithm). To use this package, you must provide an implementation of `genetic.Individual`, which has the following interface:
```go
// Individual is one member of a population
type Individual interface {
	Crossover(Individual) (Individual, error)
	Mutate(float64) (Individual, error)
	Fitness() (float64, error)
}
```
In a genetic algorithm, an *individual* is a candidate solution. The `Crossover()` and `Mutate()` methods produce *offspring* that are minor variations of the *parents*. In this way, the genetic algorithm traverses the search space, keeping inividuals with high fitness (as determined by `Fitness()`) and discarding those with low fitness.

## Quickstart
```go
package main

import (
    "log"

    "github.com/tomjcleveland/genetic"
)

func main() {
    ctrl, _ := genetic.NewController(genetic.Params{
        Elitism:          3,
        Mutation:         0.5,
        Crossover:        0.7,
        TargetFitness:    30,
        Parallelism:      10,
        SelectionMethod:  genetic.Tournament(10),
        InitPop:          testPopulation(50),
        AdaptiveMutation: true,
    })
    ctrl.Run()
    fittest, _ := ctrl.Fittest()
    log.Printf("Best solution: %v", fittest)
}

func testPopulation(n int) []genetic.Individual {
    // This is the tricky part: implementing genetic.Individual
    // See examples for inspiration
}
```

[API Documentation (GoDoc)](https://godoc.org/github.com/tomjcleveland/genetic)

## Examples

### [String Distance](examples/distance)
In this simple example, the genetic algorithm starts with randomly generated strings, and searches for strings that are similar to the target string.

### [Traveling Salesman](examples/salesman)
This example searches for solutions to the [traveling salesman problem](https://en.wikipedia.org/wiki/Travelling_salesman_problem), where the goal is to find the shortest possible between a set of points in 2D space.