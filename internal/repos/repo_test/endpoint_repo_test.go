package repo_test

import (
	"context"
	"testing"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/stretchr/testify/assert"
)

func buildEndpointEntity() *entities.Endpoint {
	return &entities.Endpoint{
		ID:                 uid.ID(),
		Name:               "Test Endpoint",
		Runtime:            "go",
		ActiveDeploymentID: "",
	}
}

func Test_endpointRepo_Add(t *testing.T) {
	endpointRepo := endpoint.NewEndpointRepo(testDB)

	testEndpointEntity := buildEndpointEntity()
	err := endpointRepo.Add(context.TODO(), testEndpointEntity)
	assert.NoError(t, err)

	err = endpointRepo.Remove(context.TODO(), testEndpointEntity.ID)
	assert.NoError(t, err)
}

func Test_endpointRepo_UpdateActiveDeployment(t *testing.T) {
	endpointRepo := endpoint.NewEndpointRepo(testDB)

	testEndpointEntity := buildEndpointEntity()
	err := endpointRepo.Add(context.TODO(), testEndpointEntity)
	assert.NoError(t, err)

	err = endpointRepo.UpdateActiveDeployment(context.TODO(), testEndpointEntity.ID, uid.ID())
	assert.NoError(t, err)

	err = endpointRepo.Remove(context.TODO(), testEndpointEntity.ID)
	assert.NoError(t, err)
}
