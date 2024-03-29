package genetic

// Individual is one member of a population
type Individual interface {
	Crossover(Individual) (Individual, error)
	Mutate(float64) (Individual, error)
	Fitness() (float64, error)
}
