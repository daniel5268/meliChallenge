package geometry

import "math"

const (
	COMPLETE_SPIN_DEGREES = float64(360)
	HALF_SPIN_DEGREES     = float64(180)
)

// Returns the radians equivalent to the provided degrees.
func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / HALF_SPIN_DEGREES
}

// GetTraveledDistance returns an entitie's traveled distance given it's speed and passed time.
func GetTraveledDistance(speed, time float64) float64 {
	return speed * time
}

// IsOriginInTriangle calculates if the triangle formed by the provided Coordinates contains the origin,
// the used algorithm is described here https://www.youtube.com/watch?v=WaYS1gEXEFE&t=542s&ab_channel=huse360.
func IsOriginInTriangle(a, b, c Coordinate) bool {
	d := b.Diff(a)
	e := c.Diff(a)

	w1 := (e.X*a.Y - e.Y*a.X) / (d.X*e.Y - d.Y*e.X)
	w2 := -(a.Y + w1*d.Y) / e.Y

	return w1 >= 0 && w2 >= 0 && w1+w2 <= 1
}

// GetTrianglePerimeter returns a triangle perimeter given it's Coordinates
func GetTrianglePerimeter(a Coordinate, b Coordinate, c Coordinate) float64 {
	return a.Distance(b) + a.Distance(c) + b.Distance(c)
}
