package service

import (
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
)

type ClimateRecordsRepository interface {
	Create(cr ...*meteorology.ClimateRecord) error
	FindByDay(day int64) (meteorology.ClimateRecord, error)
}

type ClimateRecordsService struct {
	climateRecordsRepository ClimateRecordsRepository
}

func NewClimateRecordsService(crRepository ClimateRecordsRepository) *ClimateRecordsService {
	return &ClimateRecordsService{
		climateRecordsRepository: crRepository,
	}
}

func (crs *ClimateRecordsService) GetClimateRecord(day int64) (meteorology.ClimateRecord, error) {
	return crs.climateRecordsRepository.FindByDay(day)
}
