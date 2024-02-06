package web

import (
	"context"
	std_log "log"
	"net/http"
	"strings"

	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/db"
	"github.com/htquangg/a-wasm/internal/services"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	ctx context.Context
	cfg *Config
	e   *echo.Echo
}

func New(
	ctx context.Context,
	cfg *Config,
	svc *services.Service,
	db db.DB,
) *Server {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if err == nil {
			return // no error
		}

		var apiErr *RespError

		if errors.As(err, &apiErr) {
			err = c.JSON(http.StatusOK, apiErr)
			// already an api error...
		} else if v := new(echo.HTTPError); errors.As(err, &v) {
			log.Warn().Err(err).Msgf("[API][ECHO] error: %d %T", v.Code, v.Message)
			apiErr = &RespError{
				Resp: Resp{
					Code: RespStatus(v.Code),
					Data: nil,
				},
				Message: "",
			}
			err = c.JSON(int(apiErr.Code), apiErr)
		} else {
			log.Warn().Err(err).Msgf("[API][UNKNOWN] error: %d %T", v.Code, v.Message)
			apiErr = &RespError{
				Resp: Resp{
					Code: StatusInternalServer,
					Data: nil,
				},
				Message: "",
			}
			err = c.JSON(int(apiErr.Code), apiErr)
		}

		if c.Response().Committed {
			return
		}
	}

	v1 := e.Group("/api/v1", TransactionalMiddleware(db), TransactionalMiddleware(db))

	bindHealthApi(v1)
	bindEndpointsApi(v1, svc)

	// catch all any route
	v1.Any("/*", func(c echo.Context) error {
		return echo.ErrNotFound
	}, otelecho.Middleware("a-wasm"))

	return &Server{
		ctx: ctx,
		cfg: cfg,
		e:   e,
	}
}

func (s *Server) ServeHandler() (execute func() error, interrupt func(error)) {
	server := &http.Server{
		ReadTimeout:       constants.ReadTimeout,
		ReadHeaderTimeout: constants.ReadHeaderTimeout,
		WriteTimeout:      constants.WriteTimeout,
		Handler:           s.e,
		Addr:              s.cfg.Addr,
	}

	if s.cfg.ShowStartBanner {
		addr := server.Addr
		schema := "http"

		date := new(strings.Builder)
		std_log.New(date, "", std_log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Web server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema, addr),
		)

		regular := color.New()
		regular.Printf("├─ REST API: %s\n", color.CyanString("%s://%s/api/", schema, addr))
	}

	return func() error {
			return server.ListenAndServe()
		}, func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Warn().Err(err).Msg("Web server failed to stop gracefully")
			}
		}
}
