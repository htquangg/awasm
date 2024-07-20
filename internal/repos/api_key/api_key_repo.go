package api_key

import (
	"context"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/internal/base/db"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/entities"
	"github.com/htquangg/awasm/internal/services/api_key"
)

type apiKeyRepo struct {
	db db.DB
}

func NewApiKeyRepo(db db.DB) api_key.ApiKeyRepo {
	return &apiKeyRepo{
		db: db,
	}
}

func (r *apiKeyRepo) AddApiKey(ctx context.Context, apiKey *entities.ApiKey) error {
	_, err := r.db.Engine(ctx).Insert(apiKey)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *apiKeyRepo) GetApiKeyByKey(
	ctx context.Context,
	key string,
) (apiKey *entities.ApiKey, exists bool, err error) {
	apiKey = &entities.ApiKey{}
	exists, err = r.db.Engine(ctx).Where("key = $1", key).Get(apiKey)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}
