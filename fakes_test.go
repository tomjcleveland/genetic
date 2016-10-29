package genetic

type fakeIndividual struct {
	id int

	// Stats
	mutateCount    int
	crossoverCount int
	fitnessCount   int

	// Returned values
	fitness float64
	err     error
}

func (fi fakeIndividual) Crossover(ind Individual) (Individual, error) {
	fi.crossoverCount++
	return fi, fi.err
}

func (fi fakeIndividual) Mutate() (Individual, error) {
	fi.mutateCount++
	return fi, fi.err
}

func (fi fakeIndividual) Fitness() (float64, error) {
	fi.fitnessCount++
	return fi.fitness, fi.err
}
