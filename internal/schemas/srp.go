package schemas

type GetSRPAttributeReq struct {
	Email string `validate:"required,email" query:"email"`
}

type GetSRPAttributeResp struct {
	SRPUserID string `json:"srpUserId"`
	SRPSalt   string `json:"srpSalt"`
	// MemLimit,OpsLimit,KekSalt are needed to derive the KeyEncryptionKey
	// on the client. Client generates the LoginKey from the KeyEncryptionKey
	// and treat that as UserInputPassword.
	MemLimit int    `json:"memLimit"`
	OpsLimit int    `json:"opsLimit"`
	KekSalt  string `json:"kekSalt"`
}

type SetupSRPAccountSignupReq struct {
	SRPUserID   string `validate:"required" json:"srpUserId"`
	SRPSalt     string `validate:"required" json:"srpSalt"`
	SRPVerifier string `validate:"required" json:"srpVerifier"`
	SRPA        string `validate:"required" json:"srpA"`
}

type SetupSRPAccountSignupResp struct {
	SetupID string `json:"setupId"`
	SRPB    string `json:"srpB"`
}
