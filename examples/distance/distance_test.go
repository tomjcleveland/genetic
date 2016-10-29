package distance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tomjcleveland/genetic"
	"golang.org/x/net/context"
)

func Test_Distance_Run(t *testing.T) {
	ctrl, err := genetic.NewController(genetic.Params{
		Elitism:          15,
		Mutation:         0.9,
		Crossover:        0.9,
		TargetFitness:    -1,
		Parallelism:      100,
		SelectionMethod:  genetic.Tournament(100),
		InitPop:          testPopulation(500),
		AdaptiveMutation: false,
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
		assert.True(t, -ld(target, string(fittest.(dString))) >= -1)
	}
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
