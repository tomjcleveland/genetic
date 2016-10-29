package distance

import (
	"testing"

	"github.com/tomjcleveland/kurt/genetic"
)

func Test_Distance_GA(t *testing.T) {
	ctrl, err := genetic.NewController(genetic.Params{
		Elitism:         3,
		Mutation:        0.3,
		Crossover:       0.9,
		TargetFitness:   0,
		SelectionMethod: genetic.Roulette,
		InitPop:         testPopulation(10),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ctrl.Run()
	if err != nil {
		t.Fatal(err)
	}
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
