package schemas

type ServeDeploymentReq struct {
	Method       string              `json:"method"`
	URL          string              `json:"url"`
	Body         []byte              `json:"body"`
	Header       map[string][]string `json:"header"`
	DeploymentID string              `json:"deploymentId"`
}

type ServeDeploymentResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	StatusCode int32  `json:"statusCode"`
}
