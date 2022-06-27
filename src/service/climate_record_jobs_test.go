package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/daniel5268/meliChallenge/src/service"
	"github.com/daniel5268/meliChallenge/src/service/mocks"
	"github.com/stretchr/testify/assert"
)

//go:generate mockery --name ClimateRecordsRepository
//go:generate mockery --name ClimateRecordJobsRepository

func TestClimateRecordJobsServiceCreateClimateRecordJob(t *testing.T) {
	errTest := errors.New("something happened :/")
	crJob := &meteorology.ClimateRecordJob{
		FirstDay: 89,
		LastDay:  91,
	}
	wantClimateRecords := []*meteorology.ClimateRecord{
		{
			Day:     89,
			Climate: meteorology.ClimateRain,
		},
		{
			Day:     90,
			Climate: meteorology.ClimateFollow,
		},
		{
			Day:     91,
			Climate: meteorology.ClimateRain,
		},
	}

	tests := []struct {
		name      string
		crR       service.ClimateRecordsRepository
		jR        service.ClimateRecordJobsRepository
		wantError error
	}{
		{
			name: "creates the ClimateRecordJob and ClimateRecords when jobs are found",
			crR: func() service.ClimateRecordsRepository {
				repositoryMock := &mocks.ClimateRecordsRepository{}
				repositoryMock.On("Create", wantClimateRecords[0], wantClimateRecords[1], wantClimateRecords[2]).Return(nil)

				return repositoryMock
			}(),
			jR: func() service.ClimateRecordJobsRepository {
				repositoryMock := &mocks.ClimateRecordJobsRepository{}
				repositoryMock.On("FindLast").Return(meteorology.ClimateRecordJob{
					LastDay:   88,
					CreatedAt: time.Now().AddDate(0, 0, -3),
				}, nil)
				repositoryMock.On("Create", crJob).Return(nil)
				return repositoryMock
			}(),
		},
		{
			name: "it doesn't calculate ClimateRecords when the job was executed for second time today",
			crR:  &mocks.ClimateRecordsRepository{},
			jR: func() service.ClimateRecordJobsRepository {
				repositoryMock := &mocks.ClimateRecordJobsRepository{}
				repositoryMock.On("FindLast").Return(
					meteorology.ClimateRecordJob{
						LastDay:   88,
						CreatedAt: time.Now(),
					},
					nil,
				)
				return repositoryMock
			}(),
		},
		{
			name: "returns an error when the ClimateRecordJobsRepository fails finding the last ClimateRecordJob",
			crR:  &mocks.ClimateRecordsRepository{},
			jR: func() service.ClimateRecordJobsRepository {
				repositoryMock := &mocks.ClimateRecordJobsRepository{}
				repositoryMock.On("FindLast").Return(
					meteorology.ClimateRecordJob{},
					errTest,
				)
				return repositoryMock
			}(),
			wantError: errTest,
		},
		{
			name: "returns an error when the ClimateRecordsRepository fails creating the ClimateRecord",
			crR: func() service.ClimateRecordsRepository {
				repositoryMock := &mocks.ClimateRecordsRepository{}
				repositoryMock.On("Create", wantClimateRecords[0], wantClimateRecords[1], wantClimateRecords[2]).Return(
					errTest,
				)

				return repositoryMock
			}(),
			jR: func() service.ClimateRecordJobsRepository {
				repositoryMock := &mocks.ClimateRecordJobsRepository{}
				repositoryMock.On("FindLast").Return(meteorology.ClimateRecordJob{
					LastDay:   88,
					CreatedAt: time.Now().AddDate(0, 0, -3),
				}, nil)
				repositoryMock.On("Create", crJob).Return(nil)
				return repositoryMock
			}(),
			wantError: errTest,
		},
		{
			name: "returns an error when the ClimateRecordJobsRepository fails creating the ClimateRecordJob",
			crR: func() service.ClimateRecordsRepository {
				repositoryMock := &mocks.ClimateRecordsRepository{}
				repositoryMock.On("Create", wantClimateRecords[0], wantClimateRecords[1], wantClimateRecords[2]).Return(nil)

				return repositoryMock
			}(),
			jR: func() service.ClimateRecordJobsRepository {
				repositoryMock := &mocks.ClimateRecordJobsRepository{}
				repositoryMock.On("FindLast").Return(meteorology.ClimateRecordJob{
					LastDay:   88,
					CreatedAt: time.Now().AddDate(0, 0, -3),
				}, nil)
				repositoryMock.On("Create", crJob).Return(errTest)
				return repositoryMock
			}(),
			wantError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewClimateRecordJobsService(tt.crR, tt.jR)
			gotErr := s.CreateClimateRecordJob()
			assert.Equal(t, tt.wantError, gotErr)
		})
	}
}
