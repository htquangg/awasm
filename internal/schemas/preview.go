package schemas

type ServeDeploymentReq struct {
	Header       map[string][]string `json:"header"`
	Method       string              `json:"method"`
	URL          string              `json:"url"`
	DeploymentID string              `json:"deploymentId"`
	UserID       string              `json:"-"`
	Body         []byte              `json:"body"`
}

type ServeDeploymentResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	Header     []byte `json:"-"`
	StatusCode int32  `json:"statusCode"`
}
