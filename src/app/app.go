package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/daniel5268/meliChallenge/src/handler"
	"github.com/daniel5268/meliChallenge/src/infrastructure"
	"github.com/daniel5268/meliChallenge/src/repository"
	"github.com/daniel5268/meliChallenge/src/service"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

const (
	ApplicationName = "challenge"
)

type climateRecordJobsDependencies struct {
	service *service.ClimateRecordJobsService
}

type climateRecordsDependencies struct {
	handler *handler.ClimateRecordsHandler
}

type dependencies struct {
	climateRecordJobs *climateRecordJobsDependencies
	climateRecords    *climateRecordsDependencies
}

type App struct {
	Server        *echo.Echo
	ServerReady   chan bool
	ServerStopped chan bool
	DB            *gorm.DB
	Dependencies  *dependencies
}

// NewApp initializes the app
func NewApp() *App {
	a := &App{}
	a.setupInfrastructure()
	a.setupDependencies()
	a.setupServer()
	return a
}

func (a *App) setupInfrastructure() {
	a.DB = infrastructure.NewGormPostgresClient()
}

func (a *App) setupDependencies() {
	crRepository := repository.NewClimateRecordsRepository(a.DB)
	crJobsRepository := repository.NewClimateRecordJobsRepository(a.DB)
	crJobsService := service.NewClimateRecordJobsService(crRepository, crJobsRepository)
	crJobsDependencies := &climateRecordJobsDependencies{
		service: crJobsService,
	}
	crService := service.NewClimateRecordsService(crRepository)
	crHandler := handler.NewClimateRecordsHandler(crService)
	crDependencies := &climateRecordsDependencies{
		handler: crHandler,
	}
	a.Dependencies = &dependencies{
		climateRecordJobs: crJobsDependencies,
		climateRecords:    crDependencies,
	}

}

func setupHealthCheckRoute(g *echo.Group) {
	g.GET("/health-check", func(ctx echo.Context) error {
		return ctx.NoContent(200)
	})
}

func (a *App) setupServer() {
	a.Server = echo.New()
	a.Server.HTTPErrorHandler = ErrorHandler
	baseGroup := a.Server.Group(fmt.Sprintf("/api/%s", ApplicationName))
	setupHealthCheckRoute(baseGroup)
	a.setupRoutes(baseGroup)
}

// StartApp initializes the server
func (a *App) StartApp() {
	a.Dependencies.climateRecordJobs.service.CreateClimateRecordJob()
	c := cron.New()
	c.AddFunc("@daily", func() {
		a.Dependencies.climateRecordJobs.service.CreateClimateRecordJob()
	})

	go a.startServer()
	if a.ServerReady != nil {
		a.ServerReady <- true
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	a.stopApp()
}

func (a *App) startServer() {
	port := os.Getenv("PORT")
	if err := a.Server.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Print("Shutting down the server")
	}
}

func (a *App) stopApp() {
	ctx := context.Background()
	if err := a.Server.Shutdown(ctx); err != nil {
		log.Print("Error shutting down the server", "error:", err)
	}
}
