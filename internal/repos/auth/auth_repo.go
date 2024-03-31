package auth

import (
	"context"
	"encoding/json"

	"github.com/htquangg/a-wasm/internal/base/cache"
	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/services/auth"

	"github.com/segmentfault/pacman/errors"
)

type authRepo struct {
	db     db.DB
	cacher cache.Cacher
}

func NewAuthRepo(db db.DB, cacher cache.Cacher) auth.AuthRepo {
	return &authRepo{
		db:     db,
		cacher: cacher,
	}
}

func (r *authRepo) GetUserCacheInfo(ctx context.Context, userID string) (userInfo *entities.UserCacheInfo, err error) {
	userInfoCache, exist, err := r.cacher.Get(ctx, constants.UserInfoCacheKey+userID)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil, nil
	}
	userInfo = &entities.UserCacheInfo{}
	_ = json.Unmarshal([]byte(userInfoCache), userInfo)
	return userInfo, nil
}

func (r *authRepo) SetUserCacheInfo(ctx context.Context, userID string, userInfo *entities.UserCacheInfo) error {
	cachePayload, err := json.Marshal(userInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	err = r.cacher.Set(ctx, constants.UserInfoCacheKey+userID, cachePayload, constants.UserInfoCacheTime)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *authRepo) RemoveUserCacheInfo(ctx context.Context, userID string) (err error) {
	err = r.cacher.Delete(ctx, constants.UserInfoCacheKey+userID)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
