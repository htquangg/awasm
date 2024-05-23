package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/htquangg/a-wasm/internal/schemas"
)

func CallAddApiKey(
	httpClient *resty.Client,
	req *schemas.AddApiKeyReq,
) (*schemas.AddApiKeyResp, error) {
	var result AwasmResp[*schemas.AddApiKeyResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", UserAgent).
		SetBody(req).
		Post("/v1/api-keys")
	if err != nil {
		return nil, fmt.Errorf("CallCreateApiKey: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallCreateApiKey: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
