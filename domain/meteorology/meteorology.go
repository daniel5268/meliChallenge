package meteorology

import (
	"math"
)

const (
	completeSpin = float64(360)
	halfSpin     = completeSpin / 2
	allowedError = float64(0.01) // closeness allowed to the main line for a coordinate to be considered aligned
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

type Planet struct {
	DegreesByDay   float64
	InitialDegrees float64
	DistanceToSun  uint64
}

func NewPlanet(degreesByDay float64, passedDays float64, initialDegrees float64, distanceToSun uint64) *Planet {
	return &Planet{
		DegreesByDay:   degreesByDay,
		InitialDegrees: initialDegrees,
		DistanceToSun:  distanceToSun,
	}
}

// getPositionDegrees returns the degrees traveled by the planet starting from zero everytime a full spin is passed
func (p Planet) getPositionDegrees(passedDays float64) float64 {
	return math.Mod(p.InitialDegrees+p.DegreesByDay*passedDays, completeSpin)
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / halfSpin
}

// getPositionRadians returns the radians traveled by the planet starting from zero everytime a full spin is passed
func (p Planet) getPositionRadians(passedDays float64) float64 {
	return toRadians(p.getPositionDegrees(passedDays))
}

type coordinate struct {
	x float64
	y float64
}

var sunCoordinate = coordinate{
	x: 0,
	y: 0,
}

func (p Planet) getCoordinate(passedDays float64) coordinate {
	traveledRadians := p.getPositionRadians(passedDays)
	x := math.Cos(traveledRadians) * float64(p.DistanceToSun)
	y := math.Sin(traveledRadians) * float64(p.DistanceToSun)

	return coordinate{
		x: x,
		y: y,
	}
}

type rectEquation struct {
	m          float64
	b          float64
	isVertical bool
	xVertical  float64
}

// distanceToCoordinate returns the minimun distance fron a given coordinate to the rect described by the ecuation
func (re rectEquation) distanceToCoordinate(c coordinate) float64 {
	if re.isVertical {
		return math.Abs(re.xVertical - c.x)
	}

	return math.Abs(re.m*c.x-c.y+re.b) / math.Sqrt(re.m*re.m+1)
}

// returns true if the line described by the ecuation close enougth to coordinate to be considered aligned
func (re rectEquation) contains(c coordinate) bool {
	return re.distanceToCoordinate(c) <= allowedError
}

func getRectEcuation(c1, c2 coordinate) rectEquation {
	denominator := c1.x - c2.x
	if denominator == 0 {
		return rectEquation{
			isVertical: true,
			xVertical:  c1.x,
		}
	}
	m := (c1.y - c2.y) / denominator
	b := c1.y - m*c1.x

	return rectEquation{
		m: m,
		b: b,
	}
}

func getClosestPlanet(planets ...Planet) Planet {
	closestPlanet := planets[0]

	for _, planet := range planets {
		if planet.DistanceToSun < closestPlanet.DistanceToSun {
			closestPlanet = planet
		}
	}

	return closestPlanet
}

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

	ecuation := getRectEcuation(closestCoordinate, farthestCoordinate)

	for _, p := range planets {
		c := p.getCoordinate(passedDays)
		if !ecuation.contains(c) {
			return alignmentNone
		}
	}

	if ecuation.contains(sunCoordinate) {
		return alignmentWithSun
	}

	return alignmentWithoutSun
}

func diff(c1, c2 coordinate) coordinate {
	return coordinate{
		x: c1.x - c2.x,
		y: c1.y - c2.y,
	}
}

// isSunBetween calculates if a point is inside a triangle
// this is used algorithm https://www.youtube.com/watch?v=WaYS1gEXEFE&t=542s&ab_channel=huse360
func isSunBetween(p1, p2, p3 Planet, passedDays float64) bool {
	a := p1.getCoordinate(passedDays)
	b := p2.getCoordinate(passedDays)
	c := p3.getCoordinate(passedDays)

	d := diff(b, a)
	e := diff(c, a)

	w1 := (e.x*a.y - e.y*a.x) / (d.x*e.y - d.y*e.x)
	w2 := -(a.y + w1*d.y) / e.y

	return w1 >= 0 && w2 >= 0 && w1+w2 <= 1
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
