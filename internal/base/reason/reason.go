package reason

const (
	Success               = "base.success"
	UnknownError          = "base.unknown"
	RequestFormatError    = "base.request_format_error"
	UnauthorizedError     = "base.unauthorized_error"
	DatabaseError         = "base.database_error"
	ForbiddenError        = "base.forbidden_error"
	DuplicateRequestError = "base.duplicate_request_error"
)

const (
	EndpointNotFound           = "error.endpoint.not_found"
	EndpointHasNotPublished    = "error.endpoint.has_not_published"
	DeploymentNotFound         = "error.deployment.not_found"
	DeploymentAlreadyActivated = "error.deployment.already_activated"
	EmailDuplicate             = "error.email.duplicate"
	RequiredSession            = "error.access_token.session_required"
	SessionNotFound            = "error.session.not_found"
	UserNotFound               = "error.user.not_found"
)
