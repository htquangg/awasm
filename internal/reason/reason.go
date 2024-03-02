package reason

const (
	Success               = "success"
	UnknownError          = "unknown"
	RequestFormatError    = "request_format_error"
	UnauthorizedError     = "unauthorized_error"
	DatabaseError         = "database_error"
	ForbiddenError        = "forbidden_error"
	DuplicateRequestError = "duplicate_request_error"
)

const (
	EndpointNotFound           = "error.endpoint.not_found"
	EndpointHasNotPublished    = "error.endpoint.has_not_published"
	DeploymentNotFound         = "error.deployment.not_found"
	DeploymentAlreadyActivated = "error.deployment.already_activated"
)
