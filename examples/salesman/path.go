package salesman

import (
	"fmt"

	"math/rand"

	"bytes"
	"io"

	"io/ioutil"

	"github.com/tomjcleveland/genetic"
)

// path describes a possible order in which the salesman
// can visit all cities.
type path []*city

func (p path) Crossover(ind genetic.Individual) (genetic.Individual, error) {
	p2, ok := ind.(path)
	if !ok {
		return nil, fmt.Errorf("expecting type genetic.Individual, got %T", ind)
	}
	child := make([]*city, len(p))

	// Add ordered subset of parent 1's path
	added := make(map[city]bool)
	start, end := p.subset()
	for i := start; i < end; i++ {
		child[i] = p[i]
		added[*p[i]] = true
	}

	// Add remainder of cities from parent 2
	for i := 0; i < len(p); i++ {
		p2Index := i + end
		if p2Index >= len(p) {
			p2Index -= len(p)
		}

		if !added[*p2[p2Index]] {
			for j := 0; j < len(child); j++ {
				if child[j] == nil {
					child[j] = p2[p2Index]
					added[*p2[p2Index]] = true
					break
				}
			}
		}
	}
	return path(child), nil
}

func (p path) Mutate(rate float64) (genetic.Individual, error) {
	out := p[:]
	for i := 1; i < len(p); i++ {
		if rate > rand.Float64() {
			out[i-1], out[i] = out[i], out[i-1]
		}
	}
	return path(out), nil
}

func (p path) Fitness() (float64, error) {
	totalDistance := float64(0)
	for i := 0; i < len(p)-1; i++ {
		totalDistance += p[i].distanceFrom(p[i+1])
	}
	return -totalDistance, nil
}

func (p path) subset() (start, end int) {
	a, b := rand.Intn(len(p)), rand.Intn(len(p))
	if a >= b {
		return b, a
	}
	return a, b
}

func (p path) String() string {
	out := bytes.NewBuffer(nil)
	io.WriteString(out, "[")
	for index, city := range p {
		if city == nil {
			fmt.Fprint(out, "<nil>")
		} else {
			fmt.Fprintf(out, "(%d,%d)", city.x, city.y)
		}
		if index != len(p)-1 {
			io.WriteString(out, ",")
		}
	}
	io.WriteString(out, "]")
	outBytes, _ := ioutil.ReadAll(out)
	return string(outBytes)
}
