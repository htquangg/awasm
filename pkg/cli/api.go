package cli

import (
	"fmt"
	"net/http"

	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/go-resty/resty/v2"
)

const (
	USER_AGENT = "cli"
	AWASM_URL  = "http://127.0.0.1:3000"
)

func CallBeginEmailSignupProcess(
	httpClient *resty.Client,
	req *schemas.BeginEmailSignupProcessReq,
) (*schemas.BeginEmailSignupProcessResp, error) {
	var result AwasmResp[*schemas.BeginEmailSignupProcessResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/api/v1/users/auth/email/signup", AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallBeginEmailSignupProcess: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallBeginEmailSignupProcess: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallVerifyEmailSignup(
	httpClient *resty.Client,
	req *schemas.VerifyEmailSignupReq,
) (*schemas.CommonTokenResp, error) {
	var result AwasmResp[*schemas.CommonTokenResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/api/v1/users/auth/email/verify", AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallVerifyEmailSignup: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallVerifyEmailSignup: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallSetupSRPAccountSignup(
	httpClient *resty.Client,
	req *schemas.SetupSRPAccountSignupReq,
) (*schemas.SetupSRPAccountSignupResp, error) {
	var result AwasmResp[*schemas.SetupSRPAccountSignupResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/api/v1/users/auth/email/setup-srp", AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallSetupSRPAccountSignup: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallSetupSRPAccountSignup: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallCompleteEmailAccountSignup(
	httpClient *resty.Client,
	req *schemas.CompleteEmailSignupReq,
) (*schemas.CompleteEmailSignupResp, error) {
	var result AwasmResp[*schemas.CompleteEmailSignupResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/api/v1/users/auth/email/complete", AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallCompleteEmailAccountSignup: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallCompleteEmailAccountSignup: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
