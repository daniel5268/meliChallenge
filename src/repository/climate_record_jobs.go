package repository

import (
	"errors"
	"log"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"gorm.io/gorm"
)

type ClimateRecordJobsRepository struct {
	db *gorm.DB
}

func NewClimateRecordJobsRepository(db *gorm.DB) *ClimateRecordJobsRepository {
	return &ClimateRecordJobsRepository{
		db: db,
	}
}

func (r *ClimateRecordJobsRepository) Create(crJob ...*meteorology.ClimateRecordJob) error {
	err := r.db.Create(crJob).Error

	if err != nil {
		log.Print(err)
		return domain.ErrCreateClimateRecordJob
	}

	return nil
}

func (r *ClimateRecordJobsRepository) FindLast() (meteorology.ClimateRecordJob, error) {
	crJob := meteorology.ClimateRecordJob{}
	result := r.db.Last(&crJob)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return crJob, domain.ErrNoClimateRecordJobFound
	}

	if result.Error != nil {
		log.Print(result.Error)
		return crJob, domain.ErrFindClimateRecordJob
	}

	return crJob, nil
}
