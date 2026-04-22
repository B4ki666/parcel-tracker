package model

import "time"

type Parcel struct {
	ID         int       `json:"id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Address    string    `json:"address"`
	Created_at time.Time `json:"created_at"`
}
