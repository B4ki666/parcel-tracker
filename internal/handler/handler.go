package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"parcel_tracker/internal/service"
)

const (
	ParcelStatusCreated = "created"
)

type Handler struct {
	Logger        *log.Logger
	ParcelService *service.ParcelService
	// ParcelService service.ParcelService // позже
}

func NewHandler(logger *log.Logger, p *service.ParcelService) *Handler {
	return &Handler{
		Logger:        logger,
		ParcelService: p,
	}
}

func (h *Handler) GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.Logger.Println("write error:", err)
	}
}

func (h *Handler) PostCreateParcelHabdler(w http.ResponseWriter, r *http.Request) {
	var parcel service.Parcel

	err := json.NewDecoder(r.Body).Decode(&parcel)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.Trim(parcel.Number, " ") == "" {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}
	if strings.Trim(parcel.Status, " ") == "" {
		parcel.Status = ParcelStatusCreated
	}

	resParcel := h.ParcelService.CreateParcel(parcel)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(resParcel)
}

func (h *Handler) GetParcelsHandler(w http.ResponseWriter, r *http.Request) {
	parcels := h.ParcelService.GetAllParcels()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parcels)
}
