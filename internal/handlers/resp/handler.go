package resp

import (
	std_errors "errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/segmentfault/pacman/errors"
)

// HandleResponse Handle response body
func HandleResponse(c echo.Context, err error, data interface{}) error {
	// no error
	if err == nil {
		return c.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, Success, data))
	}

	var myErr *errors.Error
	// unknown error
	if !std_errors.As(err, &myErr) {
		log.Error().Stack().Err(err)
		return c.JSON(
			http.StatusOK,
			NewRespBody(
				http.StatusInternalServerError,
				UnknownError,
			),
		)
	}

	// log internal server error
	if errors.IsInternalServer(myErr) {
		log.Error().Err(myErr)
	}

	respBody := NewRespBodyFromError(myErr)
	if data != nil {
		respBody.Data = data
	}

	return c.JSON(http.StatusOK, respBody)
}
