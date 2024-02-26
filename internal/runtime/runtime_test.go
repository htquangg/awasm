package runtime

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"
	"github.com/htquangg/a-wasm/pkg/uid"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/wazero"
	"google.golang.org/protobuf/proto"
)

func TestRuntimeInvoke(t *testing.T) {
	b, err := os.ReadFile("../_testdata/helloworld.wasm")
	require.Nil(t, err)

	req := &messages.HTTPRequest{
		Method: "get",
		URL:    "/",
		Body:   nil,
	}
	breq, err := proto.Marshal(req)
	require.Nil(t, err)

	out := &bytes.Buffer{}
	args := Args{
		Stdout:       out,
		DeploymentID: uid.ID(),
		Data:         b,
		Engine:       "go",
		Cache:        wazero.NewCompilationCache(),
	}
	r, err := New(context.Background(), args)
	require.Nil(t, err)
	require.Nil(t, r.Invoke(bytes.NewReader(breq)))
	// t.Log(out.String())
}
