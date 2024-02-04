package cluster

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/automanaged"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

type Cluster struct {
	ctx     context.Context
	cluster *cluster.Cluster
}

func New(ctx context.Context) *Cluster {
	system := actor.NewActorSystem()

	provider := automanaged.New()
	config := remote.Configure("0.0.0.0", 0)
	lookup := disthash.New()

	clusterConfig := cluster.Configure("awasm-cluster", provider, lookup, config)
	cluster := cluster.New(system, clusterConfig)

	return &Cluster{
		ctx:     ctx,
		cluster: cluster,
	}
}

func (c *Cluster) ServeHandler() (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(c.ctx)
	return func() error {
			c.cluster.StartMember()
			<-ctx.Done()
			return ctx.Err()
		}, func(err error) {
			defer cancel()
			c.cluster.Shutdown(true)
		}
}
