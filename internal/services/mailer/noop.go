package mailer

import (
	"context"

	"github.com/htquangg/a-wasm/pkg/logger"
)

type noopProvider struct{}

func newNoop() mailerProvider {
	return &noopProvider{}
}

func (n *noopProvider) Send(
	ctx context.Context,
	toEmailAddr string,
	subject string,
	body string,
) error {
	logFields := logger.Fields{
		"email": toEmailAddr,
	}
	logger.Debugw("Noop send email", logFields)

	return nil
}
