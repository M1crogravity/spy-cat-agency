package main

import (
	"flag"
	"log/slog"
	"os"
	"sync"

	"github.com/m1crogravity/spy-cat-agency/internal/service"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/memory"
)

type config struct {
	port int
}

type application struct {
	config          config
	logger          *slog.Logger
	wg              sync.WaitGroup
	spyCatsService  *service.SpyCatService
	missionsService *service.MissionsService
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	spyCatsRepo := memory.NewSpyCatRepository()
	spyCatsService := service.NewSpyCatService(spyCatsRepo)
	missionRepo := memory.NewMissionsRepository()
	missionsService := service.NewMissionsService(missionRepo)

	app := &application{
		config:          cfg,
		logger:          logger,
		spyCatsService:  spyCatsService,
		missionsService: missionsService,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
