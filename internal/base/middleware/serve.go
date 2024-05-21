package middleware

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/reason"
)

func Serve(ctx echo.Context, status int, data []byte, header []byte) error {
	headers := make(map[string][]string)

	err := json.Unmarshal(header, &headers)
	if err != nil {
		return handler.HandleResponse(
			ctx,
			errors.InternalServer(reason.UnknownError).WithError(err),
			data,
		)
	}

	for k, vv := range headers {
		for _, v := range vv {
			ctx.Response().Header().Add(k, v)
		}
	}

	_, err = ctx.Response().Write(data)
	if err != nil {
		return err
	}

	return nil
}
