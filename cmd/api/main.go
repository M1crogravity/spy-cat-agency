package main

import (
	"bytes"
	"flag"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/m1crogravity/spy-cat-agency/internal/service"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/memory"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/remote"
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
	tokensService   *service.TokenService
	agentsService   *service.AgentsService
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	spyCatsRepo := memory.NewSpyCatRepository()
	httpClient := newHttpClient(logger)
	breedsRepo := remote.NewBreedsRepository(httpClient, 5*time.Minute)
	spyCatsService := service.NewSpyCatService(spyCatsRepo, breedsRepo)
	missionRepo := memory.NewMissionsRepository()
	missionsService := service.NewMissionsService(missionRepo)
	tokensRepo := memory.NewTokensRepository()
	tokensService := service.NewTokensService(tokensRepo)
	agentsRepository := memory.NewAgentsRepository()
	agentsService := service.NewAgentsService(agentsRepository)

	app := &application{
		config:          cfg,
		logger:          logger,
		spyCatsService:  spyCatsService,
		missionsService: missionsService,
		tokensService:   tokensService,
		agentsService:   agentsService,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

type LoggableHttpClient struct {
	logger *slog.Logger
	http.Client
}

func (c *LoggableHttpClient) Do(r *http.Request) (*http.Response, error) {
	received := time.Now()
	bodyBytes := []byte{}

	if r.Body != nil {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	c.logger.Info("->request:",
		"method", r.Method,
		"url", r.URL.String(),
		"body", string(bodyBytes),
		"sent_at", received.Format(time.RFC3339),
	)
	resp, err := c.Client.Do(r)

	if err == nil {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))

		c.logger.Info("<-response:",
			"code", resp.StatusCode,
			"body", string(respBodyBytes),
			"took", time.Since(received),
		)
	} else {
		c.logger.Error(err.Error())
	}

	return resp, err
}

func newHttpClient(logger *slog.Logger) *LoggableHttpClient {
	return &LoggableHttpClient{
		logger: logger,
		Client: http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   time.Second,
				ResponseHeaderTimeout: time.Second,
			},
		},
	}
}
