package handler

import (
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/i18n"

	"github.com/htquangg/awasm/internal/base/translator"
)

type RespBody struct {
	// response data
	Data interface{} `json:"data"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
	// http code
	Code int `json:"code"`
}

// NewRespBody new response body
func NewRespBody(code int, reason string) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
	}
}

// TrMsg translate the reason cause as a message
func (r *RespBody) TrMsg(lang i18n.Language) *RespBody {
	if len(r.Message) == 0 {
		r.Message = translator.Tr(lang, r.Reason)
	}
	return r
}

func (r RespBody) Error() string {
	return r.Reason
}

// NewRespBodyFromError new response body from error
func NewRespBodyFromError(e *errors.Error) *RespBody {
	return &RespBody{
		Code:    e.Code,
		Reason:  e.Reason,
		Message: e.Message,
	}
}

// NewRespBodyData new response body with data
func NewRespBodyData(code int, reason string, data interface{}) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}
