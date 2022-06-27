package geometry_test

import (
	"testing"

	"github.com/daniel5268/meliChallenge/src/domain/geometry"
	"github.com/stretchr/testify/assert"
)

func TestCoordinateDiff(t *testing.T) {
	t.Run("returns the difference between two coordinates", func(t *testing.T) {
		c := geometry.NewCoordinate(4.5, 6.8)
		c2 := geometry.NewCoordinate(3, 8)
		want := geometry.NewCoordinate(1.5, -1.2)
		got := c.Diff(c2)

		assert.True(t, almostEqual(want.X, got.X))
		assert.True(t, almostEqual(want.Y, got.Y))
	})
}

func TestCoordinateDistance(t *testing.T) {
	t.Run("returns the distance to the provided coordinate", func(t *testing.T) {
		c := geometry.NewCoordinate(1, 2)
		c2 := geometry.NewCoordinate(4, 6)
		want := float64(5)
		got := c.Distance(c2)

		assert.Equal(t, want, got)
	})
}
