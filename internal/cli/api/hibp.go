package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/supabase/hibp"
)

const DEFAULT_MIN_PASSWORD_LENGTH = 8

func CallCheckPasswordStrength(
	httpClient *resty.Client,
	password string,
) error {
	if len(password) < DEFAULT_MIN_PASSWORD_LENGTH {
		return fmt.Errorf("Password should be at least %d characters.", DEFAULT_MIN_PASSWORD_LENGTH)
	}

	// all HIBP API requests should finish quickly to avoid
	// unnecessary slowdowns
	cc := httpClient.Clone().SetTimeout(5 * time.Second)
	hibpClient := &hibp.PwnedClient{
		UserAgent: USER_AGENT,
		HTTP:      cc.GetClient(),
	}

	pwned, err := hibpClient.Check(context.Background(), password)
	if err != nil {
		return fmt.Errorf(
			"Unable to perform password strength check with HaveIBeenPwned.org, pwned passwords are being allowed",
		)
	}
	if pwned {
		return fmt.Errorf("Password is known to be weak and easy to guess, please choose a different one.")
	}

	return nil
}
