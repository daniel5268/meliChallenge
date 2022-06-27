package app

import (
	"github.com/labstack/echo/v4"
)

func (a *App) setupRoutes(g *echo.Group) {
	crHandler := a.Dependencies.climateRecords.handler
	climateRecordsAPIGroup := g.Group("/clima")

	climateRecordsAPIGroup.GET("", crHandler.GetClimateRecord)
}
