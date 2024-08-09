package main

import (
	"net/http"
	"os"

	"booking-app/internal/httpHandler"
	"booking-app/pkg/logging"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	handler := httpHandler.NewHandler()
	handler.RegisterRoutes(r)

	logging.Infof("Сервер слушает на localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		if err == http.ErrServerClosed {
			logging.Infof("Сервер закрыт")
		} else {
			logging.Errorf("Ошибка сервера: %s", err)
			os.Exit(1)
		}
	}
}
