package domain

import "errors"

var (
	ErrNoClimateRecordFound    = errors.New("error_no_climate_record_found")
	ErrNoClimateRecordJobFound = errors.New("error_no_climate_record_job_found")
	ErrCreateClimateRecordJob  = errors.New("error_create_climate_record_job")
	ErrCreateClimateRecord     = errors.New("error_create_climate_record")
	ErrFindClimateRecordJob    = errors.New("error_find_climate_record_job")
	ErrFindClimateRecord       = errors.New("error_find_climate_record")
	ErrInvalidDayParam         = errors.New("invalid_day_param")
)
