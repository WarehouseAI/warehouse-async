package server

import (
	"sync"

	internalBroker "github.com/warehouse/mail-service/internal/broker"
	"github.com/warehouse/mail-service/internal/config"
	"github.com/warehouse/mail-service/internal/handler/broker"
	"github.com/warehouse/mail-service/internal/pkg/logger"
	"go.uber.org/zap"
)

type brokerServer struct {
	log logger.Logger
	cfg config.Rabbit
	wg  sync.WaitGroup

	broker      *internalBroker.RabbitClient
	mailHandler broker.MailHandler
}

func (g *brokerServer) Start() {
	g.log.Zap().Info("Connect to broker server", zap.String("host", g.cfg.URL))

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if err := g.mailHandler.SendEmail(); err != nil {
			g.log.Zap().Panic("Error while connecting to broker server", zap.Error(err))
		}
	}()
}

func (g *brokerServer) Stop() error {
	g.log.Zap().Info("Shutting down connection to broker")

	if err := g.broker.Chan.Close(); err != nil {
		return err
	}
	if err := g.broker.Conn.Close(); err != nil {
		return err
	}

	g.wg.Wait()
	return nil
}

func NewBrokerServer(
	log logger.Logger,
	cfg config.Config,

	broker *internalBroker.RabbitClient,
	mailHandler broker.MailHandler,
) (Server, error) {
	return &brokerServer{
		log:         log,
		cfg:         cfg.Rabbit,
		broker:      broker,
		mailHandler: mailHandler,
	}, nil
}
