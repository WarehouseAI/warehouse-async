package mail

import (
	"context"

	"github.com/warehouse/mail-service/internal/domain"
	"github.com/warehouse/mail-service/internal/pkg/errors"
	"github.com/warehouse/mail-service/internal/pkg/logger"
	"gopkg.in/gomail.v2"
)

type (
	Service interface {
		SendEmail(ctx context.Context, email domain.EmailMessage) *errors.Error
	}

	service struct {
		dialer *gomail.Dialer
		log    logger.Logger
		sender string
	}
)

func NewMailService(
	dialer *gomail.Dialer,
	log logger.Logger,
	sender string,
) Service {
	return &service{
		dialer: dialer,
		log:    log,
		sender: sender,
	}
}

func (s *service) SendEmail(ctx context.Context, email domain.EmailMessage) *errors.Error {
	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(s.sender, "WarehouseAI Team"))
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", s.generateSubjectName(email.Type))
	m.SetBody("text/plain", s.generateTextBody(email))

	if err := s.dialer.DialAndSend(m); err != nil {
		return errors.WD(errors.InternalError, err)
	}

	return nil
}
