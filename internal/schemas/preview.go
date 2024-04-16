package schemas

type ServeDeploymentReq struct {
	Header       map[string][]string `json:"header"`
	Method       string              `json:"method"`
	URL          string              `json:"url"`
	DeploymentID string              `json:"deploymentId"`
	Body         []byte              `json:"body"`
}

type ServeDeploymentResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	StatusCode int32  `json:"statusCode"`
}
