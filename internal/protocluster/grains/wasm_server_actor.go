package grains

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"

	"github.com/htquangg/awasm/internal/protocluster/grains/messages"
	"github.com/htquangg/awasm/pkg/logger"
)

const RequestRuntimeTimeout = 5 * time.Second

type (
	wasmServerActor struct {
		c                 *cluster.Cluster
		responses         map[string]chan *messages.HTTPResponse
		self              *actor.PID
		runtimeManagerPID *actor.PID
	}

	RequestWithResponse struct {
		Req  *messages.HTTPRequest
		Resp chan *messages.HTTPResponse
	}
)

func NewWasmServerActor(c *cluster.Cluster) actor.Actor {
	return &wasmServerActor{
		c:                 c,
		responses:         make(map[string]chan *messages.HTTPResponse),
		runtimeManagerPID: cluster.GetCluster(c.ActorSystem).Get("1", KindRuntimeManager),
	}
}

func (s *wasmServerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		s.initialize(ctx)
	case *actor.Stopped:
	case *RequestWithResponse:
		pid := s.requestRuntime(ctx, msg.Req.GetDeploymentID())
		if pid == nil {
			logger.Error("failed to request a runtime PID")
			return
		}
		s.responses[msg.Req.GetID()] = msg.Resp
		ctx.Request(pid, msg.Req)
	case *messages.HTTPResponse:
		if resp, ok := s.responses[msg.RequestID]; ok {
			resp <- msg
			delete(s.responses, msg.RequestID)
		}
	}
}

func (s *wasmServerActor) initialize(ctx actor.Context) {
	s.self = ctx.Self()
}

func (s *wasmServerActor) requestRuntime(_ actor.Context, key string) *actor.PID {
	resp, err := s.c.ActorSystem.Root.RequestFuture(s.runtimeManagerPID, &requestRuntime{
		key: key,
	}, RequestRuntimeTimeout).Result()
	if err != nil {
		logger.Warn("runtime manager response failed")
		return nil
	}
	pid, ok := resp.(*actor.PID)
	if !ok {
		logger.Warn("runtime manager responded with a non *actor.PID")
	}
	return pid
}

func NewRequestWithResponse(req *messages.HTTPRequest) *RequestWithResponse {
	return &RequestWithResponse{
		Req:  req,
		Resp: make(chan *messages.HTTPResponse),
	}
}
