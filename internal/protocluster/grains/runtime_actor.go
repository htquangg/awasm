package grains

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"google.golang.org/protobuf/proto"

	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/internal/protocluster/repos"
	"github.com/htquangg/a-wasm/internal/runtime"
	"github.com/htquangg/a-wasm/pkg/logger"
)

const (
	KindRuntime = "kind"

	magicLen = 1 << 8
)

type runtimeActor struct {
	deploymentRepo *repos.DeploymentRepo
	stdout         *bytes.Buffer
	runtime        *runtime.Runtime
	deploymentID   string
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
		logger.Infof("runtime handling request with request_id %s, pid %v", msg.ID, ctx.Self())

		if r.runtime == nil {
			_ = r.initialize(msg)
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
		logger.Warnf("failed to marshal incoming HTTP request: %v", err)
		return
	}

	req := bytes.NewReader(b)
	if err := r.runtime.Invoke(req); err != nil {
		logger.Warnf("failed to invoke runtime: %v", err)
		handleResponse(ctx, http.StatusBadRequest, []byte("failed to invoke runtime"), nil, msg.ID)
		return
	}
	_, res, header, status, err := ParseStdout(r.stdout)
	if err != nil {
		handleResponse(ctx, http.StatusOK, []byte(err.Error()), nil, msg.ID)
		return
	}

	handleResponse(ctx, int32(status), res, header, msg.ID)
}

func ParseStdout(stdout io.Reader) (logs []byte, resp []byte, header []byte, status int, err error) {
	stdoutb, err := io.ReadAll(stdout)
	if err != nil {
		return
	}

	outLen := len(stdoutb)
	if outLen < magicLen {
		err = fmt.Errorf("mallformed HTTP response missing last %d bytes", magicLen)
		return
	}

	wasmRespStart := outLen - magicLen

	status = int(binary.LittleEndian.Uint32(stdoutb[wasmRespStart : wasmRespStart+4]))
	respLen := binary.LittleEndian.Uint32(stdoutb[wasmRespStart+4 : wasmRespStart+8])
	headerLen := binary.LittleEndian.Uint32(stdoutb[wasmRespStart+8 : wasmRespStart+12])
	if int(respLen)+int(headerLen) > outLen-magicLen {
		err = fmt.Errorf("response length exceeds available data")
		return
	}

	respStart := outLen - magicLen - int(headerLen) - int(respLen)
	resp = stdoutb[respStart : respStart+int(respLen)]

	headerStart := outLen - magicLen - int(headerLen)
	header = stdoutb[headerStart : headerStart+int(headerLen)]

	logs = stdoutb[:respStart]

	return
}

func handleResponse(ctx actor.Context, code int32, msg []byte, header []byte, id string) {
	ctx.Respond(&messages.HTTPResponse{
		Response:   msg,
		Header:     header,
		StatusCode: code,
		RequestID:  id,
	})
}
