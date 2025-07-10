package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1crogravity/spy-cat-agency/internal/service"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/postgres"
	"github.com/m1crogravity/spy-cat-agency/internal/storage/remote"
)

type config struct {
	port int
	db   struct {
		dsn         string
		maxOpenCons int
		maxIdleCons int
		maxIdleTime time.Duration
	}
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

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://head_agent:pa55word@localhost:5432/sca?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max-open-const", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleCons, "db-max-idle-const", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	dbPool, err := openDB(context.Background(), cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	spyCatsRepo := postgres.NewSpyCatsRepository(dbPool)
	httpClient := newHttpClient(logger)
	breedsRepo := remote.NewBreedsRepository(httpClient, 5*time.Minute)
	spyCatsService := service.NewSpyCatService(spyCatsRepo, breedsRepo)
	missionRepo := postgres.NewMissionsRepository(dbPool)
	missionsService := service.NewMissionsService(missionRepo)
	tokensRepo := postgres.NewTokensRepository(dbPool)
	tokensService := service.NewTokensService(tokensRepo)
	agentsRepository := postgres.NewAgentsRepository(dbPool)
	agentsService := service.NewAgentsService(agentsRepository)

	app := &application{
		config:          cfg,
		logger:          logger,
		spyCatsService:  spyCatsService,
		missionsService: missionsService,
		tokensService:   tokensService,
		agentsService:   agentsService,
	}

	err = app.serve()
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

func openDB(ctx context.Context, cfg config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.db.maxOpenCons)
	config.MinConns = int32(cfg.db.maxIdleCons)
	config.MaxConnIdleTime = cfg.db.maxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
