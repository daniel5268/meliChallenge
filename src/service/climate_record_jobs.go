package service

import (
	"errors"
	"fmt"
	"os"
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
	Create(cr ...*meteorology.ClimateRecordJob) error
	FindLast() (meteorology.ClimateRecordJob, error)
}

func NewClimateRecordJobsService(crR ClimateRecordsRepository, crjR ClimateRecordJobsRepository) *ClimateRecordJobsService {
	return &ClimateRecordJobsService{
		climateRecordsRepository:    crR,
		climateRecordJobsRepository: crjR,
	}
}

func saveResponses(climateRecords []*meteorology.ClimateRecord, perimeters []float64) {
	followCount := 0
	idealCount := 0
	rainCount := 0

	for _, cr := range climateRecords {
		switch cr.Climate {
		case meteorology.ClimateFollow:
			followCount += 1
		case meteorology.ClimateIdeal:
			idealCount += 1
		case meteorology.ClimateRain:
			rainCount += 1
		}
	}

	var maxPerimeter float64
	var maxPerimeterIndex int

	for i, p := range perimeters {
		if p > maxPerimeter {
			maxPerimeter = p
			maxPerimeterIndex = i
		}
	}
	resultsFile := "./results"
	f, err := os.Create(resultsFile)
	if err != nil {
		panic("error creating results file")
	}
	defer f.Close()
	results := fmt.Sprintf("lluvia:%d maxima:%d\nsequia:%d\noptimo:%d", rainCount, maxPerimeterIndex, followCount, idealCount)
	f.WriteString(results)
}

// CreateClimateRecords calculates and creates a list of climateRecords given an array of three planets
// if no climate records are found it creates 10 years of climateRecords
// if climate records are found it stores the next climateRecord
func (crs *ClimateRecordJobsService) CreateClimateRecordJob() error {

	var firstDay int64 = 0 // if no jobs have been executed, those will be limits to calculate climateRecords
	var lastDay int64 = 3650

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

	newClimateRecordsLength := lastDay - firstDay
	climateRecords := make([]*meteorology.ClimateRecord, newClimateRecordsLength)
	perimeters := make([]float64, newClimateRecordsLength)

	planets := meteorology.GetPlanets()
	var wg sync.WaitGroup
	for currentDay := firstDay; currentDay < lastDay; currentDay++ {
		wg.Add(1)
		go func(currentDay int64) {
			defer wg.Done()
			dayClimate, exactMoment := meteorology.GetDayClimate(currentDay, planets)
			climateRecords[currentDay-firstDay] = &meteorology.ClimateRecord{
				Day:     currentDay,
				Climate: dayClimate,
			}
			if isFirstExecution && dayClimate == meteorology.ClimateRain { // it only calculates perimeter if it's the first execution
				perimeters[currentDay-firstDay] = meteorology.GetPerimeter(exactMoment, planets)
			}
		}(currentDay)
	}
	wg.Wait()

	if isFirstExecution {
		saveResponses(climateRecords, perimeters)
	}

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
