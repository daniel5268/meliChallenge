package app

import (
	"errors"
	"net/http"

	"github.com/daniel5268/meliChallenge/src/domain"
	"github.com/labstack/echo/v4"
)

var errDefault = errors.New("internal_server_error")

// APIError structure that represents an API error
type APIError struct {
	Type       string `json:"tipo"`
	Message    string `json:"mensaje"`
	StatusCode int    `json:"codigo"`
}

// ErrorHandler manages errors on application level
func ErrorHandler(err error, c echo.Context) {
	aErr := mapToAPIError(err)
	_ = c.JSON(aErr.StatusCode, aErr)
}

func NewAPIError(t error, m string, sc int) APIError {
	return APIError{
		Type:       t.Error(),
		Message:    m,
		StatusCode: sc,
	}
}

func mapToAPIError(err error) APIError {
	switch err {
	case domain.ErrInvalidDayParam:
		return NewAPIError(err, "El dia es invalido", http.StatusBadRequest)
	case domain.ErrNoClimateRecordFound:
		return NewAPIError(err, "No se encontraron registros para este dia", http.StatusNotFound)
	case domain.ErrFindClimateRecord:
		return NewAPIError(err, "Hubo un error encontrando este registro, por favor contacte el administrador", http.StatusInternalServerError)
	default:
		return NewAPIError(errDefault, "Error en el servidor, por favor contacte el administrador", http.StatusInternalServerError)
	}
}
