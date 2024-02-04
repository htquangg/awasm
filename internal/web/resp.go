package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type (
	RespStatus uint8
	RespData   interface{}

	Resp struct {
		Code RespStatus `json:"code"`
		Data RespData   `json:"data"`
	}

	RespMessage string

	RespError struct {
		Resp
		Message RespMessage `json:"message"`
	}
)

const (
	StatusUnknown RespStatus = iota
	StatusOK
	StatusBadRequest
	StatusNotFound
	StatusUnauthorized
	StatusForbidden
	StatusInternal
)

func handleSuccessResp(ctx echo.Context, data RespData) error {
	return handleResp(ctx, data, StatusOK, nil, "")
}

func handleBadRequestResp(ctx echo.Context, err error, msg RespMessage) error {
	return handleResp(ctx, nil, StatusBadRequest, err, msg)
}

func handlNotFoundResp(ctx echo.Context, err error, msg RespMessage) error {
	return handleResp(ctx, nil, StatusNotFound, err, msg)
}

func handleUnauthorizedResp(ctx echo.Context, msg RespMessage) error {
	return handleResp(ctx, nil, StatusUnauthorized, nil, msg)
}

func handleForbiddenResp(ctx echo.Context, msg RespMessage) error {
	return handleResp(ctx, nil, StatusForbidden, nil, msg)
}

func handleInternalResp(ctx echo.Context, err error, msg RespMessage) error {
	return handleResp(ctx, nil, StatusForbidden, err, msg)
}

func handleResp(ctx echo.Context, data RespData, status RespStatus, err error, msg RespMessage) error {
	if status == StatusOK {
		resp := new(Resp)
		resp.Code = status
		resp.Data = data
		return ctx.JSON(http.StatusOK, resp)
	}

	resp := new(RespError)
	resp.Message = msg
	resp.Code = status

	switch status {
	case StatusBadRequest, StatusNotFound, StatusUnauthorized, StatusForbidden:
	case StatusInternal:
		log.Error().Err(err).Msg("internal server error")
	default:
		resp.Code = StatusInternal
		log.Warn().Err(err).Msgf("unknown resp status code: %T", status)
	}

	return ctx.JSON(http.StatusOK, resp)
}
