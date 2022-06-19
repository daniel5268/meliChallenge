package geometry_test

import (
	"math"
	"testing"

	"github.com/daniel5268/meliChallenge/domain/geometry"
	"github.com/stretchr/testify/assert"
)

func TestRectDistanceToCoordinate(t *testing.T) {
	tests := []struct {
		name         string
		equation     geometry.RectEquation
		coordinate   geometry.Coordinate
		wantDistance float64
	}{
		{
			name: "should return the distance when a vertical equation is provided",
			equation: geometry.NewRectEquation(
				geometry.NewCoordinate(5, 2),
				geometry.NewCoordinate(5, 3),
			),
			coordinate:   geometry.NewCoordinate(10, 15),
			wantDistance: 5,
		},
		{
			name: "should return the distance when a not vertical equation is provided",
			equation: geometry.NewRectEquation(
				geometry.NewCoordinate(0, 1),
				geometry.NewCoordinate(1, 2),
			),
			coordinate:   geometry.NewCoordinate(4, 1),
			wantDistance: math.Sqrt(8),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDistance := tt.equation.DistanceToCoordinate(tt.coordinate)

			assert.True(t, almostEqual(tt.wantDistance, gotDistance))
		})
	}
}
