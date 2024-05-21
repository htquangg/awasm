package protocluster

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/htquangg/a-wasm/internal/base/db"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/protocluster/grains"
	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/protocluster/repos"
)

type Cluster struct {
	ctx           context.Context
	c             *cluster.Cluster
	wasmServerPID *actor.PID
	repos         *repos.Repos
}

func New(ctx context.Context, db db.DB) *Cluster {
	system := actor.NewActorSystem()

	repos := repos.New(db)

	kinds := initKinds(repos)

	provider := automanaged.New()
	config := remote.Configure("127.0.0.1", 0)
	lookup := disthash.New()
	clusterConfig := cluster.Configure(
		"awasm-cluster",
		provider,
		lookup,
		config,
		cluster.WithKinds(kinds...),
	)

	cluster := &Cluster{
		ctx:   ctx,
		c:     cluster.New(system, clusterConfig),
		repos: repos,
	}

	return cluster
}

func (c *Cluster) Start() error {
	c.c.StartMember()
	c.createWasmServerActor()
	return nil
}

func (c *Cluster) Shutdown() error {
	_, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	c.c.Shutdown(true)

	return nil
}

func initKinds(repos *repos.Repos) []*cluster.Kind {
	kinds := make([]*cluster.Kind, 0)

	runtimeManagerKind := getRuntimeManagerKind()
	runtimeKind := getRuntimeKind(repos.Deployment)

	kinds = append(kinds,
		runtimeManagerKind,
		runtimeKind,
	)

	return kinds
}

func getRuntimeManagerKind() *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return grains.NewRuntimeManagerActor()
	})

	kind := cluster.NewKind(grains.KindRuntimeManager, props)

	return kind
}

func getRuntimeKind(deploymentRepo *repos.DeploymentRepo) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return grains.NewRuntimeActor(deploymentRepo)
	})

	kind := cluster.NewKind(grains.KindRuntime, props)

	return kind
}

func (c *Cluster) createWasmServerActor() {
	props := actor.PropsFromProducer(func() actor.Actor {
		return grains.NewWasmServerActor(c.c)
	})

	c.wasmServerPID = c.c.ActorSystem.Root.Spawn(props)
}

func (c *Cluster) Serve(req *messages.HTTPRequest) *messages.HTTPResponse {
	reqResp := grains.NewRequestWithResponse(req)
	c.c.ActorSystem.Root.Send(c.wasmServerPID, reqResp)

	resp := <-reqResp.Resp

	return resp
}
