package schemas

import "encoding/json"

type EmailSourceType string

const (
	EmailVerificationSourceType EmailSourceType = "email-verification"
)

type EmailCodeContent struct {
	SourceType EmailSourceType `json:"sourceType"`
	Code       string          `json:"string"`
	ExpiresAt  int64           `json:"expiresAt"`
	TriesLeft  int             `json:"triesLeft"`
}

func (r *EmailCodeContent) ToJSONString() string {
	codeBytes, _ := json.Marshal(r)
	return string(codeBytes)
}

func (r *EmailCodeContent) FromJSONString(data string) error {
	return json.Unmarshal([]byte(data), &r)
}

type EmailVerificationTemplateData struct {
	Code string
}
