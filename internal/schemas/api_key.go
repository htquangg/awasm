package schemas

type AddApiKeyReq struct {
	FriendlyName string `json:"friendlyName"`
	UserID       string `                                 json:"-"`
}

type AddApiKeyResp struct {
	FriendlyName string `json:"friendlyName"`
	Key          string `json:"key"`
	KeyPreview   string `json:"keyPreview"`
	CreatedAt    int64  `json:"createdAt"`
}

type GetApiKeyResp struct {
	FriendlyName string `json:"friendlyName"`
	Key          string `json:"-"`
	KeyPreview   string `json:"keyPreview"`
	UserID       string `json:"userID"`
	CreatedAt    int64  `json:"createdAt"`
}
