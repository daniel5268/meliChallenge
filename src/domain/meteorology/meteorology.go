package meteorology

import (
	"math"

	"github.com/daniel5268/meliChallenge/src/domain/geometry"
)

const (
	ALIGNMENT_THRESHOLD = float64(0.01) // closeness allowed to the main line for a coordinate to be considered aligned
)

type alignment int16

const (
	alignmentWithoutSun alignment = iota
	alignmentWithSun
	alignmentNone
)

const (
	ClimateFollow = "follow"
	ClimateRain   = "rain"
	ClimateIdeal  = "ideal"
	ClimateNoInfo = "no_info"
)

var sunCoordinate geometry.Coordinate = geometry.NewCoordinate(0, 0)

// Planet conatins the properties needed to calculate a planets position
type Planet struct {
	DegreesByDay   float64
	InitialDegrees float64
	DistanceToSun  uint64
}

// NewPlanet returns a new Planet given it's properties
func NewPlanet(degreesByDay float64, initialDegrees float64, distanceToSun uint64) Planet {
	return Planet{
		DegreesByDay:   degreesByDay,
		InitialDegrees: initialDegrees,
		DistanceToSun:  distanceToSun,
	}
}

// getPositionDegrees returns the degrees traveled by the planet starting from zero everytime a full spin is passed
func (p Planet) getPositionDegrees(passedDays float64) float64 {
	traveledDegress := geometry.GetTraveledDistance(p.DegreesByDay, passedDays)
	return math.Mod(p.InitialDegrees+traveledDegress, geometry.COMPLETE_SPIN_DEGREES)
}

// getPositionRadians returns the radians traveled by the planet starting from zero everytime a full spin is passed
func (p Planet) getPositionRadians(passedDays float64) float64 {
	return geometry.ToRadians(p.getPositionDegrees(passedDays))
}

func (p Planet) getCoordinate(passedDays float64) geometry.Coordinate {
	traveledRadians := p.getPositionRadians(passedDays)
	x := math.Cos(traveledRadians) * float64(p.DistanceToSun)
	y := math.Sin(traveledRadians) * float64(p.DistanceToSun)

	return geometry.NewCoordinate(x, y)
}

// rectEquationContains returns true if the provided rect equation is close enougth to the provided coordinate to be considered aligned
func rectEquationContains(re geometry.RectEquation, c geometry.Coordinate) bool {
	return re.DistanceToCoordinate(c) <= ALIGNMENT_THRESHOLD
}

// getClosestPlanet returns the closest planet to the sun from a given Planet slice
func getClosestPlanet(planets ...Planet) Planet {
	closestPlanet := planets[0]

	for _, planet := range planets {
		if planet.DistanceToSun < closestPlanet.DistanceToSun {
			closestPlanet = planet
		}
	}

	return closestPlanet
}

// getFarthestPlanet returns the farthest planet to the sun from a given Planet slice
func getFarthestPlanet(planets ...Planet) Planet {
	farthestPlanet := planets[0]

	for _, planet := range planets {
		if planet.DistanceToSun >= farthestPlanet.DistanceToSun {
			farthestPlanet = planet
		}
	}

	return farthestPlanet
}

func getAlignment(passedDays float64, planets ...Planet) alignment {
	closestPlanet := getClosestPlanet(planets...)
	farthestPlanet := getFarthestPlanet(planets...)

	closestCoordinate := closestPlanet.getCoordinate(passedDays)
	farthestCoordinate := farthestPlanet.getCoordinate(passedDays)

	equation := geometry.NewRectEquation(closestCoordinate, farthestCoordinate)

	for _, p := range planets {
		c := p.getCoordinate(passedDays)
		if !rectEquationContains(equation, c) {
			return alignmentNone
		}
	}

	if rectEquationContains(equation, sunCoordinate) {
		return alignmentWithSun
	}

	return alignmentWithoutSun
}

// isSunBetween calculates if a point is inside a triangle
// this is used algorithm https://www.youtube.com/watch?v=WaYS1gEXEFE&t=542s&ab_channel=huse360
func isSunBetween(p1, p2, p3 Planet, passedDays float64) bool {
	a := p1.getCoordinate(passedDays)
	b := p2.getCoordinate(passedDays)
	c := p3.getCoordinate(passedDays)

	return geometry.IsOriginInTriangle(a, b, c)
}

func GetClimate(passedDays float64, planets [3]Planet) string {
	planetsSlice := planets[:]
	alignment := getAlignment(passedDays, planetsSlice...)

	if alignment == alignmentWithoutSun {
		return ClimateIdeal
	}

	if alignment == alignmentWithSun {
		return ClimateFollow
	}

	if isSunBetween(planets[0], planets[1], planets[2], passedDays) {
		return ClimateRain
	}

	return ClimateNoInfo
}
