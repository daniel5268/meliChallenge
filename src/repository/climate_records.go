package repository

import (
	"errors"
	"log"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"gorm.io/gorm"
)

type ClimateRecordsRepository struct {
	db *gorm.DB
}

func NewClimateRecordsRepository(db *gorm.DB) *ClimateRecordsRepository {
	return &ClimateRecordsRepository{
		db: db,
	}
}

func (r *ClimateRecordsRepository) Create(cr ...*meteorology.ClimateRecord) error {
	err := r.db.Create(cr).Error

	if err != nil {
		log.Print(err)
		return domain.ErrCreateClimateRecord
	}

	return nil
}

func (r *ClimateRecordsRepository) FindByDay(day int64) (meteorology.ClimateRecord, error) {
	cr := meteorology.ClimateRecord{}
	result := r.db.Where("day = ?", day).First(&cr)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return cr, domain.ErrNoClimateRecordFound
	}

	if result.Error != nil {
		log.Print(result.Error)
		return cr, domain.ErrFindClimateRecord
	}

	return cr, nil
}

func (r *ClimateRecordsRepository) FindBetweenDays(firstDay int64, lastDay int64) ([]meteorology.ClimateRecord, error) {
	climateRecords := []meteorology.ClimateRecord{}
	result := r.db.Where("day >= ? AND day <= ?", firstDay, lastDay).Find(&climateRecords)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return climateRecords, domain.ErrNoClimateRecordFound
	}

	if result.Error != nil {
		log.Print(result.Error)
		return climateRecords, domain.ErrFindClimateRecord
	}

	return climateRecords, nil
}
