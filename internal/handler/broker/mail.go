package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/warehouse/mail-service/internal/broker"
	"github.com/warehouse/mail-service/internal/domain"
	"github.com/warehouse/mail-service/internal/pkg/errors"
	"github.com/warehouse/mail-service/internal/pkg/logger"
	mailSvc "github.com/warehouse/mail-service/internal/service/mail"
	"gopkg.in/gomail.v2"
)

type (
	mailHandler struct {
		log         logger.Logger
		broker      broker.RabbitClient
		dialer      *gomail.Dialer
		sender      string
		mailService mailSvc.Service
	}
)

func NewMailHandler(
	log logger.Logger,
	broker broker.RabbitClient,
	sender string,
	dialer *gomail.Dialer,
	mailService mailSvc.Service,
) MailHandler {
	return &mailHandler{
		log:         log,
		broker:      broker,
		dialer:      dialer,
		sender:      sender,
		mailService: mailService,
	}
}

func (m *mailHandler) SendEmail() error {
	messages, err := m.broker.Chan.Consume(
		m.broker.Queues["mail"].Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("can't consume messages from %s queue: %w", m.broker.Queues["mail"].Name, err)
	}

	stop := make(chan bool)

	go func() {
		for message := range messages {
			var emailEvent domain.EmailMessage
			if err := json.Unmarshal(message.Body, &emailEvent); err != nil {
				m.log.ServiceError(errors.WD(errors.InternalError, err))
			}

			// TODO: добавить таймер на выполнение

			if e := m.mailService.SendEmail(context.Background(), emailEvent); e != nil {
				m.log.ServiceError(e)
			}
		}
	}()

	<-stop
	return nil
}
