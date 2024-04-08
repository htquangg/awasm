package user

import (
	"context"

	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/kong/go-srp"
	"github.com/segmentfault/pacman/errors"
)

const (
	Srp4096Params                = 4096
	MaxAttempsVerifySRPChallenge = 5
)

func (s *UserService) createAndInsertSRPChallenge(
	ctx context.Context,
	srpUserID, srpVerifier, srpA string,
) (*string, string, error) {
	serverSecret := srp.GenKey()
	srpParams := srp.GetParams(Srp4096Params)

	srpVerifierBytes, err := convertStringToBytes(srpVerifier)
	if err != nil {
		return nil, "", err
	}
	srpServer := srp.NewServer(srpParams, srpVerifierBytes, serverSecret)
	if srpServer == nil {
		return nil, "", errors.InternalServer(reason.UnknownError).WithStack()
	}

	srpB := srpServer.ComputeB()
	if srpB == nil {
		return nil, "", errors.InternalServer(reason.UnknownError).WithStack()
	}

	srpChallenge := &entities.SrpChallenge{}
	srpChallenge.ID = uid.ID()
	srpChallenge.SrpUserID = srpUserID
	srpChallenge.ServerKey = convertBytesToString(serverSecret)
	srpChallenge.SrpA = srpA
	err = s.userAuthRepo.AddSRPChallenge(ctx, srpChallenge)
	if err != nil {
		return nil, "", err
	}

	srpBBase64 := convertBytesToString(srpB)
	return &srpBBase64, srpChallenge.ID, nil
}

func (s *UserService) verifySRPChallenge(
	ctx context.Context,
	srpVerifier string,
	challengeID string,
	srpM1 string,
) (*string, error) {
	srpChallenge, exists, err := s.userAuthRepo.GetSRPChallengeByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.BadRequest(reason.SRPChallengeNotFound)
	}
	if srpChallenge.VerifiedAt != nil {
		return nil, errors.BadRequest(reason.SRPChallengeAlreadyVerified)
	} else if srpChallenge.AttemptCount >= MaxAttempsVerifySRPChallenge {
		return nil, errors.BadRequest(reason.TooManyWrongAttemptsError)
	}

	srpParams := srp.GetParams(Srp4096Params)
	srpVerifierBytes, err := convertStringToBytes(srpVerifier)
	if err != nil {
		return nil, err
	}
	srpServerKeyBytes, err := convertStringToBytes(srpChallenge.ServerKey)
	if err != nil {
		return nil, err
	}
	srpServer := srp.NewServer(srpParams, srpVerifierBytes, srpServerKeyBytes)
	if srpServer == nil {
		return nil, errors.InternalServer(reason.UnknownError).WithMsg("server is nil.").WithStack()
	}

	srpABytes, err := convertStringToBytes(srpChallenge.SrpA)
	if err != nil {
		return nil, err
	}
	srpServer.SetA(srpABytes)

	srpM1Bytes, err := convertStringToBytes(srpM1)
	if err != nil {
		return nil, err
	}
	srpM2Bytes, err := srpServer.CheckM1(srpM1Bytes)
	if err != nil {
		err2 := s.userAuthRepo.IncrementSrpChallengeAttemptCount(ctx, challengeID)
		if err2 != nil {
			return nil, err2
		}
		return nil, errors.BadRequest(reason.EmailOrPasswordWrong)
	} else {
		err2 := s.userAuthRepo.SetSrpChallengeVerified(ctx, challengeID)
		if err2 != nil {
			return nil, err2
		}
	}
	srpM2 := convertBytesToString(srpM2Bytes)

	return &srpM2, nil
}
