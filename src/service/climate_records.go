package service

import (
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
)

type ClimateRecordsRepository interface {
	Create(...*meteorology.ClimateRecord) error
	FindByDay(int64) (meteorology.ClimateRecord, error)
	FindBetweenDays(int64, int64) ([]meteorology.ClimateRecord, error)
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

func (crs *ClimateRecordsService) GetClimateRecordsSummary(firstDay, lastDay int64) (meteorology.ClimateRecordSummary, error) {
	climateRecords, err := crs.climateRecordsRepository.FindBetweenDays(firstDay, lastDay)
	if err != nil {
		return meteorology.ClimateRecordSummary{}, err
	}

	var maxPerimeter float64
	var maxPerimeterDay int64
	var followDays uint64
	var idealDays uint64
	var rainDays uint64

	for _, cr := range climateRecords {
		if cr.Perimeter > maxPerimeter {
			maxPerimeter = cr.Perimeter
			maxPerimeterDay = cr.Day
		}
		switch cr.Climate {
		case meteorology.ClimateFollow:
			followDays++
		case meteorology.ClimateIdeal:
			idealDays++
		case meteorology.ClimateRain:
			rainDays++
		}
	}

	return meteorology.ClimateRecordSummary{
		FirstDay:   firstDay,
		LastDay:    lastDay,
		MaxRainDay: maxPerimeterDay,
		FollowDays: followDays,
		IdealDays:  idealDays,
		RainDays:   rainDays,
	}, nil
}
