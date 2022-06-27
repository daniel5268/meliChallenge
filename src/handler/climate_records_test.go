package handler_test

import (
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

func buildJSONRequest(day string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/clima", nil)
	q := req.URL.Query()
	q.Add("dia", day)
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
			req := buildJSONRequest(tt.day)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			h := handler.NewClimateRecordsHandler(tt.service)
			gotErr := h.GetClimateRecord(c)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
