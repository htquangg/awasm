package auth

import (
	"context"
	"encoding/json"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/cache"
	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/constants"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/auth"
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

func (r *authRepo) GetUserCacheInfo(
	ctx context.Context,
	userID string,
) (userInfo *entities.UserCacheInfo, err error) {
	cacheKey := &cache.Key{
		Namespace: constants.UserInfoNameSpaceCache,
		Key:       userID,
	}
	userInfoCache, exist, err := r.cacher.Get(ctx, cacheKey)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil, nil
	}
	userInfo = &entities.UserCacheInfo{}
	_ = json.Unmarshal(userInfoCache, userInfo)
	return userInfo, nil
}

func (r *authRepo) SetUserCacheInfo(
	ctx context.Context,
	userID string,
	userInfo *entities.UserCacheInfo,
) error {
	cachePayload, err := json.Marshal(userInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	cacheKey := &cache.Key{
		Namespace: constants.UserInfoNameSpaceCache,
		Key:       userID,
	}
	err = r.cacher.Set(ctx, cacheKey, cachePayload, constants.UserInfoTTLCache)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *authRepo) RemoveUserCacheInfo(ctx context.Context, userID string) (err error) {
	cacheKey := &cache.Key{
		Namespace: constants.UserInfoNameSpaceCache,
		Key:       userID,
	}
	err = r.cacher.Delete(ctx, cacheKey)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
