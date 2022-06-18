package meteorology_test

import (
	"testing"

	"github.com/daniel5268/meliChallenge/domain/meteorology"
	"github.com/stretchr/testify/assert"
)

func TestGetClimate(t *testing.T) {
	planets := [3]meteorology.Planet{
		{
			InitialDegrees: 90,
			DistanceToSun:  500,
			DegreesByDay:   -1,
		},
		{
			InitialDegrees: 90,
			DistanceToSun:  2000,
			DegreesByDay:   -3,
		},
		{
			InitialDegrees: 90,
			DistanceToSun:  1000,
			DegreesByDay:   5,
		},
	}

	tests := []struct {
		name        string
		passedDays  float64
		wantClimate string
	}{
		{
			name:        "returns Follow when planets are aligned with the sun",
			passedDays:  90,
			wantClimate: meteorology.ClimateFollow,
		},
		{
			name:        "returns Ideal when planets are aligned without the sun",
			passedDays:  18.46431,
			wantClimate: meteorology.ClimateIdeal,
		},
		{
			name:        "returns Rain when the sun is between the triangle formed by planets",
			passedDays:  28,
			wantClimate: meteorology.ClimateRain,
		},
		{
			name:        "returns no information when none of the conditions are fulfilled",
			passedDays:  20,
			wantClimate: meteorology.ClimateNoInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClimate := meteorology.GetClimate(tt.passedDays, planets)
			assert.Equal(t, tt.wantClimate, gotClimate)
		})
	}
}
