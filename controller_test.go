package genetic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewController_TestTable_AllInvalidParams(t *testing.T) {
	testTable := []struct {
		label  string
		params Params
	}{
		{
			label: "negative mutation factor",
			params: Params{
				Mutation:        -0.5,
				InitPop:         make([]fakeIndividual, 3),
				SelectionMethod: Roulette(),
			},
		},
		{
			label: "out of bounds mutation factor",
			params: Params{
				Mutation:        30,
				InitPop:         make([]fakeIndividual, 3),
				SelectionMethod: Roulette(),
			},
		},
		{
			label: "negative crossover factor",
			params: Params{
				Crossover:       -0.5,
				InitPop:         make([]fakeIndividual, 3),
				SelectionMethod: Roulette(),
			},
		},
		{
			label: "out of bounds crossover factor",
			params: Params{
				Crossover:       30,
				InitPop:         make([]fakeIndividual, 3),
				SelectionMethod: Roulette(),
			},
		},
		{
			label: "negative elitism count",
			params: Params{
				Elitism:         -10,
				InitPop:         make([]fakeIndividual, 3),
				SelectionMethod: Roulette(),
			},
		},
		{
			label: "selection method nil",
			params: Params{
				InitPop: make([]fakeIndividual, 3),
			},
		},
		{
			label: "InitPop is not slice",
			params: Params{
				InitPop:         "not a slice",
				SelectionMethod: Roulette(),
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.label, func(t *testing.T) {
			_, err := NewController(testCase.params)
			assert.NotNil(t, err)
		})
	}
}

func Test_Run_PopWithTargetMet_ErrNilAndCorrectWinner(t *testing.T) {
	fittest := fakeIndividual{id: 4, fitness: 5}
	pop := []fakeIndividual{
		{
			id:      0,
			fitness: 0,
		},
		{
			id:      1,
			fitness: 1,
		},
		{
			id:      2,
			fitness: 2,
		},
		{
			id:      3,
			fitness: 3,
		},
		fittest,
	}
	ctrl, err := NewController(Params{
		TargetFitness:   4,
		SelectionMethod: Roulette(),
		InitPop:         pop,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ctrl.Run()
	if err != nil {
		t.Fatal(err)
	}

	ret, err := ctrl.Fittest()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fittest, ret)
}
