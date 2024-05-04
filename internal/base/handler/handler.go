package handler

import (
	std_errors "errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
	myErrors "github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/base/validator"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/pkg/logger"
)

// HandleResponse Handle response body
func HandleResponse(ctx echo.Context, err error, data interface{}) error {
	lang := GetLang(ctx)
	// no error
	if err == nil {
		return ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, reason.Success, data).TrMsg(lang))
	}

	var myErr *errors.Error
	// unknown error
	if !std_errors.As(err, &myErr) {
		logger.Error(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			NewRespBody(
				http.StatusInternalServerError,
				reason.UnknownError,
			).TrMsg(lang),
		)
	}

	// log internal server error
	if isInternalServer(myErr) {
		logger.Error(myErr)
		myErr.Reason = ""
		return ctx.JSON(myErr.Code,
			NewRespBody(myErr.Code, "").TrMsg(lang),
		)
	}

	respBody := NewRespBodyFromError(myErr).TrMsg(lang)
	if data != nil {
		respBody.Data = data
	}

	if http.StatusText(respBody.Code) == respBody.Reason {
		return ctx.JSON(respBody.Code, respBody)
	}

	return ctx.JSON(http.StatusOK, respBody)
}

// Bind bind request and validate
func BindAndValidate(ctx echo.Context, data interface{}) (err error, errField any) {
	lang := GetLang(ctx)
	ctx.Set(constants.AcceptLanguageFlag, lang)
	if err := ctx.Bind(data); err != nil {
		logger.Errorf("http_handle BindAndCheck fail, %s", err.Error())
		return myErrors.New(http.StatusBadRequest, reason.RequestFormatError), nil
	}

	errField, err = validator.GetValidatorByLang(lang).Check(data)

	return err, errField
}

func isInternalServer(err *myErrors.Error) bool {
	return err.Code >= http.StatusInternalServerError
}
