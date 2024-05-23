package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/supabase/hibp"
)

const DefaultMinPasswordLength = 8

func CallCheckPasswordStrength(
	httpClient *resty.Client,
	password string,
) error {
	if len(password) < DefaultMinPasswordLength {
		return fmt.Errorf("password should be at least %d characters", DefaultMinPasswordLength)
	}

	// all HIBP API requests should finish quickly to avoid
	// unnecessary slowdowns
	cc := httpClient.Clone().SetTimeout(5 * time.Second)
	hibpClient := &hibp.PwnedClient{
		UserAgent: UserAgent,
		HTTP:      cc.GetClient(),
	}

	pwned, err := hibpClient.Check(context.Background(), password)
	if err != nil {
		return fmt.Errorf(
			"unable to perform password strength check with HaveIBeenPwned.org, pwned passwords are being allowed",
		)
	}
	if pwned {
		return fmt.Errorf(
			"password is known to be weak and easy to guess, please choose a different one",
		)
	}

	return nil
}
