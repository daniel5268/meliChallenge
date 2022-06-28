package service

import (
	"errors"
	"sync"
	"time"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
)

type ClimateRecordJobsService struct {
	climateRecordsRepository    ClimateRecordsRepository
	climateRecordJobsRepository ClimateRecordJobsRepository
}

type ClimateRecordJobsRepository interface {
	Create(...*meteorology.ClimateRecordJob) error
	FindLast() (meteorology.ClimateRecordJob, error)
}

func NewClimateRecordJobsService(crR ClimateRecordsRepository, crjR ClimateRecordJobsRepository) *ClimateRecordJobsService {
	return &ClimateRecordJobsService{
		climateRecordsRepository:    crR,
		climateRecordJobsRepository: crjR,
	}
}

// CreateClimateRecordJob creates a new climate_record_job and the climate_records associated to the job
// if no climate_record_jobs are found it creates 10 years of climate_records
// if climate_record_jobs are found it creates a number of climate_records = the days that passed since the last created job
func (crs *ClimateRecordJobsService) CreateClimateRecordJob() error {
	var firstDay int64 = 0   // if no jobs have been executed, those will be limits to calculate climateRecords
	var lastDay int64 = 3650 // ten years

	lastClimateRecordJob, err := crs.climateRecordJobsRepository.FindLast()
	isFirstExecution := errors.Is(err, domain.ErrNoClimateRecordJobFound)

	if err != nil && !isFirstExecution {
		return err
	}

	if !isFirstExecution { // if a job was executed before, the new limits are calculated acording to the time that has passed
		firstDay = lastClimateRecordJob.LastDay + 1
		currentDate := time.Now()
		passedDaysFromLastJob := currentDate.Sub(lastClimateRecordJob.CreatedAt).Hours() / 24
		lastDay = firstDay + int64(passedDaysFromLastJob)
	}

	if firstDay == lastDay { // if the job was executed today there is no need to execute again
		return nil
	}

	climateRecordsLength := lastDay - firstDay
	climateRecords := make([]*meteorology.ClimateRecord, climateRecordsLength)
	planets := meteorology.GetPlanets()

	var wg sync.WaitGroup
	for currentDay := firstDay; currentDay < lastDay; currentDay++ {
		wg.Add(1)
		go func(currentDay int64) {
			defer wg.Done()
			climateRecords[currentDay-firstDay] = meteorology.GetDayClimateRecord(currentDay, planets)
		}(currentDay)
	}
	wg.Wait()

	if err = crs.climateRecordsRepository.Create(climateRecords...); err != nil {
		return err
	}

	crj := &meteorology.ClimateRecordJob{
		FirstDay: firstDay,
		LastDay:  lastDay - 1,
	}

	if err = crs.climateRecordJobsRepository.Create(crj); err != nil {
		return err
	}

	return nil
}
