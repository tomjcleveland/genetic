package salesman

import "math"

// city describes a city that the traveling salesman must visit.
// Two cities with the same coordinates are identical.
type city struct {
	x int
	y int
}

func (c *city) distanceFrom(other *city) float64 {
	delta1 := math.Pow(float64(other.x-c.x), 2)
	delta2 := math.Pow(float64(other.y-c.y), 2)
	return math.Sqrt(delta1 + delta2)
}
