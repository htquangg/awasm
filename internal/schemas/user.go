package schemas

type SignUpReq struct {
	Email     string `validate:"required,email"       json:"email"`
	OTT       string `validate:"required,gte=6,lte=6" json:"ott"`
	IP        string `                                json:"-"`
	UserAgent string `                                json:"-"`
}

type SignUpResp struct {
	Token string `json:"token"`
}

type EncryptionResult struct {
	Cipher []byte
	Nonce  []byte
}
