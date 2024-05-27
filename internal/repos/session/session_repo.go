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
)

type sessionRepo struct {
	db db.DB
}

func NewSessionRepo(db db.DB) session.SessionRepo {
	return &sessionRepo{
		db: db,
	}
}

func (s *sessionRepo) WithTx(
	parentCtx context.Context,
	f func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	return s.db.WithTx(parentCtx, f)
}

func (s *sessionRepo) AddRefreshToken(
	ctx context.Context,
	refreshToken *entities.RefreshToken,
) (err error) {
	_, err = s.db.Engine(ctx).Insert(refreshToken)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (s *sessionRepo) GetRefreshTokenByToken(
	ctx context.Context,
	token string,
	forUpdate bool,
) (*entities.RefreshToken, *entities.Session, error) {
	refreshToken := &entities.RefreshToken{}

	var exists bool
	var err error

	if forUpdate {
		exists, err = s.db.Engine(ctx).
			SQL(fmt.Sprintf("SELECT * FROM %q WHERE token = $1 LIMIT 1 FOR UPDATE SKIP LOCKED", refreshToken.TableName()), token).
			Get(refreshToken)
	} else {
		exists, err = s.db.Engine(ctx).Where("token = $1", token).Get(refreshToken)
	}
	if err != nil {
		return nil, nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exists {
		return nil, nil, errors.BadRequest(reason.RefreshTokenNotFound)
	}

	var session *entities.Session

	session, exists, err = s.GetSessionByID(ctx, refreshToken.SessionID, forUpdate)
	if err != nil {
		return nil, nil, err
	}
	if !exists {
		return nil, nil, errors.BadRequest(reason.SessionNotFound)
	}

	return refreshToken, session, nil
}

func (s *sessionRepo) AddSession(ctx context.Context, session *entities.Session) (err error) {
	_, err = s.db.Engine(ctx).Insert(session)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (s *sessionRepo) AddClaimToSession(
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
	forUpdate bool,
) (*entities.Session, bool, error) {
	session := &entities.Session{}

	var exists bool
	var err error

	if forUpdate {
		exists, err = r.db.Engine(ctx).
			SQL(fmt.Sprintf("SELECT * FROM %q WHERE id = $1 LIMIT 1 FOR UPDATE SKIP LOCKED", session.TableName()), id).
			Get(session)
	} else {
		exists, err = r.db.Engine(ctx).ID(id).Get(session)
	}
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exists {
		return nil, false, errors.BadRequest(reason.SessionNotFound)
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

func (r *sessionRepo) GetCurrentlyActiveRefreshTokenBySessionID(
	ctx context.Context,
	sessionID string,
) (*entities.RefreshToken, bool, error) {
	refreshToken := &entities.RefreshToken{}

	exists, err := r.db.Engine(ctx).
		Where("session_id = $1 AND revoked IS false", sessionID).
		Get(refreshToken)

	return refreshToken, exists, err
}

func (r *sessionRepo) UpdateRefreshTokenCols(
	ctx context.Context,
	refreshToken *entities.RefreshToken,
	cols ...string,
) error {
	_, err := r.db.Engine(ctx).ID(refreshToken.ID).Cols(cols...).Update(refreshToken)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
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
