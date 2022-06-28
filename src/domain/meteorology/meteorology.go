package meteorology

import (
	"errors"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/daniel5268/meliChallenge/src/domain/geometry"
)

var ErrInvalidEnv = errors.New("invalid_environment")

const (
	// alignmentThreshold represents the closeness allowed to the main line for a coordinate to be considered aligned
	alignmentThreshold = float64(0.01)
)

const (
	Ferengi   = "ferengi"
	Betasoide = "betasoide"
	Vulcano   = "vulcano"
)

type alignment int16

const (
	alignmentWithoutSun alignment = iota
	alignmentWithSun
	alignmentNone
)

type Climate string

const (
	ClimateFollow Climate = "sequia"
	ClimateRain   Climate = "lluvia"
	ClimateIdeal  Climate = "condiciones_optimas"
	ClimateNoInfo Climate = "sin_informacion"
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

type ClimateRecord struct {
	ID        uint64    `json:"-"`
	Day       int64     `json:"dia"`
	Perimeter float64   `json:"-"`
	Climate   Climate   `json:"clima"`
	CreatedAt time.Time `json:"-"`
}

type ClimateRecordSummary struct {
	FirstDay   int64  `json:"primer_dia"`
	LastDay    int64  `json:"ultimo_dia"`
	MaxRainDay int64  `json:"dia_lluvia_maxima,omitempty"`
	FollowDays uint64 `json:"dias_sequia"`
	RainDays   uint64 `json:"dias_lluvia"`
	IdealDays  uint64 `json:"dias_condiciones_optimas"`
}

type ClimateRecordJob struct {
	ID        uint64
	FirstDay  int64
	LastDay   int64
	CreatedAt time.Time
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

func GetPerimeter(passedDays float64, planets [3]Planet) float64 {
	c0 := planets[0].getCoordinate(passedDays)
	c1 := planets[1].getCoordinate(passedDays)
	c2 := planets[2].getCoordinate(passedDays)

	return geometry.GetTrianglePerimeter(c0, c1, c2)
}

// rectEquationContains returns true if the provided rect equation is close enougth to the provided coordinate to be considered aligned
func rectEquationContains(re geometry.RectEquation, c geometry.Coordinate) bool {
	return re.DistanceToCoordinate(c) <= alignmentThreshold
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

// GetClimate creturns the Climate given a moment (passed days) and an array of three planets
func GetMomentClimate(moment float64, planets [3]Planet) Climate {
	planetsSlice := planets[:]
	alignment := getAlignment(moment, planetsSlice...)

	if alignment == alignmentWithoutSun {
		return ClimateIdeal
	}

	if alignment == alignmentWithSun {
		return ClimateFollow
	}

	if isSunBetween(planets[0], planets[1], planets[2], moment) {
		return ClimateRain
	}

	return ClimateNoInfo
}

// GetCurrentDayClimate returns an special Climate (ClimateIdeal, ClimateRain, ClimateFollow)
// if one of this is reached in the provided day, if not, it returns ClimateNoInfo
func GetDayClimate(day int64, planets [3]Planet) (Climate, float64) {
	limit := float64(day + 1)
	momentDelta := float64(0.00002)
	for moment := float64(day); moment < limit; moment += momentDelta {
		currentClimate := GetMomentClimate(moment, planets)

		if currentClimate != ClimateNoInfo {
			return currentClimate, moment
		}
	}

	return ClimateNoInfo, 0
}

func GetDayClimateRecord(day int64, planets [3]Planet) *ClimateRecord {
	dayClimate, exactMoment := GetDayClimate(day, planets)
	climateRecord := &ClimateRecord{
		Day:     day,
		Climate: dayClimate,
	}
	if dayClimate == ClimateRain { // it only calculates perimeter if climate is rain
		climateRecord.Perimeter = GetPerimeter(exactMoment, planets)
	}

	return climateRecord
}

func getEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrInvalidEnv
	}
	return v, nil
}

func getSafeEnvFloat64(key string) float64 {
	s, err := getEnvStr(key)
	if err != nil {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func GetPlanets() [3]Planet {
	ferengiDegreesByDay := getSafeEnvFloat64("FERENGI_DEGREES_BY_DAY")
	betasoideDegreesByDay := getSafeEnvFloat64("BETASOIDE_DEGREES_BY_DAY")
	vulcanoDegreesByDay := getSafeEnvFloat64("VULCANO_DEGREES_BY_DAY")

	ferengiInitialDegrees := getSafeEnvFloat64("FERENGI_INITIAL_DEGREES")
	betasoideInitialDegrees := getSafeEnvFloat64("BETASOIDE_INITIAL_DEGREES")
	vulcanoInitialDegrees := getSafeEnvFloat64("VULCANO_INITIAL_DEGREES")

	ferengiDistanceToSun := getSafeEnvFloat64("FERENGI_DISTANCE_TO_SUN")
	betasoideDistanceToSun := getSafeEnvFloat64("BETASOIDE_DISTANCE_TO_SUN")
	vulcanoDistanceToSun := getSafeEnvFloat64("VULCANO_DISTANCE_TO_SUN")

	ferengi := NewPlanet(ferengiDegreesByDay, ferengiInitialDegrees, uint64(ferengiDistanceToSun))
	betasoide := NewPlanet(betasoideDegreesByDay, betasoideInitialDegrees, uint64(betasoideDistanceToSun))
	vulcano := NewPlanet(vulcanoDegreesByDay, vulcanoInitialDegrees, uint64(vulcanoDistanceToSun))

	return [3]Planet{ferengi, betasoide, vulcano}
}
