package schemas

type PublishEndpointReq struct {
	DeploymentID string `validate:"required" json:"deploymentId"`
}

type PublishEndpointResp struct {
	DeploymentID string `json:"deploymentId"`
	URL          string `json:"url"`
}

type ServeEndpointReq struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Body       []byte              `json:"body"`
	Header     map[string][]string `json:"header"`
	EndpointID string              `json:"endpointId"`
}

type ServeEndpointResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	StatusCode int32  `json:"statusCode"`
}
