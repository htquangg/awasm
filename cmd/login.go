package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/crypto"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Example: "awasm login",
	Use:     "login",
	Short:   "Login into your Awasm account",
	Run: func(cmd *cobra.Command, args []string) {
		var userCredentialToBeStored schemas.UserCredential
		loginCredential(&userCredentialToBeStored)

		fmt.Printf(">>>> Welcome to Awasm!\n")
	},
}

func init() {
}

func loginCredential(userCredential *schemas.UserCredential) {
	email, password, err := askForCredential()
	if err != nil {
		cli.HandleError(err, "Unable to parse email and password for authentication")
	}

	// set up resty client
	httpClient := resty.New()
	httpClient.
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	// [1]. Get srp attribute
	getSRPAttributeResp, err := api.CallGetSRPAttribute(httpClient, &schemas.GetSRPAttributeReq{
		Email: email,
	})
	if err != nil {
		cli.HandleError(err)
	}

	// [2]. Challenge email login
	kekSaltBytes, err := convertStringToBytes(getSRPAttributeResp.KekSalt)
	if err != nil {
		cli.HandleError(err)
	}
	kekBytes, _, _ := deriveKey(password, kekSaltBytes)
	loginSubKey, err := crypto.GenerateLoginSubKey(convertBytesToString(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	srpClient, err := generateSRPClient(getSRPAttributeResp.SRPSalt, getSRPAttributeResp.SRPUserID, loginSubKey)
	if err != nil {
		cli.HandleError(err)
	}

	challengeEmailLoginResp, err := api.CallChallengeEmailLogin(httpClient, &schemas.ChallengeEmailLoginReq{
		SRPUserID: getSRPAttributeResp.SRPUserID,
		SRPA:      convertBytesToString(srpClient.ComputeA()),
	})
	if err != nil {
		cli.HandleError(err)
	}

	// [3]. Verify email login
	srpBBytes, err := convertStringToBytes(challengeEmailLoginResp.SRPB)
	if err != nil {
		cli.HandleError(err)
	}
	srpClient.SetB(srpBBytes)
	srpM1 := convertBytesToString(srpClient.ComputeM1())
	verifyEmailLoginResp, err := api.CallVerifyEmailLogin(httpClient, &schemas.VerifyEmailLoginReq{
		ChallengeID: challengeEmailLoginResp.ChallengeID,
		SRPUserID:   getSRPAttributeResp.SRPUserID,
		SRPM1:       srpM1,
	})
	if err != nil {
		cli.HandleError(err)
	}

	// get access token
	encryptedKeyBytes, err := convertStringToBytes(verifyEmailLoginResp.KeyAttribute.EncryptedKey)
	if err != nil {
		cli.HandleError(err)
	}
	decryptionKeyNonceBytes, err := convertStringToBytes((verifyEmailLoginResp.KeyAttribute.KeyDecryptionNonce))
	if err != nil {
		cli.HandleError(err)
	}
	masterKey, err := crypto.Decrypt(encryptedKeyBytes, kekBytes, decryptionKeyNonceBytes)
	if err != nil {
		cli.HandleError(err)
	}
	masterKeyBytes, err := convertStringToBytes(masterKey)
	if err != nil {
		cli.HandleError(err)
	}
	encryptedSecretKeyBytes, err := convertStringToBytes(verifyEmailLoginResp.KeyAttribute.EncryptedSecretKey)
	if err != nil {
		cli.HandleError(err)
	}
	keyEncryptionNonceBytes, err := convertStringToBytes(verifyEmailLoginResp.KeyAttribute.SecretKeyDecryptionNonce)
	if err != nil {
		cli.HandleError(err)
	}
	privateKey, err := crypto.Decrypt(encryptedSecretKeyBytes, masterKeyBytes, keyEncryptionNonceBytes)
	if err != nil {
		cli.HandleError(err)
	}
	tokenEnc, err := crypto.GetDecryptedToken(
		verifyEmailLoginResp.EncryptedToken,
		verifyEmailLoginResp.KeyAttribute.PublicKey,
		privateKey,
	)
	if err != nil {
		cli.HandleError(err)
	}
	tokenEncBytes, err := convertStringToBytes(tokenEnc)
	if err != nil {
		cli.HandleError(err)
	}
	token := string(tokenEncBytes)

	// updating user credential
	kekEncrypted, err := crypto.GenerateKeyAndEncrypt(convertBytesToString(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}

	userCredential.Email = email
	userCredential.AccessToken = token
	userCredential.KeyAttribute = verifyEmailLoginResp.KeyAttribute
	userCredential.KekEncrypted = &kekEncrypted
}
