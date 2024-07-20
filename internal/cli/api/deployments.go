package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/htquangg/awasm/internal/schemas"
)

func CallAddDeployment(
	httpClient *resty.Client,
	req *schemas.AddDeploymentReq,
) (*schemas.AddDeploymentResp, error) {
	var result AwasmResp[*schemas.AddDeploymentResp]
	resp, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("User-Agent", UserAgent).
		SetBody(req.Data).
		Post(fmt.Sprintf("/v1/endpoints/%s/deployments", req.EndpointID))
	if err != nil {
		return nil, fmt.Errorf("CallAddDeployment: Unable to complete api request [err=%s]", err)
	}
	if resp.IsError() || result.Code != http.StatusOK {
		return nil, fmt.Errorf("CallAddDeployment: Unsuccessful response [response=%s]", resp)
	}

	return result.Data, nil
}
