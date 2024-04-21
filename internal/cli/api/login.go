package api

import (
	"fmt"
	"net/http"

	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/logger"

	"github.com/go-resty/resty/v2"
)

func CallGetSRPAttribute(
	httpClient *resty.Client,
	req *schemas.GetSRPAttributeReq,
) (*schemas.GetSRPAttributeResp, error) {
	var result AwasmResp[*schemas.GetSRPAttributeResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Get(fmt.Sprintf("/v1/users/srp/attributes?email=%s", req.Email))
	if err != nil {
		return nil, fmt.Errorf("CallGetSRPAttribute: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallGetSRPAttribute: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallChallengeEmailLogin(
	httpClient *resty.Client,
	req *schemas.ChallengeEmailLoginReq,
) (*schemas.ChallengeEmailLoginResp, error) {
	var result AwasmResp[*schemas.ChallengeEmailLoginResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post("/v1/users/auth/email/login/challenge")
	if err != nil {
		return nil, fmt.Errorf("CallChallengeEmailLogin: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallChallengeEmailLogin: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallVerifyEmailLogin(
	httpClient *resty.Client,
	req *schemas.VerifyEmailLoginReq,
) (*schemas.VerifyEmailLoginResp, error) {
	var result AwasmResp[*schemas.VerifyEmailLoginResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post("/v1/users/auth/email/login/verify")
	if err != nil {
		return nil, fmt.Errorf("CallVerifyEmailLogin: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallVerifyEmailLogin: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallIsAuthenticated(
	httpClient *resty.Client,
) bool {
	var result AwasmResp[any]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		Post("/v1/users/auth/check")
	if err != nil {
		return false
	}
	if resp.IsError() || result.Code != http.StatusOK {
		logger.Debugf("CallVerifyEmailLogin: Unsuccessful response [response=%s]", resp)
		return false
	}

	return true
}
