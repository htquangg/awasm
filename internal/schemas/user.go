package schemas

import (
	"github.com/jinzhu/copier"

	"github.com/htquangg/awasm/internal/entities"
)

type BeginEmailSignupProcessReq struct {
	Email string `validate:"required,email" json:"email"`
}

type BeginEmailSignupProcessResp struct{}

type VerifyEmailSignupReq struct {
	Email     string `validate:"required,email"       json:"email"`
	OTP       string `validate:"required,gte=6,lte=6" json:"otp"`
	IP        string `                                json:"-"`
	UserAgent string `                                json:"-"`
}

type VerifyEmailSignupResp struct {
	Token string `json:"token"`
}

type CompleteEmailSignupReq struct {
	SetupID      string           `validate:"required" json:"setupId"`
	SRPM1        string           `validate:"required" json:"srpM1"`
	IP           string           `                    json:"-"`
	UserAgent    string           `                    json:"-"`
	KeyAttribute KeyAttributeInfo `validate:"required" json:"keyAttribute"`
}

type CompleteEmailSignupInfo struct {
	UserID       string
	SRPUserID    string
	Salt         string
	Verifier     string
	KeyAttribute KeyAttributeInfo
}

type CompleteEmailSignupResp struct {
	KeyAttribute *KeyAttributeInfo `json:"keyAttribute"`
	SRPM2        string            `json:"srpM2"`
}

type ChallengeEmailLoginReq struct {
	SRPUserID string `validate:"required" json:"srpUserId"`
	SRPA      string `validate:"required" json:"srpA"`
}

type ChallengeEmailLoginResp struct {
	ChallengeID string `json:"challengeId" binding:"required"`
	SRPB        string `json:"srpB"        binding:"required"`
}

type VerifyEmailLoginReq struct {
	ChallengeID string `validate:"required" json:"challengeId"`
	SRPUserID   string `validate:"required" json:"srpUserId"`
	SRPM1       string `validate:"required" json:"srpM1"`
}

type VerifyEmailLoginResp struct {
	*AccessTokenResp
	KeyAttribute *KeyAttributeInfo `json:"keyAttribute"`
	SRPM2        string            `json:"srpM2"`
}

type KeyAttributeInfo struct {
	KekSalt                           string `validate:"required" json:"kekSalt"`
	EncryptedKey                      string `validate:"required" json:"encryptedKey"`
	KeyDecryptionNonce                string `validate:"required" json:"keyDecryptionNonce"`
	PublicKey                         string `validate:"required" json:"publicKey"`
	EncryptedSecretKey                string `validate:"required" json:"encryptedSecretKey"`
	SecretKeyDecryptionNonce          string `validate:"required" json:"secretKeyDecryptionNonce"`
	MasterKeyEncryptedWithRecoveryKey string `                    json:"masterKeyEncryptedWithRecoveryKey"`
	MasterKeyDecryptionNonce          string `                    json:"masterKeyDecryptionNonce"`
	RecoveryKeyEncryptedWithMasterKey string `                    json:"recoveryKeyEncryptedWithMasterKey"`
	RecoveryKeyDecryptionNonce        string `                    json:"recoveryKeyDecryptionNonce"`
	MemLimit                          int    `validate:"required" json:"memLimit"`
	OpsLimit                          int    `validate:"required" json:"opsLimit"`
}

func (i *KeyAttributeInfo) ConvertFromKeyAttributeEntity(keyAttributeInfo *entities.KeyAttribute) {
	_ = copier.Copy(i, keyAttributeInfo)
}

type EncryptionResult struct {
	Cipher []byte
	Nonce  []byte
	Key    []byte
}

type UserCredential struct {
	KeyAttribute *KeyAttributeInfo `json:"keyAttribute"`
	KekEncrypted *EncryptionResult `json:"kekEncrypted"`
	Email        string            `json:"email"`
	AccessToken  string            `json:"accessToken"`
}

type ConfigFile struct {
	LoggedInUserEmail  string          `json:"loggedInUserEmail"`
	LoggedInUserDomain string          `json:"LoggedInUserDomain,omitempty"`
	LoggedInUsers      []*LoggedInUser `json:"loggedInUsers,omitempty"`
}

type LoggedInUser struct {
	Email  string `json:"email"`
	Domain string `json:"domain"`
}

type RefreshTokenReq struct {
	RefreshToken string `validate:"required" json:"refreshToken"`
	IP           string `                    json:"-"`
	UserAgent    string `                    json:"-"`
}

type RefreshTokenResp struct {
	*AccessTokenResp
}
