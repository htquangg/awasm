package handler

import (
	std_errors "errors"
	"net/http"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/base/validator"
	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/errors"
	myErrors "github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
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
		log.Error(err)
		return ctx.JSON(
			http.StatusOK,
			NewRespBody(
				http.StatusInternalServerError,
				reason.UnknownError,
			).TrMsg(lang),
		)
	}

	// log internal server error
	if errors.IsInternalServer(myErr) {
		log.Error(myErr)
		myErr.Reason = ""
	}

	respBody := NewRespBodyFromError(myErr).TrMsg(lang)
	if data != nil {
		respBody.Data = data
	}

	return ctx.JSON(http.StatusOK, respBody)
}

// Bind bind request and validate
func BindAndValidate(ctx echo.Context, data interface{}) (err error, errField any) {
	lang := GetLang(ctx)
	ctx.Set(constants.AcceptLanguageFlag, lang)
	if err := ctx.Bind(data); err != nil {
		log.Errorf("http_handle BindAndCheck fail, %s", err.Error())
		return myErrors.New(http.StatusBadRequest, reason.RequestFormatError), nil
	}

	errField, err = validator.GetValidatorByLang(lang).Check(data)

	return err, errField
}
