package cli

import "github.com/htquangg/awasm/internal/schemas"

func CheckAuthentication() {
	configFile, _ := GetConfigFile()
	if configFile == nil {
		PrintErrorMessageAndExit(
			"You must be logged in to run this command. To login, run [awasm login]",
		)
	}

	if configFile.LoggedInUserEmail == "" {
		PrintErrorMessageAndExit(
			"You must be logged in to run this command. To login, run [awasm login]",
		)
	}
}

func ConfigContainsEmail(users []*schemas.LoggedInUser, email string) bool {
	for _, value := range users {
		if value.Email == email {
			return true
		}
	}
	return false
}
