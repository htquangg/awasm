package reason

const (
	Success                   = "base.success"
	UnknownError              = "base.unknown"
	RequestFormatError        = "base.request_format_error"
	UnauthorizedError         = "base.unauthorized_error"
	DatabaseError             = "base.database_error"
	MailServerError           = "base.mailserver_error"
	ForbiddenError            = "base.forbidden_error"
	DuplicateRequestError     = "base.duplicate_request_error"
	TooManyWrongAttemptsError = "base.too_many_wrong_attempts_error"
)

const (
	EmailOrPasswordWrong        = "error.object.email_or_password_incorrect"
	InvalidTokenError           = "error.common.invalid_token"
	InvalidScopeError           = "error.common.invalid_scope"
	EndpointNotFound            = "error.endpoint.not_found"
	EndpointAccessDenied        = "error.endpoint.access_denied"
	EndpointHasNotPublished     = "error.endpoint.has_not_published"
	DeploymentNotFound          = "error.deployment.not_found"
	DeploymentAccessDenied      = "error.deployment.access_denied"
	DeploymentAlreadyActivated  = "error.deployment.already_activated"
	EmailDuplicate              = "error.email.duplicate"
	EmailNotFound               = "error.email.not_found"
	AccessTokenSessionRequired  = "error.access_token.session_required"
	SessionNotFound             = "error.session.not_found"
	UserNotFound                = "error.user.not_found"
	SRPNotFound                 = "error.srp.not_found"
	SRPAlreadyVerified          = "error.srp.already_verified"
	SRPChallengeNotFound        = "error.srp_challenge.not_found"
	SRPChallengeAlreadyVerified = "error.srp_challenge.already_verified"
	KeyAttributeNotFound        = "error.key_attribute.not_found"
	OTPExpired                  = "error.otp.expired"
	OTPIncorrect                = "error.otp.incorrect"
	ApiKeyInvalid               = "error.api_key.invalid"
	ApiKeyRequired              = "error.api_key.required"
)
