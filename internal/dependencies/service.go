package dependencies

import (
	"github.com/warehouse/mail-service/internal/service/mail"
)

func (d *dependencies) MailService() mail.Service {
	if d.mailService == nil {
		d.mailService = mail.NewMailService(
			d.MailDialer(),
			d.log,
			d.cfg.Mail.Email,
		)
	}

	return d.mailService
}
