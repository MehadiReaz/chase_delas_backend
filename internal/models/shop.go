package models

import "time"

type Shop struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	OwnerID   string    `json:"owner_id"`
	Priority  int       `json:"priority" validate:"min=0"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ShopUpdateRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
