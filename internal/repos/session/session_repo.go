package session

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/session"
	"github.com/htquangg/a-wasm/pkg/crypto"
	"github.com/htquangg/a-wasm/pkg/uid"
)

type sessionRepo struct {
	db db.DB
}

func NewSessionRepo(db db.DB) session.SessionRepo {
	return &sessionRepo{
		db: db,
	}
}

func (s *sessionRepo) CreateRefreshToken(
	ctx context.Context,
	userID string,
	authenticationMethod entities.AuthenticationMethod,
	params *entities.GrantParams,
) (*entities.RefreshToken, error) {
	token, err := crypto.GenerateURLSafeRandomString(entities.TokenLength)
	if err != nil {
		return nil, err
	}
	refreshToken := &entities.RefreshToken{
		ID:     uid.ID(),
		UserID: userID,
		Token:  token,
	}
	defaultAAL := entities.Aal1.String()
	session := &entities.Session{
		ID:        uid.ID(),
		UserID:    userID,
		AAL:       defaultAAL,
		IP:        params.IP,
		UserAgent: params.UserAgent,
	}

	if _, err = s.db.WithTx(ctx, func(ctxTx context.Context) (interface{}, error) {
		// 1. Create session
		if _, terr := s.db.Engine(ctxTx).Insert(session); terr != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		// 2. Update session into refresh token
		refreshToken.SessionID = session.ID
		if _, terr := s.db.Engine(ctxTx).Insert(refreshToken); terr != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		// 3. Update AMR Claims
		if terr := s.addClaimToSession(ctxTx, session.ID, authenticationMethod); terr != nil {
			return nil, terr
		}

		return nil, nil
	}); err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *sessionRepo) addClaimToSession(
	ctx context.Context,
	sessionID string,
	authenticationMethod entities.AuthenticationMethod,
) error {
	amrClaim := &entities.MFAAMRClaim{
		SessionID:            sessionID,
		AuthenticationMethod: authenticationMethod.String(),
	}

	if _, err := s.db.Engine(ctx).Exec(
		fmt.Sprintf(`
		INSERT INTO %s (session_id, authentication_method) VALUES ($1, $2)
			ON CONFLICT ON CONSTRAINT pk_mfa_amr_claims_session_id_authentication_method
			DO UPDATE SET updated_at = $3;
	`, amrClaim.TableName()), amrClaim.SessionID, amrClaim.AuthenticationMethod, time.Now(),
	); err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *sessionRepo) GetSessionByID(
	ctx context.Context,
	id string,
) (*entities.Session, bool, error) {
	session := &entities.Session{}
	exists, err := r.db.Engine(ctx).ID(id).Get(session)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	if session.AMRClaims == nil {
		amrClaims, err := r.getAMRClaimsWithSessionID(ctx, session.ID)
		if err != nil {
			return nil, false,
				errors.InternalServer(reason.DatabaseError).
					WithError(err).
					WithStack()
		}
		session.AMRClaims = amrClaims
	}

	return session, exists, nil
}

func (r *sessionRepo) getAMRClaimsWithSessionID(
	ctx context.Context,
	sessionID string,
) ([]*entities.MFAAMRClaim, error) {
	var amrClaims []*entities.MFAAMRClaim

	if err := r.db.Engine(ctx).Where("session_id = $1", sessionID).Find(&amrClaims); err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return amrClaims, nil
}
