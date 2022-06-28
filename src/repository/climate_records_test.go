package repository_test

import (
	"testing"

	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/daniel5268/meliChallenge/src/infrastructure"
	"github.com/daniel5268/meliChallenge/src/repository"
	"github.com/stretchr/testify/assert"
)

func TestClimateRecordsRepositoryCreate(t *testing.T) {
	db := infrastructure.NewGormPostgresClient()
	repository := repository.NewClimateRecordsRepository(db)
	tests := []struct {
		name    string
		cr      *meteorology.ClimateRecord
		wantErr bool
	}{
		{
			name: "creates a ClimateRecord",
			cr: &meteorology.ClimateRecord{
				Day:       2,
				Climate:   meteorology.ClimateIdeal,
				Perimeter: 6.4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := repository.Create(tt.cr)

			assert.Equal(t, tt.wantErr, gotErr != nil)
		})
	}
}
