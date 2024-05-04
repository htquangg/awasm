package constants

import "time"

const (
	UserInfoNameSpaceCache         = "users:by_id:"
	MailerOTPNameSpaceCache        = "mailer:otp"
	AuthorizedApiKeyNameSpaceCache = "authorized:by_api_key"
	UserInfoTTLCache               = 7 * 24 * time.Hour
	AuthorizedApiKeyTTLCache       = 5 * time.Minute
)
