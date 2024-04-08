package user

import (
	"context"
	"time"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/internal/services/user"
	"github.com/htquangg/a-wasm/pkg/uid"
	"github.com/jinzhu/copier"

	"github.com/segmentfault/pacman/errors"
)

type userAuthRepo struct {
	cfg *config.Config
	db  db.DB
}

func NewUserAuthRepo(
	cfg *config.Config,
	db db.DB,
) user.UserAuthRepo {
	return &userAuthRepo{
		cfg: cfg,
		db:  db,
	}
}

func (r *userAuthRepo) AddSRPChallenge(ctx context.Context, srpChallege *entities.SrpChallenge) (err error) {
	_, err = r.db.Engine(ctx).Insert(srpChallege)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userAuthRepo) AddSrpAuthTemp(ctx context.Context, srpAuthTemp *entities.SrpAuthTemp) (err error) {
	_, err = r.db.Engine(ctx).Insert(srpAuthTemp)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userAuthRepo) CompleteEmailAccountSignup(
	ctx context.Context,
	accountInfo *schemas.CompleteEmailSignupInfo,
) (err error) {
	var isSrpSetupDone bool
	isSrpSetupDone, err = r.isSRPSetupDone(ctx, accountInfo.UserID)
	if err != nil {
		return err
	}
	if isSrpSetupDone {
		return errors.BadRequest(reason.SRPAlreadyVerified)
	}

	keyAttribute := &entities.KeyAttribute{}
	err = copier.Copy(keyAttribute, accountInfo.KeyAttribute)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	keyAttribute.UserID = accountInfo.UserID

	_, err = r.db.WithTx(ctx, func(ctx context.Context) (interface{}, error) {
		srpAuth := &entities.SrpAuth{
			ID:        uid.ID(),
			UserID:    accountInfo.UserID,
			SrpUserID: accountInfo.SrpUserID,
			Salt:      accountInfo.Salt,
			Verifier:  accountInfo.Verifier,
		}
		_, errTx := r.db.Engine(ctx).Insert(srpAuth)
		if errTx != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		_, errTx = r.db.Engine(ctx).Insert(keyAttribute)
		if errTx != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		return nil, nil
	})

	return err
}

func (r *userAuthRepo) GetKeyAttributeWithUserID(
	ctx context.Context,
	userID string,
) (keyAttribute *entities.KeyAttribute, exists bool, err error) {
	keyAttribute = &entities.KeyAttribute{}
	exists, err = r.db.Engine(ctx).Where("user_id = $1", userID).Get(keyAttribute)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return keyAttribute, exists, err
}

func (r *userAuthRepo) GetSRPChallengeByID(
	ctx context.Context,
	id string,
) (srpChallenge *entities.SrpChallenge, exists bool, err error) {
	srpChallenge = &entities.SrpChallenge{}
	exists, err = r.db.Engine(ctx).Where("id = $1", id).Get(srpChallenge)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return srpChallenge, exists, err
}

func (r *userAuthRepo) GetTempSRPSetupByID(
	ctx context.Context,
	id string,
) (srpAuthTemp *entities.SrpAuthTemp, exists bool, err error) {
	srpAuthTemp = &entities.SrpAuthTemp{}
	exists, err = r.db.Engine(ctx).Where("id = $1", id).Get(srpAuthTemp)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return srpAuthTemp, exists, err
}

func (r *userAuthRepo) IncrementSrpChallengeAttemptCount(ctx context.Context, id string) (err error) {
	_, err = r.db.Engine(ctx).Incr("attempt_count").ID(id).Update(new(entities.SrpChallenge))
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return err
}

func (r *userAuthRepo) SetSrpChallengeVerified(ctx context.Context, challengeID string) (err error) {
	now := time.Now()
	srpChallenge := &entities.SrpChallenge{
		VerifiedAt: &now,
	}
	_, err = r.db.Engine(ctx).Where("id = ?", challengeID).Cols("verified_at").Update(srpChallenge)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *userAuthRepo) isSRPSetupDone(ctx context.Context, userID string) (exists bool, err error) {
	srpAuth := &entities.SrpAuth{}
	exists, err = r.db.Engine(ctx).Where("user_id = $1", userID).Get(srpAuth)
	if err != nil {
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return exists, nil
}
