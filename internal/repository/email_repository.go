package repository

import (
	"booking-app/internal/model"
	"booking-app/pkg/logging"
	"errors"
)

// интерфейс для управления логами писем.
type EmailRepository interface {
	LogEmail(email string, order model.Order) error
	GetEmailLogs() []EmailLog
}

type emailRepository struct {
	emailLogs []EmailLog
}

// хранит информацию об отправленном письме.
type EmailLog struct {
	Email string
	Order model.Order
}

func NewEmailRepository() EmailRepository {
	return &emailRepository{
		emailLogs: []EmailLog{},
	}
}

// логирует отправленное письмо и сохраняет его в памяти.
func (r *emailRepository) LogEmail(email string, order model.Order) error {
	if email == "" {
		return errors.New("email address is empty")
	}
	emailLog := EmailLog{
		Email: email,
		Order: order,
	}
	r.emailLogs = append(r.emailLogs, emailLog)
	logging.Infof("Email logged for: %s, Order: %+v", email, order)
	return nil
}

// возвращает все логи отправленных писем.
func (r *emailRepository) GetEmailLogs() []EmailLog {
	return r.emailLogs
}
