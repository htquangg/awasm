package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/cli/api"
	"github.com/htquangg/a-wasm/internal/schemas"

	"github.com/go-resty/resty/v2"
	"github.com/zalando/go-keyring"
)

type LoggedInUserDetails struct {
	UserCredentials *schemas.UserCredential
}

func StoreUserCredsInKeyRing(userCred *schemas.UserCredential) error {
	userCredMarshalled, err := json.Marshal(userCred)
	if err != nil {
		return fmt.Errorf("StoreUserCredsInKeyRing: something went wrong when marshalling user creds [err=%s]", err)
	}

	err = SetValueInKeyring(userCred.Email, string(userCredMarshalled))
	if err != nil {
		return fmt.Errorf("StoreUserCredsInKeyRing: unable to store user credentials because [err=%s]", err)
	}

	return err
}

func GetUserCredsFromKeyRing(userEmail string) (credentials *schemas.UserCredential, err error) {
	credentialsValue, err := GetValueInKeyring(userEmail)
	if err != nil {
		if err == keyring.ErrUnsupportedPlatform {
			return nil, errors.New("your OS does not support keyring. Consider using a service token https://infisical.com/docs/documentation/platform/token")
		} else if err == keyring.ErrNotFound {
			return nil, errors.New("credentials not found in system keyring")
		} else {
			return nil, fmt.Errorf("something went wrong, failed to retrieve value from system keyring [error=%v]", err)
		}
	}

	var userCredentials schemas.UserCredential

	err = json.Unmarshal([]byte(credentialsValue), &userCredentials)
	if err != nil {
		return nil, fmt.Errorf("getUserCredsFromKeyRing: Something went wrong when unmarshalling user creds [err=%s]", err)
	}

	return &userCredentials, err
}

func GetCurrentLoggedInUserDetails() (*LoggedInUserDetails, error) {
	if ConfigFileExists() {
		configFile, err := GetConfigFile()
		if err != nil {
			return nil, fmt.Errorf("getCurrentLoggedInUserDetails: unable to get logged in user from config file [err=%s]", err)
		}
		if configFile.LoggedInUserEmail == "" {
			return nil, nil
		}

		userCreds, err := GetUserCredsFromKeyRing(configFile.LoggedInUserEmail)
		if err != nil {
			if strings.Contains(err.Error(), "credentials not found in system keyring") {
				return nil, errors.New("we couldn't find your logged in details, try running [infisical login] then try again")
			} else {
				return nil, fmt.Errorf("failed to fetch creditnals from keyring because [err=%s]", err)
			}
		}

		httpClient := resty.New().
			SetAuthToken(userCreds.AccessToken).
			SetHeader("Accept", "application/json")

		config.AWASM_URL_MANUAL_OVERRIDE = config.AWASM_URL

		if configFile.LoggedInUserDomain != "" {
			config.AWASM_URL = configFile.LoggedInUserDomain
		}

		isAuthenticated := api.CallIsAuthenticated(httpClient)

		// TODO: add refresh token

		if !isAuthenticated {
			return nil, nil
		}

		return &LoggedInUserDetails{
			UserCredentials: userCreds,
		}, nil
	} else {
		return nil, nil
	}
}
