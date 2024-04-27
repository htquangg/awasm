package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"os"
	"time"

	"github.com/htquangg/a-wasm/config"
	"github.com/htquangg/a-wasm/internal/base/handler"
	"github.com/htquangg/a-wasm/internal/base/reason"
	"github.com/htquangg/a-wasm/internal/base/translator"
	"github.com/htquangg/a-wasm/internal/constants"
	"github.com/htquangg/a-wasm/internal/schemas"
	"github.com/htquangg/a-wasm/pkg/logger"

	"github.com/segmentfault/pacman/errors"
	"gopkg.in/gomail.v2"
)

type (
	MailerRepo interface {
		SetCode(ctx context.Context, code, content string, duration time.Duration) error
		GetCode(ctx context.Context, code string) (string, error)
		DeleteCode(ctx context.Context, code string) error
	}

	MailerService struct {
		cfg        *config.Config
		mailerRepo MailerRepo
	}
)

func NewMailerService(cfg *config.Config, mailerRepo MailerRepo) *MailerService {
	return &MailerService{
		cfg:        cfg,
		mailerRepo: mailerRepo,
	}
}

func (s *MailerService) SendAndSaveCode(
	ctx context.Context,
	toEmailAddr string,
	subject string,
	body string,
	codeContent *schemas.EmailCodeContent,
) error {
	err := s.Send(ctx, toEmailAddr, subject, body)
	if err != nil {
		return err
	}

	// TODO: hash code
	codeHash := fmt.Sprintf("%s.%s", toEmailAddr, codeContent.SourceType)
	err = s.mailerRepo.SetCode(ctx, codeHash, codeContent.ToJSONString(), 10*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (s *MailerService) SendAndSaveCodeWithTime(
	ctx context.Context,
	toEmailAddr string,
	subject string,
	body string,
	codeContent *schemas.EmailCodeContent,
	duration time.Duration,
) error {
	// TOIMPROVE: rate limit to prevent spam
	err := s.Send(ctx, toEmailAddr, subject, body)
	if err != nil {
		return err
	}

	// TODO: hash code
	codeHash := fmt.Sprintf("%s.%s", toEmailAddr, codeContent.SourceType)
	err = s.mailerRepo.SetCode(ctx, codeHash, codeContent.ToJSONString(), duration)
	if err != nil {
		return err
	}

	return nil
}

func (s *MailerService) Send(ctx context.Context, toEmailAddr, subject, body string) error {
	logger.Infof("try to send email to %s", toEmailAddr)

	msg := gomail.NewMessage()
	fromName := mime.QEncoding.Encode("utf-8", s.cfg.SMTP.FromName)
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, s.cfg.SMTP.FromEmail))
	msg.SetHeader("To", toEmailAddr)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.SMTP.Host, s.cfg.SMTP.Port, s.cfg.SMTP.User, s.cfg.SMTP.Password)
	if s.cfg.SMTP.RequireTLS {
		d.SSL = true
	}
	if len(os.Getenv("SKIP_SMTP_TLS_VERIFY")) > 0 {
		d.TLSConfig = &tls.Config{ServerName: d.Host, InsecureSkipVerify: true}
	}

	err := d.DialAndSend(msg)
	if err != nil {
		return errors.InternalServer(reason.MailServerError).WithError(err).WithStack()
	}

	return nil
}

func (s *MailerService) EmailVerificationTemplate(
	ctx context.Context,
	code string,
) (string, string, error) {
	templateData := &schemas.EmailVerificationTemplateData{
		Code: code,
	}

	lang := handler.GetLangByCtx(ctx)
	title := translator.TrWithData(lang, constants.EmailTplEmailVerificationTitle, templateData)
	body := translator.TrWithData(lang, constants.EmailTplEmailVerificationBody, templateData)

	return title, body, nil
}

func (s *MailerService) VerifyCode(
	ctx context.Context,
	toEmailAddr string,
	sourceType schemas.EmailSourceType,
	code string,
) (bool, error) {
	// TOIMPROVE: rate limit to prevent spam, brute force
	// TODO: hash code
	codeHash := fmt.Sprintf("%s.%s", toEmailAddr, sourceType)
	content, err := s.mailerRepo.GetCode(ctx, codeHash)
	if err != nil {
		return false, err
	}

	codeContent := &schemas.EmailCodeContent{}
	err = codeContent.FromJSONString(content)
	if err != nil {
		return false, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	isExpired := time.Now().UTC().After(time.Unix(codeContent.ExpiresAt, 0))
	if isExpired {
		return false, errors.BadRequest(reason.OTPExpired)
	}

	if codeContent.Code == code {
		return true, nil
	}

	return false, nil
}
