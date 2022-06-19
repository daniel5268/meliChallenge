package geometry_test

import (
	"math"
	"testing"

	"github.com/daniel5268/meliChallenge/domain/geometry"
	"github.com/stretchr/testify/assert"
)

const FLOAT_64_EQUALITY_THRESHOLD = 0.00000000001

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= FLOAT_64_EQUALITY_THRESHOLD
}
func TestToRadians(t *testing.T) {
	t.Run("should convert degrees to radians", func(t *testing.T) {
		assert.Equal(t, math.Pi, geometry.ToRadians(180))
	})
}

func TestGetTraveledDistance(t *testing.T) {
	t.Run("should return the traveled distance", func(t *testing.T) {
		assert.Equal(t, float64(60), geometry.GetTraveledDistance(5, 12))
	})
}

func TestIsOriginInTriangle(t *testing.T) {
	test := []struct {
		name string
		a    geometry.Coordinate
		b    geometry.Coordinate
		c    geometry.Coordinate
		want bool
	}{
		{
			name: "should return true when the triangle contains the origin",
			a:    geometry.NewCoordinate(-1, 1),
			b:    geometry.NewCoordinate(1, 0),
			c:    geometry.NewCoordinate(0, -1),
			want: true,
		},
		{
			name: "should return false when the triangle doesn't contain the origin",
			a:    geometry.NewCoordinate(-1, 1),
			b:    geometry.NewCoordinate(1, 1),
			c:    geometry.NewCoordinate(0.5, 0.5),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, geometry.IsOriginInTriangle(tt.a, tt.b, tt.c))
		})
	}
}

func TestGetTrianglePerimeter(t *testing.T) {
	t.Run("should return the perimeter", func(t *testing.T) {
		want := float64(12)
		got := geometry.GetTrianglePerimeter(
			geometry.NewCoordinate(1, 1),
			geometry.NewCoordinate(5, 1),
			geometry.NewCoordinate(5, 4),
		)

		assert.Equal(t, want, got)
	})
}
