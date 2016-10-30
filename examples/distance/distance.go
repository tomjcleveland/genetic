package distance

import (
	"fmt"
	"math/rand"

	"github.com/tomjcleveland/genetic"
)

// target is the string we're trying to get the genetic algorithm
// to reproduce.
var target = "this is the target string"

// dString implements Individual
type dString string

func (d dString) Mutate(rate float64) (genetic.Individual, error) {
	stringBytes := []byte(d)
	for i := 0; i < len(d); i++ {
		if rand.Float64() < rate {
			stringBytes[i] = randCharacter()
		}
	}

	return dString(stringBytes), nil
}

func (d dString) Crossover(i genetic.Individual) (genetic.Individual, error) {
	child := dString("")
	partner, ok := i.(dString)
	if !ok {
		return &child, fmt.Errorf("expected Individual to be *dString, got %T", i)
	}
	for i := 0; i < len(d); i++ {
		if rand.Intn(2) > 0 {
			child += dString((d)[i])
		} else {
			child += dString((partner)[i])
		}
	}
	return child, nil
}

func (d dString) Fitness() (float64, error) {
	return float64(-ld(target, string(d))), nil
}

func ld(s, t string) int {
	d := make([][]int, len(s)+1)
	for i := range d {
		d[i] = make([]int, len(t)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}

	}
	return d[len(s)][len(t)]
}

func (d *dString) String() string {
	return string(*d)
}

var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 ?!.&%^"

func randCharacter() byte {
	return chars[rand.Intn(len(chars))]
}
