package dependencies

import (
	"github.com/warehouse/mail-service/internal/handler/broker"
)

func (d *dependencies) MailHandler() broker.MailHandler {
	if d.mailHandler == nil {
		d.mailHandler = broker.NewMailHandler(
			d.log,
			*d.rabbitClient,
			d.Cfg().Mail.Email,
			d.MailDialer(),
			d.MailService(),
		)
	}

	return d.mailHandler
}
