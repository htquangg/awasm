package grains

import (
	"bytes"
	"context"
	"fmt"

	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/protocluster/repos"
	"github.com/htquangg/a-wasm/internal/runtime"
	"github.com/rs/zerolog/log"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"
)

const KindRuntime = "kind"

type runtimeActor struct {
	deploymentRepo *repos.DeploymentRepo
	deploymentID   string
	stdout         *bytes.Buffer
	runtime        *runtime.Runtime
}

func NewRuntimeActor(deploymentRepo *repos.DeploymentRepo) actor.Actor {
	return &runtimeActor{
		deploymentRepo: deploymentRepo,
		stdout:         &bytes.Buffer{},
	}
}

func (r *runtimeActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.HTTPRequest:
		log.Info().
			Str("request_id", msg.ID).
			Any("pid", ctx.Self()).
			Msg("runtime handling request")

		if r.runtime == nil {
			r.initialize(msg)
		}

		r.handleHTTPRequest(ctx, msg)
	}
}

func (r *runtimeActor) initialize(msg *messages.HTTPRequest) error {
	deployment, exists, err := r.deploymentRepo.GetByID(context.Background(), msg.DeploymentID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("runtime: could not find deployment (%s)", r.deploymentID)
	}

	args := runtime.Args{
		DeploymentID: deployment.ID,
		Engine:       msg.Runtime,
		Stdout:       r.stdout,
		Data:         deployment.Data,
	}

	run, err := runtime.New(context.Background(), args)
	if err != nil {
		return err
	}
	r.runtime = run

	return nil
}

func (r *runtimeActor) handleHTTPRequest(ctx actor.Context, msg *messages.HTTPRequest) {
	b, err := proto.Marshal(msg)
	if err != nil {
		log.Warn().Err(err).Msg("failed to marshal incoming HTTP request")
		return
	}

	req := bytes.NewReader(b)
	if err := r.runtime.Invoke(req); err != nil {
		log.Warn().Err(err).Msg("runtime invoke error")
		return
	}

	resp := &messages.HTTPResponse{
		Response:   []byte{},
		RequestID:  msg.ID,
		StatusCode: int32(200),
	}

	ctx.Respond(resp)
}
