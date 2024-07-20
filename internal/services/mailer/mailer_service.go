package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentfault/pacman/errors"

	"github.com/htquangg/awasm/config"
	"github.com/htquangg/awasm/internal/base/handler"
	"github.com/htquangg/awasm/internal/base/reason"
	"github.com/htquangg/awasm/internal/base/translator"
	"github.com/htquangg/awasm/internal/constants"
	"github.com/htquangg/awasm/internal/schemas"
)

type (
	MailerRepo interface {
		SetCode(ctx context.Context, code, content string, duration time.Duration) error
		GetCode(ctx context.Context, code string) (string, error)
		DeleteCode(ctx context.Context, code string) error
	}

	MailerService struct {
		cfg        *config.Config
		provider   mailerProvider
		mailerRepo MailerRepo
	}

	mailerProvider interface {
		Send(ctx context.Context, toEmailAddr, subject, body string) error
	}
)

func NewMailerService(cfg *config.Config, mailerRepo MailerRepo) *MailerService {
	provider := providerFor(cfg.Mailer)

	return &MailerService{
		cfg:        cfg,
		provider:   provider,
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
	return s.provider.Send(ctx, toEmailAddr, subject, body)
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

func providerFor(cfg *config.Mailer) mailerProvider {
	switch typ := cfg.ProviderType; typ {
	case config.ProviderTypeSMTP:
		return newSMTP(cfg)
	case config.ProviderTypeNoop:
		fallthrough
	default:
		return newNoop()
	}
}
