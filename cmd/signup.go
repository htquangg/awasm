package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"

	"github.com/fatih/color"
	"github.com/kong/go-srp"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"

	"github.com/htquangg/awasm/internal/cli"
	"github.com/htquangg/awasm/internal/cli/api"
	"github.com/htquangg/awasm/internal/schemas"
	"github.com/htquangg/awasm/pkg/converter"
	"github.com/htquangg/awasm/pkg/crypto"
	"github.com/htquangg/awasm/pkg/uid"
)

const (
	Srp4096Params = 4096
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

	client := api.NewClient(&api.ClientOptions{
		Debug: viper.GetBool("cli.debug"),
	})

	// check password strength
	err = api.CallCheckPasswordStrength(client.HTTPClient, password)
	if err != nil {
		cli.HandleError(err)
	}

	// [1]. Send verification code to email
	_, err = api.CallBeginEmailSignupProcess(client.HTTPClient, &schemas.BeginEmailSignupProcessReq{
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
	verifyEmailSignupResp, err := api.CallVerifyEmailSignup(
		client.HTTPClient,
		&schemas.VerifyEmailSignupReq{
			Email: email,
			OTP:   otp,
		},
	)
	if err != nil {
		cli.HandleError(err)
	}

	// set the jwt token to request
	client.HTTPClient.SetAuthToken(verifyEmailSignupResp.AccessToken)

	// [3]. Setup SRP account
	masterKey, keyAttribute, srpAttribute, err := generateKeyAndSRPAttributes(password)
	if err != nil {
		cli.HandleError(err)
	}

	kekSaltBytes, err := converter.FromB64(keyAttribute.KekSalt)
	if err != nil {
		cli.HandleError(err)
	}
	kekBytes, _, _ := deriveKey(password, kekSaltBytes)
	loginSubKey, err := crypto.GenerateLoginSubKey(converter.ToB64(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	srpClient, err := generateSRPClient(srpAttribute.SRPSalt, srpAttribute.SRPUserID, loginSubKey)
	if err != nil {
		cli.HandleError(err)
	}

	setupSRPAccountSignupResp, err := api.CallSetupSRPAccountSignup(
		client.HTTPClient,
		&schemas.SetupSRPAccountSignupReq{
			SRPUserID:   srpAttribute.SRPUserID,
			SRPSalt:     srpAttribute.SRPSalt,
			SRPVerifier: srpAttribute.SRPVerifier,
			SRPA:        converter.ToB64(srpClient.ComputeA()),
		},
	)
	if err != nil {
		cli.HandleError(err)
	}

	// [4]. Complete signup account
	srpBBytes, err := converter.FromB64(setupSRPAccountSignupResp.SRPB)
	if err != nil {
		cli.HandleError(err)
	}
	srpClient.SetB(srpBBytes)
	srpM1 := converter.ToB64(srpClient.ComputeM1())
	_, err = api.CallCompleteEmailAccountSignup(
		client.HTTPClient,
		&schemas.CompleteEmailSignupReq{
			SetupID:      setupSRPAccountSignupResp.SetupID,
			SRPM1:        srpM1,
			KeyAttribute: *keyAttribute,
		},
	)
	if err != nil {
		cli.HandleError(err)
	}

	masterKeyBytes, err := converter.FromB64(masterKey)
	if err != nil {
		cli.HandleError(err)
	}
	encryptedSecretKeyBytes, err := converter.FromB64(keyAttribute.EncryptedSecretKey)
	if err != nil {
		cli.HandleError(err)
	}
	keyEncryptionNonceBytes, err := converter.FromB64(keyAttribute.SecretKeyDecryptionNonce)
	if err != nil {
		cli.HandleError(err)
	}
	_, err = crypto.Decrypt(
		encryptedSecretKeyBytes,
		masterKeyBytes,
		keyEncryptionNonceBytes,
	)
	if err != nil {
		cli.HandleError(err)
	}

	// updating user credential
	kekEncrypted, err := crypto.GenerateKeyAndEncrypt(converter.ToB64(kekBytes))
	if err != nil {
		cli.HandleError(err)
	}
	userCredential.Email = email
	userCredential.KeyAttribute = keyAttribute
	userCredential.KekEncrypted = &kekEncrypted
}

func askForCredential() (email string, password string, err error) {
	validateEmail := func(input string) error {
		matched, err := regexp.MatchString(
			"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$",
			input,
		)
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
		if len(input) < 8 {
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
	masterKeyBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		return "", nil, nil, err
	}
	masterKey := converter.ToB64(masterKeyBytes)

	recoveryKeyBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		return "", nil, nil, err
	}
	recoveryKey := converter.ToB64(recoveryKeyBytes)

	kekSaltBytes, err := crypto.GenerateRandomBytes(blake2b.Size)
	if err != nil {
		return "", nil, nil, err
	}
	kekBytes, memLimit, opsLimit := deriveKey(password, kekSaltBytes)
	kek := converter.ToB64(kekBytes)

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
	if err != nil {
		return "", nil, nil, err
	}

	setupSRPAccountSignupReq, err := generateSRPSetupAttribute(loginSubKey)
	if err != nil {
		return "", nil, nil, err
	}

	return masterKey, &schemas.KeyAttributeInfo{
		MemLimit:                 memLimit,
		OpsLimit:                 opsLimit,
		KekSalt:                  converter.ToB64(kekSaltBytes),
		EncryptedKey:             converter.ToB64(masterKeyEncryptedWithKek.Cipher),
		KeyDecryptionNonce:       converter.ToB64(masterKeyEncryptedWithKek.Nonce),
		PublicKey:                publicKey,
		EncryptedSecretKey:       converter.ToB64(encryptedKeyPairAttributes.Cipher),
		SecretKeyDecryptionNonce: converter.ToB64(encryptedKeyPairAttributes.Nonce),
		MasterKeyEncryptedWithRecoveryKey: converter.ToB64(
			masterKeyEncryptedWithRecoveryKey.Cipher,
		),
		MasterKeyDecryptionNonce: converter.ToB64(masterKeyEncryptedWithKek.Nonce),
		RecoveryKeyEncryptedWithMasterKey: converter.ToB64(
			recoveryKeyEncryptedWithMasterKey.Cipher,
		),
		RecoveryKeyDecryptionNonce: converter.ToB64(recoveryKeyEncryptedWithMasterKey.Nonce),
	}, setupSRPAccountSignupReq, nil
}

func generateSRPSetupAttribute(loginSubKey string) (*schemas.SetupSRPAccountSignupReq, error) {
	srpParams := srp.GetParams(Srp4096Params)

	loginSubKeyBytes, err := converter.FromB64(loginSubKey)
	if err != nil {
		return nil, err
	}

	srpSaltBytes, err := crypto.GenerateRandomBytes(blake2b.Size256)
	if err != nil {
		return nil, err
	}

	srpUserID := uid.ID()

	srpVerifierBytes := srp.ComputeVerifier(
		srpParams,
		srpSaltBytes,
		[]byte(srpUserID),
		loginSubKeyBytes,
	)

	return &schemas.SetupSRPAccountSignupReq{
		SRPUserID:   srpUserID,
		SRPSalt:     converter.ToB64(srpSaltBytes),
		SRPVerifier: converter.ToB64(srpVerifierBytes),
	}, nil
}

func generateSRPClient(
	srpSalt string,
	srpUserID string,
	loginSubKey string,
) (*srp.SRPClient, error) {
	clientSecret := srp.GenKey()
	srpParams := srp.GetParams(Srp4096Params)

	loginSubKeyBytes, err := converter.FromB64(loginSubKey)
	if err != nil {
		return nil, err
	}

	srpSaltBytes, err := converter.FromB64(srpSalt)
	if err != nil {
		return nil, err
	}

	return srp.NewClient(
		srpParams,
		srpSaltBytes,
		[]byte(srpUserID),
		loginSubKeyBytes,
		clientSecret,
	), nil
}

func deriveKey(
	passphrase string,
	salt []byte,
) (kekBytes []byte, memLimit int, opsLimit int) {
	memLimit = 64 * 1024
	opsLimit = runtime.NumCPU()
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		1,
		uint32(memLimit),
		uint8(opsLimit),
		32,
	), memLimit, opsLimit
}
