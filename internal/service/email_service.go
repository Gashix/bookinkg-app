package service

import (
	"booking-app/internal/model"
	"booking-app/internal/repository"
	"fmt"
)

type EmailService interface {
	SendBookingConfirmation(email string, order model.Order) error
}

type emailService struct {
	repo repository.EmailRepository
}

func NewEmailService() EmailService {
	return &emailService{
		repo: repository.NewEmailRepository(),
	}
}

func (s *emailService) SendBookingConfirmation(email string, order model.Order) error {
	// Симуляяция отправки email
	fmt.Printf("Sending booking confirmation to %s for order: %+v\n", email, order)

	// логируем письмо
	err := s.repo.LogEmail(email, order)
	if err != nil {
		return fmt.Errorf("failed to log email: %w", err)
	}

	return nil
}
