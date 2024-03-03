package handler

import (
	std_errors "errors"
	"net/http"

	"github.com/htquangg/a-wasm/internal/base/reason"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

// HandleResponse Handle response body
func HandleResponse(ctx echo.Context, err error, data interface{}) error {
	// no error
	if err == nil {
		return ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, reason.Success, data))
	}

	var myErr *errors.Error
	// unknown error
	if !std_errors.As(err, &myErr) {
		log.Error(err)
		return ctx.JSON(
			http.StatusOK,
			NewRespBody(
				http.StatusInternalServerError,
				reason.UnknownError,
			),
		)
	}

	// log internal server error
	if errors.IsInternalServer(myErr) {
		log.Error(myErr)
		myErr.Reason = ""
	}

	respBody := NewRespBodyFromError(myErr)
	if data != nil {
		respBody.Data = data
	}

	return ctx.JSON(http.StatusOK, respBody)
}
