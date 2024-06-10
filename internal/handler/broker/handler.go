package broker

type (
	MailHandler interface {
		SendEmail() error
	}
)
