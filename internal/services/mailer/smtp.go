package mailer

import (
	"context"
	"fmt"
	"mime"

	"github.com/segmentfault/pacman/errors"
	"gopkg.in/gomail.v2"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/pkg/logger"
)

type smtpProvider struct {
	cfg *config.Mailer
}

func newSMTP(cfg *config.Mailer) mailerProvider {
	return &smtpProvider{
		cfg: cfg,
	}
}

func (s *smtpProvider) Send(
	ctx context.Context,
	toEmailAddr string,
	subject string,
	body string,
) error {
	logger.Infof("try to send email to %s", toEmailAddr)

	msg := gomail.NewMessage()
	fromName := mime.QEncoding.Encode("utf-8", s.cfg.FromName)
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, s.cfg.FromEmail))
	msg.SetHeader("To", toEmailAddr)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.Host, s.cfg.Port, s.cfg.User, s.cfg.Password)
	if s.cfg.RequireTLS {
		d.SSL = true
	}

	err := d.DialAndSend(msg)
	if err != nil {
		return errors.InternalServer(reason.MailServerError).WithError(err).WithStack()
	}

	return nil
}
