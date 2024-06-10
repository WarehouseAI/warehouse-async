package dependencies

import (
	"github.com/warehouse/mail-service/internal/broker"

	"go.uber.org/zap"
)

func (d *dependencies) RabbitClient() *broker.RabbitClient {
	if d.rabbitClient == nil {
		var err error
		msg := "connection to rabbitmq broker"
		if d.rabbitClient, err = broker.NewRabbitClient(d.cfg.Rabbit.URL, d.cfg.Rabbit.MailQueue, d.cfg.Rabbit.UserQueue); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}
		d.closeCallbacks = append(d.closeCallbacks, func() {
			if err := d.rabbitClient.Conn.Close(); err != nil {
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
		})
	}
	return d.rabbitClient
}
