package web

import (
	"context"
	std_errors "errors"
	std_log "log"
	"net/http"
	"strings"

	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/handlers"
	"github.com/htquangg/a-wasm/internal/handlers/resp"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	handlers *handlers.Handlers,
	mws ...echo.MiddlewareFunc,
) *Server {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if v := new(echo.HTTPError); std_errors.As(err, &v) {
			err = c.JSON(v.Code, resp.NewRespBody(
				v.Code,
				"",
			))
		}

		resp.HandleResponse(c, err, nil)
	}

	v1 := e.Group("/api/v1", mws...)

	bindHealthApi(v1, handlers)
	bindEndpointsApi(v1)

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

	return func() error {
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

			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Warn().Err(err).Msg("Web server failed to stop gracefully")
			}
		}
}
