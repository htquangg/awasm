package api

import (
	"fmt"
	"net/http"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/go-resty/resty/v2"
)

func CallAddEndpoint(
	httpClient *resty.Client,
	req *schemas.AddEndpointReq,
) (*schemas.AddEndpointResp, error) {
	var result AwasmResp[*schemas.AddEndpointResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/v1/endpoints", config.AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallCreateEndpoint: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallCreateEndpoint: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}

func CallPublishEndpoint(
	httpClient *resty.Client,
	req *schemas.PublishEndpointReq,
) (*schemas.PublishEndpointResp, error) {
	var result AwasmResp[*schemas.PublishEndpointResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req).
		Post(fmt.Sprintf("%s/v1/live/publish", config.AWASM_URL))
	if err != nil {
		return nil, fmt.Errorf("CallPublishEndpoint: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallPublishEndpoint: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
