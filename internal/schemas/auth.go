package schemas

type CommonTokenResp struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"` // Bearer
	ExpiresIn   int    `json:"expiresIn"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type AccessTokenResp struct {
	EncryptedToken string `json:"encryptedToken"`
	RefreshToken   string `json:"refreshToken"`
	CommonTokenResp
}
