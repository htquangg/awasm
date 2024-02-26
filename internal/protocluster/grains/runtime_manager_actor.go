package grains

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
)

const KindRuntimeManager = "runtime_manager"

type (
	runtimeManagerActor struct {
		runtimes map[string]*actor.PID
	}

	requestRuntime struct {
		key string
	}
)

func NewRuntimeManagerActor() actor.Actor {
	return &runtimeManagerActor{
		runtimes: make(map[string]*actor.PID),
	}
}

func (rm *runtimeManagerActor) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *requestRuntime:
		pid := rm.runtimes[msg.key]
		if pid == nil {
			pid = cluster.GetCluster(c.ActorSystem()).Get(msg.key, KindRuntime)
			rm.runtimes[msg.key] = pid
		}
		c.Respond(pid)
	case actor.Started:
	case actor.Stopped:
	}
}
