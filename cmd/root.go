package cmd

import (
	"context"
	"os"
	"strings"
	"syscall"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/cluster"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services"
	"github.com/htquangg/a-wasm/internal/web"

	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var g run.Group

func init() {
	initLog()
}

func Execute() {
	if err := execute(); err != nil {
		log.Error().Err(err).Msg("the service exitted abnormally")
		os.Exit(1)
	}
}

func execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := db.New(ctx, cfg.DB)
	if err != nil {
		return err
	}

	svc := services.New(db)

	g.Add(web.
		New(ctx, &web.Config{
			ShowStartBanner: true,
			Addr:            cfg.Server.HTTP.Addr,
		},
			svc,
			db,
		).
		ServeHandler(),
	)

	g.Add(cluster.New(ctx).ServeHandler())

	g.Add(run.SignalHandler(ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	))

	var se run.SignalError
	if err := g.Run(); err != nil && !errors.As(err, &se) {
		return err
	}

	return nil
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logLevel := os.Getenv(constants.LogLevel)
	switch strings.ToLower(logLevel) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
