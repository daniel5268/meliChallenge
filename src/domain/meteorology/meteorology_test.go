package meteorology_test

import (
	"testing"

	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/stretchr/testify/assert"
)

func TestGetMomentClimate(t *testing.T) {
	planets := [3]meteorology.Planet{
		meteorology.NewPlanet(-1, 90, 500),
		meteorology.NewPlanet(-3, 90, 2000),
		meteorology.NewPlanet(5, 90, 1000),
	}

	tests := []struct {
		name        string
		passedDays  float64
		wantClimate meteorology.Climate
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
			gotClimate := meteorology.GetMomentClimate(tt.passedDays, planets)
			assert.Equal(t, tt.wantClimate, gotClimate)
		})
	}
}

func TestGetGetDayClimate(t *testing.T) {
	planets := [3]meteorology.Planet{
		meteorology.NewPlanet(-1, 90, 500),
		meteorology.NewPlanet(-3, 90, 2000),
		meteorology.NewPlanet(5, 90, 1000),
	}

	tests := []struct {
		name        string
		passedDays  int64
		wantClimate meteorology.Climate
	}{
		{
			name:        "returns Follow when planets are aligned with the sun in the provided day",
			passedDays:  90,
			wantClimate: meteorology.ClimateFollow,
		},
		{
			name:        "returns Ideal when planets are aligned without the sun in the provided day",
			passedDays:  18,
			wantClimate: meteorology.ClimateIdeal,
		},
		{
			name:        "returns Rain when the sun is between the triangle formed by planets in the provided day",
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
			gotClimate, _ := meteorology.GetDayClimate(tt.passedDays, planets)
			assert.Equal(t, tt.wantClimate, gotClimate)
		})
	}
}
