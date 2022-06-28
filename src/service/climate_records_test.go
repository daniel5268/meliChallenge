package service_test

import (
	"errors"
	"testing"

	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/daniel5268/meliChallenge/src/service"
	"github.com/daniel5268/meliChallenge/src/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestClimateRecordsServiceGetClimateRecordsSummary(t *testing.T) {
	errTest := errors.New("test error")
	tests := []struct {
		name        string
		firstDay    int64
		lastDay     int64
		repository  service.ClimateRecordsRepository
		wantSummary meteorology.ClimateRecordSummary
		wantErr     error
	}{
		{
			name: "returns an error if the repository fails",
			repository: func() service.ClimateRecordsRepository {
				repositoryMock := &mocks.ClimateRecordsRepository{}
				repositoryMock.On("FindBetweenDays", int64(0), int64(0)).Return([]meteorology.ClimateRecord{}, errTest)
				return repositoryMock
			}(),
			wantErr: errTest,
		},
		{
			name:     "returns the ClimateRecordSummary",
			firstDay: 5,
			lastDay:  10,
			repository: func() service.ClimateRecordsRepository {
				repositoryMock := &mocks.ClimateRecordsRepository{}
				repositoryMock.On("FindBetweenDays", int64(5), int64(10)).Return([]meteorology.ClimateRecord{
					{Day: 5, Climate: meteorology.ClimateFollow},
					{Day: 6, Climate: meteorology.ClimateRain, Perimeter: 5.5},
					{Day: 7, Climate: meteorology.ClimateRain, Perimeter: 5.6},
					{Day: 8, Climate: meteorology.ClimateIdeal},
					{Day: 9, Climate: meteorology.ClimateRain, Perimeter: 5.4},
					{Day: 10, Climate: meteorology.ClimateIdeal},
				}, nil)
				return repositoryMock
			}(),
			wantSummary: meteorology.ClimateRecordSummary{
				FirstDay:   5,
				LastDay:    10,
				MaxRainDay: 7,
				FollowDays: 1,
				RainDays:   3,
				IdealDays:  2,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewClimateRecordsService(tt.repository)
			gotSummary, gotErr := s.GetClimateRecordsSummary(tt.firstDay, tt.lastDay)
			assert.Equal(t, tt.wantSummary, gotSummary)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}

}
