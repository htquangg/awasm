package cmd

import (
	"fmt"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/converter"
	"github.com/htquangg/a-wasm/pkg/crypto"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/htquangg/a-wasm/pkg/logger"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	ADD_USER       = "Add a new account login"
	REPLACE_USER   = "Override current logged in user"
	EXIT_USER_MENU = "Exit"
)

var loginCmd = &cobra.Command{
	Example:               "awasm login",
	Use:                   "login",
	Short:                 "Login into your Awasm account",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		currentLoggedInUserDetails, isAuthenticated, err := cli.GetCurrentLoggedInUserDetails()
		if err != nil && (strings.Contains(err.Error(), "we couldn't find your logged in details")) {
			logger.Debug(err)
		} else if err != nil {
			cli.HandleError(err)
		}

		if isAuthenticated {
			shouldOverride, err := userLoginMenu(currentLoggedInUserDetails.UserCredentials.Email)
			if err != nil {
				cli.HandleError(err)
			}

			if err != nil {
				cli.HandleError(err)
			}

			if !shouldOverride {
				return
			}
		}

		var userCredentialToBeStored schemas.UserCredential
		loginCredential(&userCredentialToBeStored)

		err = cli.StoreUserCredsInKeyRing(&userCredentialToBeStored)
		if err != nil {
			logger.Errorf("Unable to store your credential in system [%s]", err)
			cli.HandleError(err)
		}

		err = cli.WriteInitalConfig(&userCredentialToBeStored)
		if err != nil {
			cli.HandleError(err, "Unable to write write to Awasm Config file. Please try again")
		}

		green := color.New(color.FgGreen)
		boldGreen := green.Add(color.Bold)
		boldGreen.Printf(">>>> Welcome to Awasm!\n")
		boldGreen.Printf("You are now logged in as %v <<<< \n", userCredentialToBeStored.Email)
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
	kekSaltBytes, err := converter.FromB64(getSRPAttributeResp.KekSalt)
	if err != nil {
		cli.HandleError(err)
	}
	kekBytes, _, _ := deriveKey(password, kekSaltBytes)
	loginSubKey, err := crypto.GenerateLoginSubKey(converter.ToB64(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	srpClient, err := generateSRPClient(getSRPAttributeResp.SRPSalt, getSRPAttributeResp.SRPUserID, loginSubKey)
	if err != nil {
		cli.HandleError(err)
	}

	challengeEmailLoginResp, err := api.CallChallengeEmailLogin(httpClient, &schemas.ChallengeEmailLoginReq{
		SRPUserID: getSRPAttributeResp.SRPUserID,
		SRPA:      converter.ToB64(srpClient.ComputeA()),
	})
	if err != nil {
		cli.HandleError(err)
	}

	// [3]. Verify email login
	srpBBytes, err := converter.FromB64(challengeEmailLoginResp.SRPB)
	if err != nil {
		cli.HandleError(err)
	}
	srpClient.SetB(srpBBytes)
	srpM1 := converter.ToB64(srpClient.ComputeM1())
	verifyEmailLoginResp, err := api.CallVerifyEmailLogin(httpClient, &schemas.VerifyEmailLoginReq{
		ChallengeID: challengeEmailLoginResp.ChallengeID,
		SRPUserID:   getSRPAttributeResp.SRPUserID,
		SRPM1:       srpM1,
	})
	if err != nil {
		cli.HandleError(err)
	}

	// get access token
	encryptedKeyBytes, err := converter.FromB64(verifyEmailLoginResp.KeyAttribute.EncryptedKey)
	if err != nil {
		cli.HandleError(err)
	}
	decryptionKeyNonceBytes, err := converter.FromB64(
		(verifyEmailLoginResp.KeyAttribute.KeyDecryptionNonce),
	)
	if err != nil {
		cli.HandleError(err)
	}
	masterKey, err := crypto.Decrypt(encryptedKeyBytes, kekBytes, decryptionKeyNonceBytes)
	if err != nil {
		cli.HandleError(err)
	}
	masterKeyBytes, err := converter.FromB64(masterKey)
	if err != nil {
		cli.HandleError(err)
	}
	encryptedSecretKeyBytes, err := converter.FromB64(verifyEmailLoginResp.KeyAttribute.EncryptedSecretKey)
	if err != nil {
		cli.HandleError(err)
	}
	keyEncryptionNonceBytes, err := converter.FromB64(
		verifyEmailLoginResp.KeyAttribute.SecretKeyDecryptionNonce,
	)
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
	tokenEncBytes, err := converter.FromB64(tokenEnc)
	if err != nil {
		cli.HandleError(err)
	}
	token := string(tokenEncBytes)

	// updating user credential
	kekEncrypted, err := crypto.GenerateKeyAndEncrypt(converter.ToB64(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}

	userCredential.Email = email
	userCredential.AccessToken = token
	userCredential.KeyAttribute = verifyEmailLoginResp.KeyAttribute
	userCredential.KekEncrypted = &kekEncrypted
}

func userLoginMenu(currentLoggedInUserEmail string) (bool, error) {
	label := fmt.Sprintf("Current logged in user email: %s on domain: %s", currentLoggedInUserEmail, config.AWASM_URL)

	prompt := promptui.Select{
		Label: label,
		Items: []string{ADD_USER, REPLACE_USER, EXIT_USER_MENU},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result != EXIT_USER_MENU, err
}
