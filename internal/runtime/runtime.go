package runtime

import (
	"context"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type (
	Args struct {
		DeploymentID string
		Engine       string
		Stdout       io.Writer
		Data         []byte
		Cache        wazero.CompilationCache
	}

	Runtime struct {
		runtime wazero.Runtime
		ctx     context.Context
		stdout  io.Writer
		engine  string
		data    []byte
		mod     wazero.CompiledModule
	}
)

func New(ctx context.Context, args Args) (*Runtime, error) {
	cfg := wazero.NewRuntimeConfigCompiler().WithCompilationCache(args.Cache)

	r := &Runtime{
		ctx:     ctx,
		stdout:  args.Stdout,
		runtime: wazero.NewRuntimeWithConfig(ctx, cfg),
		engine:  args.Engine,
		data:    args.Data,
	}
	wasi_snapshot_preview1.MustInstantiate(ctx, r.runtime)

	mod, err := r.runtime.CompileModule(ctx, r.data)
	if err != nil {
		return nil, err
	}
	r.mod = mod

	return r, nil
}

func (r *Runtime) Invoke(stdin io.Reader) error {
	modCfg := wazero.NewModuleConfig().
		WithStdin(stdin).
		WithStdout(r.stdout).
		WithStderr(os.Stderr)

	_, err := r.runtime.InstantiateModule(r.ctx, r.mod, modCfg)

	return err
}
