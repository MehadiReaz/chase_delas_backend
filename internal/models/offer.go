package models

import "time"

type Offer struct {
	ID            string    `json:"id"`
	ShopID        string    `json:"shop_id" validate:"required"`
	Title         string    `json:"title" validate:"required,min=5"`
	Description   string    `json:"description"`
	DiscountValue float64   `json:"discount_value" validate:"required,gt=0"`
	StartDate     time.Time `json:"start_date" validate:"required"`
	EndDate       time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	ImageURL      string    `json:"image_url"`
	IsFeatured    bool      `json:"is_featured"`
	CreatedAt     time.Time `json:"created_at"`
}

type OfferUpdateRequest struct {
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	DiscountValue float64   `json:"discount_value"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}
