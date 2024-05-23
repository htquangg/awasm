package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/htquangg/a-wasm/internal/schemas"
)

func CallAddEndpoint(
	httpClient *resty.Client,
	req *schemas.AddEndpointReq,
) (*schemas.AddEndpointResp, error) {
	var result AwasmResp[*schemas.AddEndpointResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", UserAgent).
		SetBody(req).
		Post("/v1/endpoints")
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
		SetHeader("User-Agent", UserAgent).
		SetBody(req).
		Post("/v1/live/publish")
	if err != nil {
		return nil, fmt.Errorf("CallPublishEndpoint: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallPublishEndpoint: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
