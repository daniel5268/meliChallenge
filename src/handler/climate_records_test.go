package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/daniel5268/meliChallenge/src/domain/meteorology"
	"github.com/daniel5268/meliChallenge/src/handler"
	"github.com/daniel5268/meliChallenge/src/handler/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func buildJSONRequest(path string, queryKeys []string, queryValues []string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	q := req.URL.Query()

	for i, key := range queryKeys {
		q.Add(key, queryValues[i])
	}

	req.URL.RawQuery = q.Encode()

	return req
}

//go:generate mockery --name ClimateRecordsService
func TestClimateRecordsHandlerGetClimateRecord(t *testing.T) {
	wantClimateRecord := meteorology.ClimateRecord{
		Day:     5,
		Climate: meteorology.ClimateIdeal,
	}
	tests := []struct {
		name    string
		day     string
		service handler.ClimateRecordsService
		wantErr error
	}{
		{
			name: "returns the climate record",
			day:  "5",
			service: func() handler.ClimateRecordsService {
				crsMock := &mocks.ClimateRecordsService{}
				crsMock.On("GetClimateRecord", int64(5)).Return(wantClimateRecord, nil)

				return crsMock
			}(),
		},
		{
			name: "returns the first climate record if no day is provided",
			service: func() handler.ClimateRecordsService {
				crsMock := &mocks.ClimateRecordsService{}
				crsMock.On("GetClimateRecord", int64(0)).Return(wantClimateRecord, nil)

				return crsMock
			}(),
		},
		{
			name:    "returns an error when the day query param is not integer",
			day:     "notInteger",
			service: &mocks.ClimateRecordsService{},
			wantErr: domain.ErrInvalidDayParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := buildJSONRequest("/clima", []string{"dia"}, []string{tt.day})
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			h := handler.NewClimateRecordsHandler(tt.service)
			gotErr := h.GetClimateRecord(c)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func TestClimateRecordsHandlerGetClimateRecordsSummary(t *testing.T) {
	wantClimateRecordSummary := meteorology.ClimateRecordSummary{
		FirstDay:   89,
		LastDay:    91,
		RainDays:   2,
		FollowDays: 1,
		IdealDays:  0,
		MaxRainDay: 91,
	}
	errTest := errors.New("test error")
	tests := []struct {
		name     string
		firstDay string
		lastDay  string
		service  handler.ClimateRecordsService
		wantErr  error
	}{
		{
			name:     "returns the climate record summary",
			firstDay: "89",
			lastDay:  "91",
			service: func() handler.ClimateRecordsService {
				crsMock := &mocks.ClimateRecordsService{}
				crsMock.On("GetClimateRecordsSummary", int64(89), int64(91)).Return(wantClimateRecordSummary, nil)

				return crsMock
			}(),
		},
		{
			name:     "returns an invalid_day_param when the first_day query param is not integer",
			firstDay: "notInteger",
			service:  &mocks.ClimateRecordsService{},
			wantErr:  domain.ErrInvalidDayParam,
		},
		{
			name:    "returns an invalid_day_param when the last_day query param is not integer",
			lastDay: "notInteger",
			service: &mocks.ClimateRecordsService{},
			wantErr: domain.ErrInvalidDayParam,
		},
		{
			name:     "returns an error when the service fails",
			firstDay: "89",
			lastDay:  "91",
			service: func() handler.ClimateRecordsService {
				crsMock := &mocks.ClimateRecordsService{}
				crsMock.On("GetClimateRecordsSummary", int64(89), int64(91)).Return(meteorology.ClimateRecordSummary{}, errTest)

				return crsMock
			}(),
			wantErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := buildJSONRequest("/clima/resumen", []string{"primer_dia", "ultimo_dia"}, []string{tt.firstDay, tt.lastDay})
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			h := handler.NewClimateRecordsHandler(tt.service)
			gotErr := h.GetClimateRecordsSummary(c)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
