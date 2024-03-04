package repo_test

import (
	"context"
	"testing"

	"github.com/htquangg/a-wasm/internal/entities"
	"github.com/htquangg/a-wasm/internal/repos/deployment"
	"github.com/htquangg/a-wasm/internal/repos/endpoint"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/stretchr/testify/assert"
)

func buildDeploymentEntity() *entities.Deployment {
	return &entities.Deployment{
		ID:         uid.ID(),
		EndpointID: "",
		Hash:       "",
		Data:       []byte("Hello World!!!"),
	}
}

func Test_deploymentRepo_Add(t *testing.T) {
	endpointRepo := endpoint.NewEndpointRepo(testDB)
	deploymentRepo := deployment.NewDeploymentRepo(testDB)

	testEndpointEntity := buildEndpointEntity()
	err := endpointRepo.Add(context.TODO(), testEndpointEntity)
	assert.NoError(t, err)

	testDeploymentEntity := buildDeploymentEntity()
	testDeploymentEntity.EndpointID = testEndpointEntity.ID
	err = deploymentRepo.Add(context.TODO(), testDeploymentEntity)
	assert.NoError(t, err)

	err = deploymentRepo.Remove(context.TODO(), testDeploymentEntity.ID)
	assert.NoError(t, err)

	err = endpointRepo.Remove(context.TODO(), testEndpointEntity.ID)
	assert.NoError(t, err)
}
