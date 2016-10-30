package salesman

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"math/rand"

	"github.com/tomjcleveland/genetic"
)

func Test_Salesman_Run(t *testing.T) {
	ctrl, err := genetic.NewController(genetic.Params{
		Elitism:          3,
		Mutation:         0.8,
		Crossover:        0.7,
		TargetFitness:    10,
		Parallelism:      10,
		SelectionMethod:  genetic.Tournament(10),
		InitPop:          testPopulation(50),
		AdaptiveMutation: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	ctrl.Start(ctx)
	done := make(chan error)
	go func() {
		done <- ctrl.Wait()
	}()

	select {
	case <-time.After(time.Second * 5):
		cancelFunc()
		t.Fatal("search did not finish after five seconds")
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
		fittest, err := ctrl.Fittest()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Fittest: %+v", fittest)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func testPopulation(n int) []genetic.Individual {
	mapScale := 50
	numCities := 10
	paths := make([]genetic.Individual, n)

	for i := 0; i < n; i++ {
		seen := make(map[city]bool)
		cities := make([]*city, numCities)
		for j := 0; j < numCities; j++ {
			curr := &city{
				x: rand.Intn(mapScale),
				y: rand.Intn(mapScale),
			}
			for seen[*curr] {
				curr = &city{
					x: rand.Intn(mapScale),
					y: rand.Intn(mapScale),
				}
			}
			seen[*curr] = true
			cities[j] = curr
		}
		paths[i] = genetic.Individual(path(cities))
	}

	return paths
}
