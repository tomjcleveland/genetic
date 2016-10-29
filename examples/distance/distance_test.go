package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomjcleveland/kurt/genetic"
)

func Test_Distance_Run(t *testing.T) {
	ctrl, err := genetic.NewController(genetic.Params{
		Elitism:         15,
		Mutation:        0.3,
		Crossover:       0.9,
		TargetFitness:   -1,
		SelectionMethod: genetic.Roulette,
		InitPop:         testPopulation(100),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ctrl.Run()
	if err != nil {
		t.Fatal(err)
	}

	fittest, err := ctrl.Fittest()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, -ld(target, string(fittest.(dString))) >= -1)
}

func testPopulation(n int) []genetic.Individual {
	var popStrings []string
	for i := 0; i < n; i++ {
		curr := ""
		for j := 0; j < len(target); j++ {
			curr += string(randCharacter())
		}
		popStrings = append(popStrings, curr)
	}

	var out []genetic.Individual
	for _, i := range popStrings {
		ds := dString(i)
		out = append(out, ds)
	}
	return out
}
