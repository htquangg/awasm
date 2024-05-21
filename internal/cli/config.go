package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/logger"
)

func WriteInitalConfig(userCredentials *schemas.UserCredential) error {
	fullConfigFilePath, fullConfigFileDirPath, err := GetFullConfigFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(fullConfigFileDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fullConfigFileDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	existingConfigFile, err := GetConfigFile()
	if err != nil {
		return fmt.Errorf("writeInitalConfig: unable to write config file because [err=%s]", err)
	}

	loggedInUser := &schemas.LoggedInUser{
		Email:  userCredentials.Email,
		Domain: config.AWASM_URL,
	}
	if existingConfigFile == nil {
		existingConfigFile = &schemas.ConfigFile{}
		existingConfigFile.LoggedInUsers = append(existingConfigFile.LoggedInUsers, loggedInUser)
	} else {
		if ConfigContainsEmail(existingConfigFile.LoggedInUsers, userCredentials.Email) {
			for idx, user := range existingConfigFile.LoggedInUsers {
				if user.Email == userCredentials.Email {
					existingConfigFile.LoggedInUsers[idx] = loggedInUser
				}
			}
		} else {
			existingConfigFile.LoggedInUsers = append(existingConfigFile.LoggedInUsers, loggedInUser)
		}
	}

	configFile := &schemas.ConfigFile{
		LoggedInUserEmail:  userCredentials.Email,
		LoggedInUserDomain: config.AWASM_URL,
		LoggedInUsers:      existingConfigFile.LoggedInUsers,
	}

	configFileMarshalled, err := json.Marshal(configFile)
	if err != nil {
		return err
	}

	err = WriteToFile(fullConfigFilePath, configFileMarshalled, 0o600)
	if err != nil {
		return err
	}

	return err
}

func ConfigFileExists() bool {
	fullConfigFileURI, _, err := GetFullConfigFilePath()
	if err != nil {
		logger.Debugf("There was an error when creating the full path to config file: %v", err)
		return false
	}

	if _, err := os.Stat(fullConfigFileURI); err == nil {
		return true
	} else {
		return false
	}
}

func GetFullConfigFilePath() (fullPathToFile string, fullPathToDirectory string, err error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", "", err
	}

	fullPath := fmt.Sprintf(
		"%s/%s/%s",
		homeDir,
		constants.CONFIG_FOLDER_NAME,
		constants.CONFIG_FILE_NAME,
	)
	fullDirPath := fmt.Sprintf("%s/%s", homeDir, constants.CONFIG_FOLDER_NAME)
	return fullPath, fullDirPath, err
}

func GetConfigFile() (*schemas.ConfigFile, error) {
	fullConfigFilePath, _, err := GetFullConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFileAsBytes, err := os.ReadFile(fullConfigFilePath)
	if err != nil {
		return nil, err
	}

	var configFile schemas.ConfigFile
	err = json.Unmarshal(configFileAsBytes, &configFile)
	if err != nil {
		return nil, err
	}

	return &configFile, nil
}

func WriteConfigFile(configFile *schemas.ConfigFile) error {
	fullConfigFilePath, fullConfigFileDirPath, err := GetFullConfigFilePath()
	if err != nil {
		return fmt.Errorf(
			"writeConfigFile: unable to write config file because an error occurred when getting config file path [err=%s]",
			err,
		)
	}

	configFileMarshalled, err := json.Marshal(configFile)
	if err != nil {
		return fmt.Errorf(
			"writeConfigFile: unable to write config file because an error occurred when marshaling the config file [err=%s]",
			err,
		)
	}

	if _, err := os.Stat(fullConfigFileDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fullConfigFileDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fullConfigFilePath, configFileMarshalled, 0o600)
	if err != nil {
		return fmt.Errorf(
			"writeConfigFile: unable to write config file because an error occurred when write the config to file [err=%s]",
			err,
		)
	}

	return nil
}
