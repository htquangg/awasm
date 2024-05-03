package web

import (
	"context"
	std_errors "errors"
	"fmt"
	std_log "log"
	"net/http"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/middleware"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/controllers"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/fatih/color"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	middleware_echo "github.com/labstack/echo/v4/middleware"
	"github.com/segmentfault/pacman/errors"
)

type Server struct {
	ctx context.Context
	cfg *config.Server
	e   *echo.Echo
	srv *http.Server
}

func New(
	ctx context.Context,
	cfg *config.Server,
	controllers *controllers.Controllers,
	mws *middleware.Middleware,
) *Server {
	e := echo.New()

	e.Pre(
		echoprometheus.NewMiddleware("awasm"),
		middleware_echo.RemoveTrailingSlashWithConfig(middleware_echo.TrailingSlashConfig{
			Skipper: func(c echo.Context) bool {
				// enable by default only for the API routes
				return !strings.HasPrefix(c.Request().URL.Path, constants.LiveIngressPath) ||
					!strings.HasPrefix(c.Request().URL.Path, constants.PreviewIngressPath) ||
					!strings.HasPrefix(c.Request().URL.Path, "/api/")
			},
		}),
		middleware_echo.RequestIDWithConfig(
			middleware_echo.RequestIDConfig{
				Generator: func() string {
					return uid.ID()
				},
			}),
		middleware_echo.Logger(),
		middleware_echo.Gzip(),
		middleware_echo.Recover(),
		middleware_echo.Secure(),
	)

	e.GET("/metrics", echoprometheus.NewHandler())

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		if v := new(echo.HTTPError); std_errors.As(err, &v) {
			reason := ""
			if v.Internal != nil {
				reason = v.Error()
			} else {
				reason = fmt.Sprintf("%v", v.Message)
			}
			err = errors.New(v.Code, reason)
		}

		_ = handler.HandleResponse(ctx, err, nil)
	}

	v1Group := e.Group("/api/v1")

	bindHealthApi(v1Group, controllers)
	bindEndpointsApi(v1Group, controllers, mws)
	bindPreviewApi(v1Group, controllers, mws)
	bindLiveApi(v1Group, controllers, mws)
	bindUserApi(v1Group, controllers, mws)
	bindApiKeyApi(v1Group, controllers, mws)

	// catch all any route
	v1Group.Any("/*", func(ctx echo.Context) error {
		return echo.ErrNotFound
	})

	srv := &http.Server{
		ReadTimeout:       constants.ReadTimeout,
		ReadHeaderTimeout: constants.ReadHeaderTimeout,
		WriteTimeout:      constants.WriteTimeout,
		Handler:           e,
		Addr:              cfg.Addr,
	}

	return &Server{
		ctx: ctx,
		cfg: cfg,
		e:   e,
		srv: srv,
	}
}

func (s *Server) Start() error {
	if s.cfg.ShowStartBanner {
		addr := s.srv.Addr
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

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	if s.srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("Web server failed to stop gracefully: %v", err)
		}

		return nil
	}

	return nil
}
