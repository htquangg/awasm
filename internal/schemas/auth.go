package schemas

type AccessTokenResp struct {
	Token        string `json:"accessToken"`
	TokenType    string `json:"tokenType"` // Bearer
	ExpiresIn    int    `json:"expiresIn"`
	ExpiresAt    int64  `json:"expiresAt"`
	RefreshToken string `json:"refreshToken"`
}
