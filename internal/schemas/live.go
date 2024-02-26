package schemas

type ServeLiveReq struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Body       []byte              `json:"body"`
	Header     map[string][]string `json:"header"`
	EndpointID string              `json:"endpointId"`
}

type ServeLiveResp struct {
	RequestID  string `json:"requestId"`
	Response   []byte `json:"response"`
	StatusCode int32  `json:"statusCode"`
}
