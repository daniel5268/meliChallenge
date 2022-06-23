package geometry

import "math"

// Coordinate describes the components of a coordinate in a cartesian plane.
type Coordinate struct {
	X float64
	Y float64
}

// NewCoordinate returns a new Coordinate given x and y.
func NewCoordinate(x, y float64) Coordinate {
	return Coordinate{
		X: x,
		Y: y,
	}
}

// Diff returns the difference between this coordinate and the provided one
func (c Coordinate) Diff(c2 Coordinate) Coordinate {
	return NewCoordinate(c.X-c2.X, c.Y-c2.Y)
}

func (c Coordinate) Distance(c2 Coordinate) float64 {
	return math.Sqrt(math.Pow(c.X-c2.X, 2) + math.Pow(c.Y-c2.Y, 2))
}
