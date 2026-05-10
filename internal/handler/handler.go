package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"parcel_tracker/internal/model"
	"parcel_tracker/internal/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	Logger        *log.Logger
	ParcelService *service.ParcelService
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

func (h *Handler) PostCreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var client model.Client

	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resClient, err := h.ParcelService.CreateClient(client)
	if err != nil {
		var httpErr *model.AppError

		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resClient); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}

func (h *Handler) PostCreateParcelHandler(w http.ResponseWriter, r *http.Request) {
	var parcel model.Parcel

	err := json.NewDecoder(r.Body).Decode(&parcel)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resParcel, err := h.ParcelService.CreateParcel(parcel)

	if err != nil {
		var httpErr *model.AppError

		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resParcel); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}

func (h *Handler) GetParcelsHandler(w http.ResponseWriter, r *http.Request) {
	parcels, err := h.ParcelService.GetAllParcels()
	if err != nil {
		h.Logger.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(parcels); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}

func (h *Handler) GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	clients, err := h.ParcelService.GetAllClients()
	if err != nil {
		h.Logger.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(clients); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}

func (h *Handler) GetParcelHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parcelID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "incorrect id", http.StatusBadRequest)
		return
	}

	parcel, err := h.ParcelService.GetParcel(parcelID)
	if err != nil {
		var httpErr *model.AppError

		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		h.Logger.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(parcel); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}

func (h *Handler) DeleteParcelHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parcelID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "incorrect id", http.StatusBadRequest)
		return
	}

	err = h.ParcelService.DeleteParcel(parcelID)

	if err != nil {
		var httpErr *model.AppError

		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		h.Logger.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchParcelHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	parcelID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "incorrect id", http.StatusBadRequest)
		return
	}

	var patchParcel model.UpdateParcelRequest

	err = json.NewDecoder(r.Body).Decode(&patchParcel)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if patchParcel.Status == nil && patchParcel.Address == nil {
		http.Error(w, "nothing to update", http.StatusBadRequest)
		return
	}

	modifiedParcel, err := h.ParcelService.PatсhParcel(parcelID, patchParcel)

	if err != nil {
		var httpErr *model.AppError

		if errors.As(err, &httpErr) {
			http.Error(w, httpErr.Message, httpErr.StatusCode)
			return
		}
		h.Logger.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(modifiedParcel); err != nil {
		h.Logger.Printf("Error encoding parcel: %v", err)
	}
}
