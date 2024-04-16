package schemas

type PublishEndpointReq struct {
	DeploymentID string `validate:"required" json:"deploymentId"`
	UserID       string `                    json:"-"`
}

type PublishEndpointResp struct {
	DeploymentID string `json:"deploymentId"`
	IngressURL   string `json:"ingressUrl"`
}

type ServeEndpointReq struct {
	Header     map[string][]string `json:"header"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	EndpointID string              `json:"endpointId"`
	Body       []byte              `json:"body"`
}

type ServeEndpointResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	StatusCode int32  `json:"statusCode"`
}
