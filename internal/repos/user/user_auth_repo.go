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

func (r *userAuthRepo) GetSRPAuthWithSRPUserID(
	ctx context.Context,
	srpUserID string,
) (srpAuth *entities.SrpAuth, exists bool, err error) {
	srpAuth = &entities.SrpAuth{}
	exists, err = r.db.Engine(ctx).Where("srp_user_id = $1", srpUserID).Get(srpAuth)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return srpAuth, exists, err
}

func (r *userAuthRepo) GetSRPAttribute(ctx context.Context, userID string) (*schemas.GetSRPAttributeResp, bool, error) {
	respFromDB := make([]*struct {
		SRPUserID string `json:"srpUserId"`
		Salt      string `json:"salt"`
		KekSalt   string `json:"kekSalt"`
		MemLimit  int    `json:"memLimit"`
		OpsLimit  int    `json:"opsLimit"`
	}, 0, 1)
	err := r.db.Engine(ctx).
		Join("INNER", "key_attributes", "`key_attributes`.user_id =`srp_auth`.user_id").
		Select("srp_user_id, salt, mem_limit, ops_limit, kek_salt").
		Table("srp_auth").
		Where("`key_attributes`.user_id = $1", userID).
		Limit(1).
		Find(&respFromDB)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	// no records
	if len(respFromDB) == 0 {
		return nil, false, nil
	}

	getSRPAttributeResp := &schemas.GetSRPAttributeResp{}
	_ = copier.Copy(getSRPAttributeResp, &respFromDB[0])
	getSRPAttributeResp.SRPSalt = respFromDB[0].Salt

	return getSRPAttributeResp, true, nil
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
			SrpUserID: accountInfo.SRPUserID,
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

		now := time.Now()
		_, errTx = r.db.Engine(ctx).ID(accountInfo.UserID).Update(&entities.User{
			EmailAcceptedAt: &now,
		})
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
