package grains

import (
	"time"

	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/rs/zerolog/log"
)

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
			log.Error().Msg("failed to request a runtime PID")
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
	}, time.Second*5).Result()
	if err != nil {
		log.Warn().Err(err).Msg("runtime manager response failed")
		return nil
	}
	pid, ok := resp.(*actor.PID)
	if !ok {
		log.Warn().Msg("runtime manager responded with a non *actor.PID")
	}
	return pid
}

func NewRequestWithResponse(req *messages.HTTPRequest) *RequestWithResponse {
	return &RequestWithResponse{
		Req:  req,
		Resp: make(chan *messages.HTTPResponse),
	}
}
