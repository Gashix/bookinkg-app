package repository

import (
	"errors"
	"time"

	"booking-app/internal/model"
)

// интерфейс для управления данными бронирования.
type BookingRepository interface {
	CreateOrder(order model.Order) error
	GetRoomAvailability(hotelID, roomID string, dates []time.Time) (map[time.Time]*model.RoomAvailability, error)
	UpdateRoomAvailability(availability map[time.Time]*model.RoomAvailability) error
}

type bookingRepository struct {
	orders       []model.Order
	availability map[string]map[string][]model.RoomAvailability // Availability by hotelID and roomID
}

func NewBookingRepository() BookingRepository {
	return &bookingRepository{
		orders: []model.Order{},
		availability: map[string]map[string][]model.RoomAvailability{
			"reddison": {
				"lux": {
					{Date: model.Date(2024, 1, 1), Quota: 1},
					{Date: model.Date(2024, 1, 2), Quota: 1},
					{Date: model.Date(2024, 1, 3), Quota: 1},
					{Date: model.Date(2024, 1, 4), Quota: 1},
					{Date: model.Date(2024, 1, 5), Quota: 0},
				},
			},
		},
	}
}

// сохраняет новый заказ на бронирование.
func (r *bookingRepository) CreateOrder(order model.Order) error {
	r.orders = append(r.orders, order)
	return nil
}

// возвращает информацию о доступности комнат для заданных дат.
func (r *bookingRepository) GetRoomAvailability(hotelID, roomID string, dates []time.Time) (map[time.Time]*model.RoomAvailability, error) {
	if r.availability[hotelID] == nil || r.availability[hotelID][roomID] == nil {
		return nil, errors.New("hotel or room not found")
	}

	availabilityMap := make(map[time.Time]*model.RoomAvailability)
	for _, date := range dates {
		found := false
		for i, avail := range r.availability[hotelID][roomID] {
			if avail.Date.Equal(date) {
				availabilityMap[date] = &r.availability[hotelID][roomID][i]
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("room not available for some of the selected dates")
		}
	}

	return availabilityMap, nil
}

// обновляет информацию о доступности комнат.
func (r *bookingRepository) UpdateRoomAvailability(availability map[time.Time]*model.RoomAvailability) error {
	for date, avail := range availability {
		for i, storedAvail := range r.availability[avail.HotelID][avail.RoomID] {
			if storedAvail.Date.Equal(date) {
				r.availability[avail.HotelID][avail.RoomID][i] = *avail
			}
		}
	}
	return nil
}
