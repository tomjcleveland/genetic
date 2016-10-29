package genetic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPopulation_InputNotSlice_ErrNotNil(t *testing.T) {
	_, err := NewPopulation("not a slice")

	assert.NotNil(t, err)
}

func Test_NewPopulation_InputSliceNotIndividuals_ErrNotNil(t *testing.T) {
	_, err := NewPopulation([]string{"these", "aren't", "individuals"})

	assert.NotNil(t, err)
}

func Test_NewPopulation_InputSliceHasZeroLength_ErrNotNil(t *testing.T) {
	_, err := NewPopulation([]string{})

	assert.NotNil(t, err)
}

func Test_NewPopulation_InputNil_ErrNotNil(t *testing.T) {
	_, err := NewPopulation(nil)

	assert.NotNil(t, err)
}

func Test_NewPopulation_ValidInput_ErrNil(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	_, err := NewPopulation(popIn)

	assert.Nil(t, err)
}

func Test_NewPopulation_ValidInput_InputLengthMatchesOutputLength(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, pop.pop, 7)
}

func Test_TargetMet_PopulationBelowTarget_ReturnFalse(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, pop.TargetMet(1))
}

func Test_TargetMet_IndividualAboveTarget_ReturnTrue(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	pop.pop = append(pop.pop, indWithScore{
		Individual: fakeIndividual{},
		score:      2,
	})

	assert.True(t, pop.TargetMet(1))
}

func Test_TargetMet_IndividualEqualsTarget_ReturnTrue(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	pop.pop = append(pop.pop, indWithScore{
		Individual: fakeIndividual{},
		score:      1,
	})

	assert.True(t, pop.TargetMet(1))
}

func Test_Fittest_FittestIndividualReturned(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	fittest := indWithScore{score: 1, Individual: fakeIndividual{id: 3}}
	pop.pop = append(pop.pop, fittest)

	ret, err := pop.Fittest()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fittest.Individual, ret)
}

func Test_FittestScore_FittestScoreReturned(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	fittest := indWithScore{score: 1, Individual: fakeIndividual{id: 3, fitness: 1}}
	pop.pop = append(pop.pop, fittest)

	ret, err := pop.FittestScore()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, float64(1), ret)
}

func Test_TotalFitness_SumOfFitnessesReturned(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	pop.pop = append(pop.pop, indWithScore{score: 1, Individual: fakeIndividual{id: 3, fitness: 1}})
	pop.pop = append(pop.pop, indWithScore{score: 5, Individual: fakeIndividual{id: 3, fitness: 5}})

	ret, err := pop.TotalFitness()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, float64(6), ret)
}

func Test_scoreAndSort_PopulationScoredAndSorted(t *testing.T) {
	popIn := make([]fakeIndividual, 7)
	pop, err := NewPopulation(popIn)
	if err != nil {
		t.Fatal(err)
	}
	id3 := indWithScore{score: 5, Individual: fakeIndividual{id: 3, fitness: 5}}
	id4 := indWithScore{score: 1, Individual: fakeIndividual{id: 4, fitness: 1}}
	pop.pop = append(pop.pop, id3, id4)

	err = pop.scoreAndSort()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, id3, pop.pop[0])
	assert.Equal(t, id4, pop.pop[1])
	assert.Equal(t, indWithScore{Individual: fakeIndividual{}}, pop.pop[2])
}
