package service

import (
	"errors"
	"sync"
	"time"

	"booking-app/internal/model"
	"booking-app/internal/repository"
	"booking-app/pkg/utils"
)

// определяет методы для управления бронированием.
type BookingService interface {
	CreateOrder(order model.Order) error
}

// реализация BookingService.
type bookingService struct {
	repo         repository.BookingRepository
	emailService EmailService
	mu           sync.Mutex
}

func NewBookingService() BookingService {
	return &bookingService{
		repo:         repository.NewBookingRepository(),
		emailService: NewEmailService(),
	}
}

// создает новый заказ на бронирование номера.
func (s *bookingService) CreateOrder(order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка на корректность дат
	if order.From.After(order.To) || order.From.Before(time.Now()) {
		return errors.New("invalid date range")
	}

	daysToBook := utils.DaysBetween(order.From, order.To)
	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	// Получение информации о доступности комнат
	availability, err := s.repo.GetRoomAvailability(order.HotelID, order.RoomID, daysToBook)
	if err != nil {
		return err
	}

	// Проверка доступности номеров
	for _, dayToBook := range daysToBook {
		if avail, ok := availability[dayToBook]; !ok || avail.Quota < 1 {
			return errors.New("hotel room is not available for selected dates")
		}
		availability[dayToBook].Quota -= 1
		delete(unavailableDays, dayToBook)
	}

	// Обновление информации о доступности номеров
	if err := s.repo.UpdateRoomAvailability(availability); err != nil {
		return err
	}

	// Создание заказа
	if err := s.repo.CreateOrder(order); err != nil {
		return err
	}

	// Отправка подтверждения бронирования по электронной почте
	if err := s.emailService.SendBookingConfirmation(order.UserEmail, order); err != nil {
		return err
	}

	return nil
}
