package handler

import (
	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/reason"
)

func IsNotfoundError(err error) bool {
	perr, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	if errors.IsNotFound(perr) {
		return true
	}

	switch perr.Reason {
	case reason.EndpointNotFound:
		return true
	case reason.DeploymentNotFound:
		return true
	case reason.EmailNotFound:
		return true
	case reason.SessionNotFound:
		return true
	case reason.UserNotFound:
		return true
	case reason.SRPNotFound:
		return true
	case reason.SRPChallengeNotFound:
		return true
	case reason.KeyAttributeNotFound:
		return true
	case reason.RefreshTokenNotFound:
		return true
	}

	return false
}
