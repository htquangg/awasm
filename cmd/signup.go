package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"runtime"

	"github.com/htquangg/a-wasm/internal/cli"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"

	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	secretbox "github.com/GoKillers/libsodium-go/cryptosecretbox"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/kong/go-srp"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
)

const (
	SRP_4096_PARAMS = 4096
)

var signupCmd = &cobra.Command{
	Example:               "awasm signup",
	Use:                   "signup",
	Short:                 "Signup into your Awasm account",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var userCredentialToBeStored schemas.UserCredential
		signupCredential(&userCredentialToBeStored)

		magenta := color.New(color.FgHiMagenta)
		boldMagenta := magenta.Add(color.Bold)
		boldMagenta.Printf(">>>> Welcome to Awasm!\n")
		boldMagenta.Printf("You are now signup as %v <<<< \n", userCredentialToBeStored.Email)
	},
}

func init() {
}

func signupCredential(userCredential *schemas.UserCredential) {
	email, password, err := askForCredential()
	if err != nil {
		cli.HandleError(err, "Unable to parse email and password for authentication")
	}
	// set up resty client
	httpClient := resty.New()
	httpClient.
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	// [1]. Send verification code to email
	_, err = api.CallBeginEmailSignupProcess(httpClient, &schemas.BeginEmailSignupProcessReq{
		Email: email,
	})
	if err != nil {
		cli.HandleError(err)
	}

	// [2]. Verify email
	fmt.Println("Enter OTP...")
	otpPrompt := promptui.Prompt{
		Label: "otp",
		Validate: func(input string) error {
			if len(input) != 6 {
				return errors.New("please enter a valid otp")
			}
			return nil
		},
	}
	otp, err := otpPrompt.Run()
	if err != nil {
		cli.HandleError(err)
	}
	verifyEmailSignupResp, err := api.CallVerifyEmailSignup(httpClient, &schemas.VerifyEmailSignupReq{
		Email: email,
		OTP:   otp,
	})

	// set the jwt token to request
	httpClient.SetAuthToken(verifyEmailSignupResp.AccessToken)

	// [3]. Setup SRP account
	masterKey, keyAttribute, srpAttribute, err := generateKeyAndSRPAttributes(password)
	if err != nil {
		cli.HandleError(err)
	}

	kekSaltBytes, err := convertStringToBytes(keyAttribute.KekSalt)
	if err != nil {
		cli.HandleError(err)
	}
	kekBytes, _, _ := deriveKey(password, kekSaltBytes)
	loginSubKey, err := crypto.GenerateLoginSubKey(convertBytesToString(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	srpClient, err := generateSRPClient(srpAttribute.SRPSalt, srpAttribute.SRPUserID, loginSubKey)
	if err != nil {
		cli.HandleError(err)
	}

	setupSRPAccountSignupResp, err := api.CallSetupSRPAccountSignup(httpClient, &schemas.SetupSRPAccountSignupReq{
		SRPUserID:   srpAttribute.SRPUserID,
		SRPSalt:     srpAttribute.SRPSalt,
		SRPVerifier: srpAttribute.SRPVerifier,
		SRPA:        convertBytesToString(srpClient.ComputeA()),
	})

	// [4]. Complete signup account
	srpBBytes, err := convertStringToBytes(setupSRPAccountSignupResp.SRPB)
	if err != nil {
		cli.HandleError(err)
	}
	srpClient.SetB(srpBBytes)
	srpM1 := convertBytesToString(srpClient.ComputeM1())
	completeEmailAccountSignupResp, err := api.CallCompleteEmailAccountSignup(
		httpClient,
		&schemas.CompleteEmailSignupReq{
			SetupID:      setupSRPAccountSignupResp.SetupID,
			SRPM1:        srpM1,
			KeyAttribute: *keyAttribute,
		},
	)
	if err != nil {
		cli.HandleError(err)
	}

	masterKeyBytes, err := convertStringToBytes(masterKey)
	if err != nil {
		cli.HandleError(err)
	}
	encryptedSecretKeyBytes, err := convertStringToBytes(keyAttribute.EncryptedSecretKey)
	if err != nil {
		cli.HandleError(err)
	}
	keyEncryptionNonceBytes, err := convertStringToBytes(keyAttribute.SecretKeyDecryptionNonce)
	if err != nil {
		cli.HandleError(err)
	}
	privateKey, err := crypto.Decrypt(encryptedSecretKeyBytes, masterKeyBytes, keyEncryptionNonceBytes)
	if err != nil {
		cli.HandleError(err)
	}
	tokenEnc, err := crypto.GetDecryptedToken(
		completeEmailAccountSignupResp.EncryptedToken,
		keyAttribute.PublicKey,
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

	// set the jwt token to request
	httpClient.SetAuthToken(token)

	// updating user credential
	kekEncrypted, err := crypto.GenerateKeyAndEncrypt(convertBytesToString(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	userCredential.Email = email
	userCredential.AccessToken = token
	userCredential.KeyAttribute = keyAttribute
	userCredential.KekEncrypted = &kekEncrypted
}

func askForCredential() (email string, password string, err error) {
	validateEmail := func(input string) error {
		matched, err := regexp.MatchString("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$", input)
		if err != nil || !matched {
			return errors.New("this doesn't look like an email address")
		}
		return nil
	}

	fmt.Println("Enter Credentials...")
	emailPrompt := promptui.Prompt{
		Label:    "Email",
		Validate: validateEmail,
	}
	userEmail, err := emailPrompt.Run()
	if err != nil {
		return "", "", err
	}

	validatePassword := func(input string) error {
		if len(input) < 1 {
			return errors.New("please enter a valid password")
		}
		return nil
	}
	passwordPrompt := promptui.Prompt{
		Label:    "Password",
		Validate: validatePassword,
		Mask:     '*',
	}
	userPassword, err := passwordPrompt.Run()
	if err != nil {
		return "", "", err
	}

	return userEmail, userPassword, nil
}

func generateKeyAndSRPAttributes(
	password string,
) (string, *schemas.KeyAttributeInfo, *schemas.SetupSRPAccountSignupReq, error) {
	masterKeyBytes, err := crypto.GenerateRandomBytes(secretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		return "", nil, nil, err
	}
	masterKey := convertBytesToString(masterKeyBytes)

	recoveryKeyBytes, err := crypto.GenerateRandomBytes(secretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		return "", nil, nil, err
	}
	recoveryKey := convertBytesToString(recoveryKeyBytes)

	kekSaltBytes, err := crypto.GenerateRandomBytes(generichash.CryptoGenericHashBytesMax())
	if err != nil {
		return "", nil, nil, err
	}
	kekBytes, memLimit, opsLimit := deriveKey(password, kekSaltBytes)
	kek := convertBytesToString(kekBytes)

	masterKeyEncryptedWithKek, err := crypto.Encrypt(masterKey, kekBytes)
	if err != nil {
		return "", nil, nil, err
	}

	masterKeyEncryptedWithRecoveryKey, err := crypto.Encrypt(masterKey, recoveryKeyBytes)
	if err != nil {
		return "", nil, nil, err
	}

	recoveryKeyEncryptedWithMasterKey, err := crypto.Encrypt(recoveryKey, masterKeyBytes)
	if err != nil {
		return "", nil, nil, err
	}

	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return "", nil, nil, err
	}
	encryptedKeyPairAttributes, err := crypto.Encrypt(privateKey, masterKeyBytes)
	if err != nil {
		return "", nil, nil, err
	}

	loginSubKey, err := crypto.GenerateLoginSubKey(kek)

	setupSRPAccountSignupReq, err := generateSRPSetupAttribute(loginSubKey)
	if err != nil {
		return "", nil, nil, err
	}

	return masterKey, &schemas.KeyAttributeInfo{
		MemLimit:                          memLimit,
		OpsLimit:                          opsLimit,
		KekSalt:                           convertBytesToString(kekSaltBytes),
		EncryptedKey:                      convertBytesToString(masterKeyEncryptedWithKek.Cipher),
		KeyDecryptionNonce:                convertBytesToString(masterKeyEncryptedWithKek.Nonce),
		PublicKey:                         publicKey,
		EncryptedSecretKey:                convertBytesToString(encryptedKeyPairAttributes.Cipher),
		SecretKeyDecryptionNonce:          convertBytesToString(encryptedKeyPairAttributes.Nonce),
		MasterKeyEncryptedWithRecoveryKey: convertBytesToString(masterKeyEncryptedWithRecoveryKey.Cipher),
		MasterKeyDecryptionNonce:          convertBytesToString(masterKeyEncryptedWithKek.Nonce),
		RecoveryKeyEncryptedWithMasterKey: convertBytesToString(recoveryKeyEncryptedWithMasterKey.Cipher),
		RecoveryKeyDecryptionNonce:        convertBytesToString(recoveryKeyEncryptedWithMasterKey.Nonce),
	}, setupSRPAccountSignupReq, nil
}

func generateSRPSetupAttribute(loginSubKey string) (*schemas.SetupSRPAccountSignupReq, error) {
	srpParams := srp.GetParams(SRP_4096_PARAMS)

	loginSubKeyBytes, err := convertStringToBytes(loginSubKey)
	if err != nil {
		return nil, err
	}

	srpSaltBytes, err := crypto.GenerateRandomBytes(secretbox.CryptoSecretBoxKeyBytes())
	if err != nil {
		return nil, err
	}

	srpUserID := uid.ID()

	srpVerifierBytes := srp.ComputeVerifier(srpParams, srpSaltBytes, []byte(srpUserID), loginSubKeyBytes)

	return &schemas.SetupSRPAccountSignupReq{
		SRPUserID:   srpUserID,
		SRPSalt:     convertBytesToString(srpSaltBytes),
		SRPVerifier: convertBytesToString(srpVerifierBytes),
	}, nil
}

func generateSRPClient(srpSalt string, srpUserID string, loginSubKey string) (*srp.SRPClient, error) {
	clientSecret := srp.GenKey()
	srpParams := srp.GetParams(SRP_4096_PARAMS)

	loginSubKeyBytes, err := convertStringToBytes(loginSubKey)
	if err != nil {
		return nil, err
	}

	srpSaltBytes, err := convertStringToBytes(srpSalt)
	if err != nil {
		return nil, err
	}

	return srp.NewClient(srpParams, srpSaltBytes, []byte(srpUserID), loginSubKeyBytes, clientSecret), nil
}

func convertStringToBytes(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func convertBytesToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func deriveKey(
	passphrase string,
	salt []byte,
) (kekBytes []byte, memLimit int, opsLimit int) {
	memLimit = 64 * 1024
	opsLimit = runtime.NumCPU()
	return argon2.IDKey([]byte(passphrase), salt, 1, uint32(memLimit), uint8(opsLimit), 32), memLimit, opsLimit
}
