package handler

import (
	"net/http"
	"strconv"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/labstack/echo/v4"
)

const (
	dayQueryParam      = "dia"
	fisrtDayQueryParam = "primer_dia"
	lastDayQueryParam  = "ultimo_dia"
)

type ClimateRecordsService interface {
	GetClimateRecord(int64) (meteorology.ClimateRecord, error)
	GetClimateRecordsSummary(int64, int64) (meteorology.ClimateRecordSummary, error)
}

type ClimateRecordsHandler struct {
	climateRecordsService ClimateRecordsService
}

func NewClimateRecordsHandler(crs ClimateRecordsService) *ClimateRecordsHandler {
	return &ClimateRecordsHandler{
		climateRecordsService: crs,
	}
}

// int64QueryParam tries to get the given query param value from the echo Context provided, as an integer
func int64QueryParam(c echo.Context, name string) (int64, error) {
	strValue := c.QueryParam(name)
	if strValue == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		return 0, domain.ErrInvalidDayParam
	}

	return int64(intValue), nil
}

func (crh *ClimateRecordsHandler) GetClimateRecord(c echo.Context) error {
	day, err := int64QueryParam(c, dayQueryParam)

	if err != nil {
		return err
	}

	cr, err := crh.climateRecordsService.GetClimateRecord(day)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cr)
}

func (crh *ClimateRecordsHandler) GetClimateRecordsSummary(c echo.Context) error {
	firstDay, err := int64QueryParam(c, fisrtDayQueryParam)
	if err != nil {
		return err
	}

	lastDay, err := int64QueryParam(c, lastDayQueryParam)
	if err != nil {
		return err
	}

	crs, err := crh.climateRecordsService.GetClimateRecordsSummary(firstDay, lastDay)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, crs)
}
