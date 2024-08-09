package httpHandler

import (
	"encoding/json"
	"net/http"
	"time"

	"booking-app/internal/model"
	"booking-app/internal/service"
	"booking-app/pkg/logging"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	bookingService service.BookingService
}

func NewHandler() *Handler {
	return &Handler{
		bookingService: service.NewBookingService(),
	}
}

// регистрирует HTTP маршруты и их обработчики.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/orders", h.createOrder)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		logging.Errorf("Ошибка при декодировании тела запроса: %v", err)
		return
	}

	// Валидация дат заказа
	if newOrder.From.After(newOrder.To) || newOrder.From.Before(time.Now()) {
		http.Error(w, "Некорректный диапазон дат", http.StatusBadRequest)
		logging.Errorf("Некорректный диапазон дат: от %v до %v", newOrder.From, newOrder.To)
		return
	}

	if err := h.bookingService.CreateOrder(newOrder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logging.Errorf("Не удалось создать заказ: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newOrder); err != nil {
		logging.Errorf("Ошибка при кодировании ответа: %v", err)
	}

	logging.Infof("Заказ успешно создан: %v", newOrder)
}
