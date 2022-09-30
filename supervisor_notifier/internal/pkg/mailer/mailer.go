package mailer

import "ur-services/spv-notif/internal/pkg/models"

type Mailer interface {
	SendMail(mail models.Mail) error
}
