package handler

import (
	"context"

	"github.com/htquangg/a-wasm/internal/constants"

	"github.com/labstack/echo/v4"
	"github.com/segmentfault/pacman/i18n"
)

// GetLang get language from header
func GetLang(ctx echo.Context) i18n.Language {
	acceptLanguage := ctx.Request().Header.Get(constants.AcceptLanguageFlag)
	if len(acceptLanguage) == 0 {
		return i18n.DefaultLanguage
	}
	return i18n.Language(acceptLanguage)
}

// GetLangByCtx get language from header
func GetLangByCtx(ctx context.Context) i18n.Language {
	acceptLanguage, ok := ctx.Value(constants.AcceptLanguageFlag).(i18n.Language)
	if ok {
		return acceptLanguage
	}
	return i18n.DefaultLanguage
}
