package mail

import (
	"fmt"

	"github.com/warehouse/mail-service/internal/domain"
)

func (s *service) generateSubjectName(subject domain.EmailType) string {
	switch subject {
	case domain.VerificationType:
		return "Подтверждение электронной почты"
	case domain.ResetType:
		return "Запрос на восстановление пароля"
	}

	return ""
}

func (s *service) generateTextBody(email domain.EmailMessage) string {
	switch email.Type {
	case domain.VerificationType:
		return fmt.Sprintf(`
      Здравствуйте, %s!
      
      Для завершения регистрации перейдите по ссылке:
      %s
      
      Если вы не знаете, что это за письмо – проигнорируйте это сообщение.
      
      WarehouseAI Team
      `,
			email.Payload.Firstname,
			email.Payload.VerifyPayload.Token,
		)

	case domain.ResetType:
		return fmt.Sprintf(`
      Здравствуйте, %s!
      
      Мы получили запрос на восстановление пароля от аккаунта, связанного с почтой %s.
      Ваш код верификации: %s
      
      Если это не вы - проигнорируйте данное письмо.
      
      WarehouseAI Team
    `,
			email.Payload.Firstname,
			email.To,
			email.Payload.ResetPayload.Token,
		)
	}

	return ""
}
