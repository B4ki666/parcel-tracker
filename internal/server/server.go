package server

import (
	"log"
	"net/http"
	"parcel_tracker/internal/handler"
	"parcel_tracker/internal/service"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	Logger *log.Logger
	Router *chi.Mux
	Server *http.Server
	// позже сюда:
	// parcelService service.ParcelService
}

func NewServer(logger *log.Logger, parcelService *service.ParcelService) *Server {
	r := chi.NewRouter()

	h := handler.NewHandler(logger, parcelService)

	srv := &Server{
		Logger: logger,
		Router: r,
	}

	//Routes
	r.Get("/health", h.GetHealthHandler)
	r.Post("/parcel", h.PostCreateParcelHandler)
	r.Get("/parcels", h.GetParcelsHandler)

	srv.Server = &http.Server{
		Addr:         ":8080",
		Handler:      srv.Router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return srv
}

func (s *Server) Start() error {
	s.Logger.Println("server started on :8080")

	return s.Server.ListenAndServe()
}
