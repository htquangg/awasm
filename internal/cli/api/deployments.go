package api

import (
	"fmt"
	"net/http"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/go-resty/resty/v2"
)

func CallAddDeployment(
	httpClient *resty.Client,
	req *schemas.AddDeploymentReq,
) (*schemas.AddDeploymentResp, error) {
	var result AwasmResp[*schemas.AddDeploymentResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", USER_AGENT).
		SetBody(req.Data).
		Post(fmt.Sprintf("%s/v1/endpoints/%s/deployments", config.AWASM_URL, req.EndpointID))
	if err != nil {
		return nil, fmt.Errorf("CallAddDeployment: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallAddDeployment: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
