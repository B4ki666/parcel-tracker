package model

import (
	"fmt"
	"time"
)

type Parcel struct {
	ID        int       `json:"id"`
	ClientId  *int      `json:"client_id"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateParcelRequest struct {
	Status  *string
	Address *string
}

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func NewError(code int, message string) error {
	return &AppError{
		StatusCode: code,
		Message:    message,
	}
}
