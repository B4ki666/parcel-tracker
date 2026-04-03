package handler

import (
	"log"
	"net/http"
)

type Handler struct {
	Logger *log.Logger
	// ParcelService service.ParcelService // позже
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		Logger: logger,
	}
}

func (h *Handler) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.Logger.Println("write error:", err)
	}
}
