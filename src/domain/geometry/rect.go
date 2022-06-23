package geometry

import "math"

// Describes the components of a rect equation,
// as a consequence of not being able to describe an infinity slope
// isVertical and xVertical properties were added.
type RectEquation struct {
	m          float64 // slope
	b          float64 // cut point in Y
	isVertical bool    // true if the rect is vertical
	xVertical  float64 // cut point with X when the rect is vertical
}

// NewRectEquation returns a new RectEquation given two coordinates contained by it
func NewRectEquation(c1, c2 Coordinate) RectEquation {
	mDenominator := c1.X - c2.X
	if mDenominator == 0 {
		return RectEquation{
			isVertical: true,
			xVertical:  c1.X,
		}
	}
	m := (c1.Y - c2.Y) / mDenominator
	b := c1.Y - m*c1.X

	return RectEquation{
		m: m,
		b: b,
	}
}

// DistanceToCoordinate returns the minimun distance fron a given coordinate to the rect described by the equation
func (re RectEquation) DistanceToCoordinate(c Coordinate) float64 {
	if re.isVertical {
		return math.Abs(re.xVertical - c.X)
	}

	return math.Abs(re.m*c.X-c.Y+re.b) / math.Sqrt(re.m*re.m+1)
}
