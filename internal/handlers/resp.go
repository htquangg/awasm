package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type (
	RespStatus uint
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
	StatusInternalServer
)

func (s RespStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "BAD REQUEST"
	case StatusNotFound:
		return "NOT FOUND"
	case StatusUnauthorized:
		return "UNAUTHORIZED"
	case StatusForbidden:
		return "FORBIDDEN"
	case StatusInternalServer:
		return "INTERNAL SERVER"
	}

	return ""
}

func (e *RespError) Error() string {
	return string(e.Message)
}

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

func handleInternalServerResp(ctx echo.Context, err error, msg RespMessage) error {
	return handleResp(ctx, nil, StatusForbidden, err, msg)
}

func handleResp(ctx echo.Context, data RespData, status RespStatus, err error, msg RespMessage) error {
	if status == StatusOK {
		resp := new(Resp)
		resp.Code = status
		resp.Data = data
		return ctx.JSON(http.StatusOK, resp)
	}

	respErr := new(RespError)
	respErr.Message = msg
	respErr.Code = status

	switch status {
	case StatusBadRequest, StatusNotFound, StatusUnauthorized, StatusForbidden:
	case StatusInternalServer:
		log.Error().Err(err).Msg("internal server error")
	default:
		respErr.Code = StatusInternalServer
		log.Warn().Err(err).Msgf("unknown resp status code: %T", status)
	}

	return respErr
}
